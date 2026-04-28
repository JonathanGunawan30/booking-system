package cmd

import (
	"field-service/clients"
	"field-service/common/cloudflare"
	"field-service/common/response"
	"field-service/config"
	"field-service/controllers"
	"field-service/domain/models"
	"field-service/middleware"
	"field-service/repositories"
	"field-service/routes"
	devRoute "field-service/routes/dev"
	"field-service/services"
	"fmt"
	"net/http"
	"time"

	_ "field-service/docs"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Field Service API
// @version 1.0
// @description A microservice for managing booking system fields, operational times, and schedules.
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8002
// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

var command = &cobra.Command{
	Use:   "serve",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		_ = godotenv.Load()
		config.AppConfig = config.LoadConfig()
		db, err := config.InitDatabase()
		if err != nil {
			panic(err)
		}

		loc, err := time.LoadLocation(config.AppConfig.Timezone)
		if err != nil {
			panic(err)
		}

		time.Local = loc

		err = db.AutoMigrate(
			&models.Field{},
			&models.FieldSchedule{},
			models.Time{},
		)

		s3 := config.InitR2()
		r2 := cloudflare.NewR2Client(s3)

		client := clients.NewClientRegistry()

		if err != nil {
			panic(err)
		}

		repository := repositories.NewRegistryRepository(db)
		services := services.NewServiceRegistry(repository, *r2)
		controllers := controllers.NewControllerRegistry(services)

		router := gin.Default()
		router.Use(middleware.HandlePanic())
		router.NoRoute(func(c *gin.Context) {
			message := fmt.Sprintf("route %s %s not found", c.Request.Method, c.Request.URL.Path)
			response.Error(c, http.StatusNotFound, nil, &message, nil)
		})

		router.GET("/", func(c *gin.Context) {
			response.Success(c, http.StatusOK, "Hello From Field Service", nil)
		})

		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

		router.Use(func(c *gin.Context) {
			c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
			c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Authorization, x-service-name, x-request-at, x-api-key,")
			c.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Length")
			c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
			if c.Request.Method == "OPTIONS" {
				c.AbortWithStatus(http.StatusNoContent)
			} else {
				c.Next()
			}
		})

		lmt := tollbooth.NewLimiter(
			config.AppConfig.RateLimiterMaxRequest,
			&limiter.ExpirableOptions{
				DefaultExpirationTTL: time.Duration(config.AppConfig.RateLimiterTimeSecond) * time.Second,
			},
		)

		router.Use(middleware.RateLimiter(lmt))

		devRoute.RegisterDevRoutes(router)

		group := router.Group("/api/v1")
		route := routes.NewRouteRegistry(group, controllers, client)
		route.Serve()

		port := fmt.Sprintf(":%d", config.AppConfig.Port)
		_ = router.Run(port)
	},
}

func Run() {
	if err := command.Execute(); err != nil {
		fmt.Println(err)
		return
	}
}
