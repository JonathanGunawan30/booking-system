package kafka

import (
	"order-service/config"
	"order-service/controllers/kafka"
	kafka2 "order-service/controllers/kafka/payment"
	"slices"
)

type Kafka struct {
	consumer *ConsumerGroup
	kafka    kafka.KafkaRegistryInterface
}

type KafkaInterface interface {
	Register()
}

func NewKafkaConsumer(consumer *ConsumerGroup, kafka kafka.KafkaRegistryInterface) KafkaInterface {
	return &Kafka{consumer: consumer, kafka: kafka}
}

func (k *Kafka) paymentHandler() {
	if slices.Contains(config.AppConfig.Kafka.Topics, kafka2.PaymentTopic) {
		k.consumer.RegisterHandler(kafka2.PaymentTopic, k.kafka.GetPayment().HandlePayment)
	}
}

func (k *Kafka) Register() {
	k.paymentHandler()
}
