package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"payment-service/clients"
	"payment-service/clients/midtrans"
	"payment-service/common/cloudflare"
	"payment-service/common/response"
	"payment-service/config"
	_ "payment-service/docs"
	kafka2 "payment-service/controllers/kafka"
	controllers "payment-service/controllers/payment"
	"payment-service/domain/models"
	middleware "payment-service/middlewares"
	"payment-service/repositories"
	"payment-service/routes"
	devRoute "payment-service/routes/dev"
	"payment-service/services"
	"syscall"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Payment Service API
// @version 1.0
// @description A robust, scalable, and professional payment gateway integration service.
// @host localhost:8003
// @BasePath /api/v1
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
			&models.Payment{},
			models.PaymentHistory{},
		)
		if err != nil {
			logrus.Errorf("failed to migrate database: %v", err)
			return
		}

		s3 := config.InitR2()
		r2 := cloudflare.NewR2Client(s3)
		kafka := kafka2.NewKafkaRegistry(config.AppConfig.Kafka.Brokers)
		midtrans := midtrans.NewMidtransClient(
			config.AppConfig.Midtrans.ServerKey,
			config.AppConfig.Midtrans.Production,
			config.AppConfig.Midtrans.SuccessCallbackURL,
		)

		client := clients.NewClientRegistry()

		repository := repositories.NewRepositoryRegistry(db)
		services := services.NewServiceRegistry(repository, *r2, kafka, midtrans)
		controllers := controllers.NewRegistryController(services)

		router := gin.Default()
		router.Use(middleware.HandlePanic())
		router.NoRoute(func(c *gin.Context) {
			message := fmt.Sprintf("route %s %s not found", c.Request.Method, c.Request.URL.Path)
			response.Error(c, http.StatusNotFound, nil, &message, nil)
		})

		router.GET("/", func(c *gin.Context) {
			response.Success(c, http.StatusOK, "Hello From Payment Service", nil)
		})

		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
