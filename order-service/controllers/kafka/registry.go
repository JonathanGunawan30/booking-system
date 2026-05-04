package kafka

import (
	kafka "order-service/controllers/kafka/payment"
	"order-service/services"
)

type Registry struct {
	service services.ServiceRegistryInterface
}

type KafkaRegistryInterface interface {
	GetPayment() kafka.PaymentKafkaInterface
}

func NewKafkaRegistry(service services.ServiceRegistryInterface) KafkaRegistryInterface {
	return &Registry{service: service}
}

func (r *Registry) GetPayment() kafka.PaymentKafkaInterface {
	return kafka.NewPaymentKafka(r.service)
}
