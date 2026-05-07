package cmd

import (
	"context"
	"order-service/clients"
	"order-service/config"
	kafka2 "order-service/controllers/kafka"
	kafka "order-service/controllers/kafka/config"
	"order-service/repositories"
	"order-service/services"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var workerCmd = &cobra.Command{
	Use:   "worker",
	Short: "Start the order Kafka worker",
	Run: func(cmd *cobra.Command, args []string) {
		initialize()

		db, err := config.InitDatabase()
		if err != nil {
			logrus.Errorf("failed to initialize database: %v", err)
			return
		}

		client := clients.NewClientRegistry()
		repository := repositories.NewRepositoryRegistry(db)
		service := services.NewServiceRegistry(repository, client)

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		wg := &sync.WaitGroup{}
		wg.Add(1)
		go serveKafkaConsumer(ctx, wg, service)

		logrus.Info("Order Worker is running...")
		<-ctx.Done()
		logrus.Info("Shutting down Order Worker...")
		wg.Wait()
		logrus.Info("Order Worker stopped cleanly")
	},
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
