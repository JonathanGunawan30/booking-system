package clients

import (
	"order-service/clients/config"
	clients3 "order-service/clients/field"
	clients2 "order-service/clients/payment"
	clients "order-service/clients/user"
	config2 "order-service/config"
)

type ClientRegistry struct {
}

type ClientRegistryInterface interface {
	GetUser() clients.UserClientInterface
	GetPayment() clients2.PaymentClientInterface
	GetField() clients3.FieldClientInterface
}

func NewClientRegistry() ClientRegistryInterface {
	return &ClientRegistry{}
}

func (r *ClientRegistry) GetUser() clients.UserClientInterface {
	return clients.NewUserClient(config.NewClientConfig(
		config.WithBaseURL(config2.AppConfig.InternalService.User.Host),
		config.WithSignatureKey(config2.AppConfig.SignatureKey),
	))
}

func (r *ClientRegistry) GetPayment() clients2.PaymentClientInterface {
	return clients2.NewPaymentClient(config.NewClientConfig(
		config.WithBaseURL(config2.AppConfig.InternalService.Payment.Host),
		config.WithSignatureKey(config2.AppConfig.SignatureKey),
	))
}

func (r *ClientRegistry) GetField() clients3.FieldClientInterface {
	return clients3.NewFieldClient(config.NewClientConfig(
		config.WithBaseURL(config2.AppConfig.InternalService.Field.Host),
		config.WithSignatureKey(config2.AppConfig.SignatureKey),
	))
}
