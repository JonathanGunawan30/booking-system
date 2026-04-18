package cmd

import (
	"fmt"
	"net/http"
	"time"
	"user-service/common/response"
	"user-service/config"
	"user-service/controllers"
	"user-service/database/seeders"
	"user-service/domain/models"
	"user-service/middleware"
	"user-service/repositories"
	"user-service/routes"
	routes2 "user-service/routes/dev"
	"user-service/services"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

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
			&models.Role{},
			&models.User{},
		)

		if err != nil {
			panic(err)
		}

		seeders.NewSeederRegistry(db).Run()

		repository := repositories.NewRepositoryRegistry(db)
		services := services.NewServiceRegistry(repository)
		controllers := controllers.NewControllerRegistry(services)

		router := gin.Default()
		router.Use(middleware.HandlePanic())
		router.NoRoute(func(c *gin.Context) {
			message := fmt.Sprintf("route %s %s not found", c.Request.Method, c.Request.URL.Path)
			response.Error(c, http.StatusNotFound, nil, &message, nil)
		})

		router.GET("/", func(c *gin.Context) {
			response.Success(c, http.StatusOK, "Hello From User Service", nil)
		})

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

		routes2.RegisterDevRoutes(router)

		group := router.Group("/api/v1")
		route := routes.NewRouteRegistry(controllers, group)
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
