package kafka

import (
	"context"
	"encoding/json"
	"order-service/common/util"
	"order-service/domain/dto"
	"order-service/services"

	"github.com/IBM/sarama"
	"github.com/sirupsen/logrus"
)

const PaymentTopic = "payment-service-callback"

type PaymentKafka struct {
	service services.ServiceRegistryInterface
}

type PaymentKafkaInterface interface {
	HandlePayment(ctx context.Context, message *sarama.ConsumerMessage) error
}

func NewPaymentKafka(service services.ServiceRegistryInterface) PaymentKafkaInterface {
	return &PaymentKafka{service: service}
}

func (p *PaymentKafka) HandlePayment(ctx context.Context, message *sarama.ConsumerMessage) error {
	defer util.Recover()
	var body dto.PaymentContent
	err := json.Unmarshal(message.Value, &body)
	if err != nil {
		logrus.Errorf("failed to unmarshal message: %s", err.Error())
		return err
	}

	data := body.Body.Data
	err = p.service.GetOrder().HandlePayment(ctx, &data)
	if err != nil {
		logrus.Errorf("failed to handle payment: %s", err.Error())
		return err
	}

	logrus.Infof("payment handled successfully")
	return nil

}
