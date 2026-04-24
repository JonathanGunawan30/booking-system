package middleware

import (
	"context"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"fmt"
	"net/http"
	"payment-service/clients"
	"payment-service/common/response"
	"payment-service/config"
	"payment-service/constants"
	errConstant "payment-service/constants/error"
	"strings"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
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
			logrus.Warnf("[RateLimiter] rate limit exceeded from %s", c.ClientIP())
			response.Error(c, http.StatusTooManyRequests, errConstant.ErrTooManyRequest, nil, err)
			c.Abort()
		}
		c.Next()
	}
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
		logrus.Warnf("[validateAPIKey] missing required headers from %s", c.ClientIP())
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

func contains(roles []string, role string) bool {
	for _, r := range roles {
		if r == role {
			return true
		}
	}
	return false
}

func CheckRole(roles []string, client clients.ClientRegistryInterface) gin.HandlerFunc {
	return func(c *gin.Context) {
		user, err := client.GetUser().GetUserByToken(c.Request.Context())
		if err != nil {
			logrus.Errorf("[CheckRole] failed to get user from user-service: %v", err)
			responseUnauthorized(c, errConstant.ErrUnauthorized.Error())
			return
		}

		if !contains(roles, user.Role) {
			logrus.Warnf("[CheckRole] unauthorized access attempt, role=%s required=%v ip=%s", user.Role, roles, c.ClientIP())
			responseUnauthorized(c, errConstant.ErrUnauthorized.Error())
			return
		}

		c.Next()
	}
}

func Authenticate() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader(constants.Authorization)
		if token == "" {
			logrus.Warnf("[Authenticate] missing authorization header ip=%s path=%s", c.ClientIP(), c.Request.URL.Path)
			responseUnauthorized(c, errConstant.ErrUnauthorized.Error())
			return
		}

		if err := validateAPIKey(c); err != nil {
			responseUnauthorized(c, errConstant.ErrUnauthorized.Error())
			return
		}

		cleanToken := strings.TrimPrefix(token, "Bearer ")
		ctx := context.WithValue(c.Request.Context(), constants.Token, cleanToken)
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}
