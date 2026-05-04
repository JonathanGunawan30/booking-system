package kafka

import (
	"context"
	"order-service/config"
	"time"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

type (
	TopicName string
	Handler   func(ctx context.Context, message *sarama.ConsumerMessage) error
)

type ConsumerGroup struct {
	handler map[TopicName]Handler
}

func NewConsumerGroup() *ConsumerGroup {
	return &ConsumerGroup{handler: make(map[TopicName]Handler)}
}

func (c *ConsumerGroup) Setup(group sarama.ConsumerGroupSession) error {
	logrus.Infof("setup consumer group: %v", group.Claims())
	return nil
}

func (c *ConsumerGroup) Cleanup(group sarama.ConsumerGroupSession) error {
	logrus.Infof("cleanup consumer group: %v", group.Claims())
	return nil
}

func (c *ConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	messages := claim.Messages()
	for message := range messages {
		handler, ok := c.handler[TopicName(message.Topic)]
		if !ok {
			logrus.Errorf("topic %s not found", message.Topic)
			continue
		}

		var err error
		maxRetry := config.AppConfig.Kafka.MaxRetry
		for attempt := 1; attempt <= maxRetry; attempt++ {
			err = handler(context.Background(), message)
			if err == nil {
				break
			}
			logrus.Errorf("error handling message on topic %s: %v, attempt: %d", message.Topic, err, attempt)
			if attempt == maxRetry {
				logrus.Errorf("max retry reached, skipping message on topic %s", message.Topic)
			}
		}

		if err != nil {
			logrus.Errorf("error handling message on topic %s: %v", message.Topic, err)
			break
		}

		session.MarkMessage(message, time.Now().UTC().String())
	}
	return nil
}

func (c *ConsumerGroup) RegisterHandler(topic TopicName, handler Handler) {
	c.handler[topic] = handler
	logrus.Infof("register handler for topic %s", topic)
}
