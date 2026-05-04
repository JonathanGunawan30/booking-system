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
	kafka2 "order-service/controllers/kafka"
	kafka "order-service/controllers/kafka/config"
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

	"github.com/IBM/sarama"
	"github.com/didip/tollbooth"
	"github.com/didip/tollbooth/limiter"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var command = &cobra.Command{
	Use:   "serve",
	Short: "Start the order service",
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
			&models.Order{},
			models.OrderHistory{},
			models.OrderField{},
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

		wg.Add(1)
		go serveKafkaConsumer(ctx, wg, service)

		logrus.Info("Application is running...")

		<-ctx.Done()
		logrus.Info("Shutting down application...")

		wg.Wait()
		logrus.Info("Application stopped cleanly")
	},
}

func Run() {
	if err := command.Execute(); err != nil {
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
			Message: "Order Service",
		})
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

func serveKafkaConsumer(ctx context.Context, wg *sync.WaitGroup, service services.ServiceRegistryInterface) {
	defer wg.Done()

	kafkaConsumerConfig := sarama.NewConfig()
	kafkaConsumerConfig.Consumer.MaxWaitTime = time.Duration(config.AppConfig.Kafka.MaxWaitTime) * time.Millisecond
	kafkaConsumerConfig.Consumer.MaxProcessingTime = time.Duration(config.AppConfig.Kafka.MaxProcessingTime) * time.Millisecond
	kafkaConsumerConfig.Consumer.Retry.Backoff = time.Duration(config.AppConfig.Kafka.BackOffTime) * time.Millisecond
	kafkaConsumerConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	kafkaConsumerConfig.Consumer.Offsets.AutoCommit.Enable = true
	kafkaConsumerConfig.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second
	kafkaConsumerConfig.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
		sarama.NewBalanceStrategyRoundRobin(),
	}

	brokers := config.AppConfig.Kafka.Brokers
	groupID := config.AppConfig.Kafka.GroupID
	topics := config.AppConfig.Kafka.Topics

	consumerGroup, err := sarama.NewConsumerGroup(brokers, groupID, kafkaConsumerConfig)
	if err != nil {
		logrus.Errorf("failed to create consumer group: %v", err)
		return
	}

	consumer := kafka.NewConsumerGroup()
	kafkaRegistry := kafka2.NewKafkaRegistry(service)
	kafkaConsumer := kafka.NewKafkaConsumer(consumer, kafkaRegistry)
	kafkaConsumer.Register()

	go func() {
		for {
			if err := consumerGroup.Consume(ctx, topics, consumer); err != nil {
				logrus.Errorf("failed to consume message: %v", err)
				time.Sleep(2 * time.Second)
			}

			if ctx.Err() != nil {
				return
			}
		}
	}()

	logrus.Info("Kafka consumer started")

	<-ctx.Done()
	logrus.Info("Shutting down Kafka consumer...")

	if err := consumerGroup.Close(); err != nil {
		logrus.Errorf("failed to close consumer group: %v", err)
	}
	logrus.Info("Kafka consumer stopped")
}
