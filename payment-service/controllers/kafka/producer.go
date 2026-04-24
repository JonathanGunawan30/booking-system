package kafka

import (
	configApp "payment-service/config"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type Kafka struct {
	brokers []string
}

type KafkaInterface interface {
	ProduceMessage(string, []byte) error
}

func NewKafkaProducer(brokers []string) KafkaInterface {
	return &Kafka{brokers: brokers}
}

func (k *Kafka) ProduceMessage(topic string, bytes []byte) error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = configApp.AppConfig.Kafka.MaxRetry
	producer, err := sarama.NewSyncProducer(k.brokers, config)
	if err != nil {
		logrus.Errorf("Failed to create a new producer: %v", err)
		return err
	}
	defer producer.Close()

	message := &sarama.ProducerMessage{
		Topic:   topic,
		Headers: nil,
		Value:   sarama.ByteEncoder(bytes),
	}

	partition, offset, err := producer.SendMessage(message)
	if err != nil {
		logrus.Errorf("Failed to send message: %v", err)
		return err
	}
	logrus.Infof("Message sent to partition %d at offset %d", partition, offset)
	return nil
}
