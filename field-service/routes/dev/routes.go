package routes

import (
	"crypto/sha256"
	"encoding/hex"
	"field-service/common/response"
	"field-service/config"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func RegisterDevRoutes(router *gin.Engine) {
	if config.AppConfig.AppEnv != "development" {
		return
	}

	devGroup := router.Group("/dev")
	{
		devGroup.GET("/api-key", generateAPIKeyHandler)
	}
}

func generateAPIKeyHandler(c *gin.Context) {
	serviceName := c.Query("service_name")
	requestAt := time.Now().UTC().Format(time.RFC3339)

	validateKey := fmt.Sprintf("%s:%s:%s", serviceName, config.AppConfig.SignatureKey, requestAt)
	hash := sha256.New()
	hash.Write([]byte(validateKey))
	apiKey := hex.EncodeToString(hash.Sum(nil))

	response.Success(c, http.StatusOK, gin.H{
		"x-api-key":      apiKey,
		"x-service-name": serviceName,
		"x-request-at":   requestAt,
	}, nil)
}
