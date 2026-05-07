package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	"github.com/sirupsen/logrus"
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
			logrus.Errorf("failed to initialize database: %v", err)
			return
		}

		loc, err := time.LoadLocation(config.AppConfig.Timezone)
		if err != nil {
			logrus.Errorf("failed to load timezone: %v", err)
			return
		}

		time.Local = loc

		err = db.AutoMigrate(
			&models.Role{},
			&models.User{},
		)
		if err != nil {
			logrus.Errorf("failed to migrate database: %v", err)
			return
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

		srv := &http.Server{
			Addr:    fmt.Sprintf(":%d", config.AppConfig.Port),
			Handler: router,
		}

		go func() {
			logrus.Infof("HTTP server starting on port %d", config.AppConfig.Port)
			if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				logrus.Errorf("HTTP server failed: %v", err)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
		<-quit

		logrus.Info("Shutting down server...")

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := srv.Shutdown(ctx); err != nil {
			logrus.Errorf("Server forced to shutdown: %v", err)
		}

		logrus.Info("Server exiting")
	},
}

func Run() {
	if err := command.Execute(); err != nil {
		fmt.Println(err)
		return
	}
}
