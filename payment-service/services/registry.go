package services

import (
	"payment-service/clients/midtrans"
	"payment-service/common/cloudflare"
	"payment-service/controllers/kafka"
	"payment-service/repositories"
	services "payment-service/services/payment"
)

type Registry struct {
	repository repositories.RepositoryRegistryInterface
	r2         cloudflare.R2Client
	kafka      kafka.KafkaRegistryInterface
	midtrans   midtrans.MidtransClientInterface
}

type ServiceRegistryInterface interface {
	GetPayment() services.PaymentServiceInterface
}

func NewServiceRegistry(repository repositories.RepositoryRegistryInterface, r2 cloudflare.R2Client, kafka kafka.KafkaRegistryInterface, midtrans midtrans.MidtransClientInterface) ServiceRegistryInterface {
	return &Registry{repository: repository, r2: r2, kafka: kafka, midtrans: midtrans}
}

func (r *Registry) GetPayment() services.PaymentServiceInterface {
	return services.NewPaymentService(r.repository, r.r2, r.kafka, r.midtrans)
}
