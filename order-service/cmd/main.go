package cmd

import (
	"context"
	"fmt"
	"net/http"
	"order-service/clients"
	"order-service/common/response"
	"order-service/config"
	"order-service/constants"
	controllers "order-service/controllers/http"
	"order-service/domain/models"
	middleware "order-service/middlewares"
	"order-service/repositories"
	"order-service/routes"
	"order-service/services"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "order-service",
	Short: "Order Service CLI",
}

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Start the order HTTP API",
	Run: func(cmd *cobra.Command, args []string) {
		initialize()

		db, err := config.InitDatabase()
		if err != nil {
			logrus.Errorf("failed to initialize database: %v", err)
			return
		}

		err = db.AutoMigrate(
			&models.Order{},
			&models.OrderHistory{},
			&models.OrderField{},
			&models.OrderSequence{},
		)
		if err != nil {
			logrus.Errorf("failed to migrate database: %v", err)
			return
		}

		client := clients.NewClientRegistry()
		repository := repositories.NewRepositoryRegistry(db)
		service := services.NewServiceRegistry(repository, client)
		controller := controllers.NewControllerRegistry(service)

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		wg := &sync.WaitGroup{}
		wg.Add(1)
		go serveHttp(ctx, wg, controller, client)

		logrus.Info("Order API is running...")
		<-ctx.Done()
		logrus.Info("Shutting down Order API...")
		wg.Wait()
		logrus.Info("Order API stopped cleanly")
	},
}

func initialize() {
	_ = godotenv.Load()
	config.AppConfig = config.LoadConfig()

	loc, err := time.LoadLocation(config.AppConfig.Timezone)
	if err != nil {
		logrus.Errorf("failed to load timezone: %v", err)
		return
	}
	time.Local = loc
}

func init() {
	rootCmd.AddCommand(apiCmd)
	rootCmd.AddCommand(workerCmd)
}

func Run() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func serveHttp(ctx context.Context, wg *sync.WaitGroup, controller controllers.ControllerRegistryInterface, client clients.ClientRegistryInterface) {
	defer wg.Done()

	router := gin.Default()
	router.Use(middleware.HandlePanic())
	router.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, response.Response{
			Status:  constants.Error,
			Message: fmt.Sprintf("route %s not found", c.Request.URL.Path),
		})
	})

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, response.Response{
			Status:  constants.Success,
			Message: "Order Service API",
		})
	})

	lmt := tollbooth.NewLimiter(
		config.AppConfig.RateLimiterMaxRequest,
		&limiter.ExpirableOptions{
			DefaultExpirationTTL: time.Duration(config.AppConfig.RateLimiterTimeSecond) * time.Second,
		})
	router.Use(middleware.RateLimiter(lmt))

	group := router.Group("/api/v1")
	route := routes.NewRouteRegistry(controller, client, group)
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

	<-ctx.Done()
	logrus.Info("Shutting down HTTP server...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logrus.Errorf("HTTP server forced to shutdown: %v", err)
	}
	logrus.Info("HTTP server stopped")
}
