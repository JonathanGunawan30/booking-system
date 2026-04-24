package cmd

import (
	"fmt"
	"net/http"
	"payment-service/clients"
	"payment-service/clients/midtrans"
	"payment-service/common/cloudflare"
	"payment-service/common/response"
	"payment-service/config"
	kafka2 "payment-service/controllers/kafka"
	controllers "payment-service/controllers/payment"
	"payment-service/domain/models"
	middleware "payment-service/middlewares"
	"payment-service/repositories"
	"payment-service/routes"
	devRoute "payment-service/routes/dev"
	"payment-service/services"
	"time"

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
			&models.Payment{},
			models.PaymentHistory{},
		)

		s3 := config.InitR2()
		r2 := cloudflare.NewR2Client(s3)
		kafka := kafka2.NewKafkaRegistry(config.AppConfig.Kafka.Brokers)
		midtrans := midtrans.NewMidtransClient(config.AppConfig.Midtrans.ServerKey, config.AppConfig.Midtrans.Production)

		client := clients.NewClientRegistry()

		if err != nil {
			panic(err)
		}

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
			response.Success(c, http.StatusOK, "Hello From Field Service", nil)
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
