package middleware

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"net/http"
	"strings"
	"time"
	"user-service/common/response"
	"user-service/config"
	"user-service/constants"
	errConstant "user-service/constants/error"
	services "user-service/services/user"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

func HandlePanic() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logrus.Errorf("[HandlePanic] panic recovered on %s %s: %v", c.Request.Method, c.Request.URL.Path, err)
				response.Error(c, http.StatusInternalServerError, errConstant.ErrInternalServerError, nil, err)
				c.Abort()
			}
		}()
		c.Next()
	}
}

func RateLimiter(lmt *limiter.Limiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := tollbooth.LimitByRequest(lmt, c.Writer, c.Request)
		if err != nil {
			logrus.Warnf("[RateLimiter] rate limit exceeded from ip=%s", c.ClientIP())
			response.Error(c, http.StatusTooManyRequests, errConstant.ErrTooManyRequest, nil, err)
			c.Abort()
		}
		c.Next()
	}
}

func extractBearerToken(token string) string {
	const prefix = "Bearer "
	if !strings.HasPrefix(token, prefix) {
		return ""
	}
	return strings.TrimPrefix(token, prefix)
}

func responseUnauthorized(c *gin.Context, message string) {
	response.Error(c, http.StatusUnauthorized, nil, &message, nil)
	c.Abort()
}

func validateAPIKey(c *gin.Context) error {
	apiKey := c.GetHeader(constants.XApiKey)
	requestAt := c.GetHeader(constants.XRequestAt)
	serviceName := c.GetHeader(constants.XServiceName)

	if apiKey == "" || requestAt == "" || serviceName == "" {
		logrus.Warnf("[validateAPIKey] missing required headers from ip=%s", c.ClientIP())
		return errConstant.ErrUnauthorized
	}

	requestTime, err := time.Parse(time.RFC3339, requestAt)
	if err != nil {
		logrus.Warnf("[validateAPIKey] invalid requestAt format from service=%s", serviceName)
		return errConstant.ErrUnauthorized
	}

	if time.Since(requestTime) > 5*time.Minute {
		logrus.Warnf("[validateAPIKey] expired request from service=%s", serviceName)
		return errConstant.ErrUnauthorized
	}

	signatureKey := config.AppConfig.SignatureKey
	validateKey := fmt.Sprintf("%s:%s:%s", serviceName, signatureKey, requestAt)

	hash := sha256.New()
	hash.Write([]byte(validateKey))
	resultHash := hex.EncodeToString(hash.Sum(nil))

	if subtle.ConstantTimeCompare([]byte(apiKey), []byte(resultHash)) != 1 {
		logrus.Warnf("[validateAPIKey] api key mismatch from service=%s ip=%s", serviceName, c.ClientIP())
		return errConstant.ErrUnauthorized
	}

	return nil
}

func ValidateBearerToken(c *gin.Context, token string) error {
	if !strings.Contains(token, "Bearer ") {
		logrus.Warnf("[ValidateBearerToken] missing Bearer prefix from ip=%s", c.ClientIP())
		return errConstant.ErrUnauthorized
	}

	tokenString := extractBearerToken(token)
	if tokenString == "" {
		logrus.Warnf("[ValidateBearerToken] empty token after extract from ip=%s", c.ClientIP())
		return errConstant.ErrUnauthorized
	}

	claims := &services.Claims{}
	tokenJwt, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			logrus.Warnf("[ValidateBearerToken] invalid signing method from ip=%s", c.ClientIP())
			return nil, errConstant.ErrUnauthorized
		}
		return []byte(config.AppConfig.JwtSecretKey), nil
	})

	if err != nil || !tokenJwt.Valid {
		logrus.Warnf("[ValidateBearerToken] invalid or expired token from ip=%s: %v", c.ClientIP(), err)
		return errConstant.ErrUnauthorized
	}

	ctx := context.WithValue(c.Request.Context(), constants.UserLogin, claims.User)
	c.Request = c.Request.WithContext(ctx)
	return nil
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(constants.Authorization)
		if token == "" {
			logrus.Warnf("[Authenticate] missing authorization header ip=%s path=%s", c.ClientIP(), c.Request.URL.Path)
			responseUnauthorized(c, errConstant.ErrUnauthorized.Error())
			return
		}

		if err := ValidateBearerToken(c, token); err != nil {
			responseUnauthorized(c, errConstant.ErrUnauthorized.Error())
			return
		}

		if err := validateAPIKey(c); err != nil {
			responseUnauthorized(c, errConstant.ErrUnauthorized.Error())
			return
		}

		c.Next()
	}
}
