package clients

import (
	"field-service/clients/config"
	clients "field-service/clients/user"
	config2 "field-service/config"
)

type ClientRegistry struct {
}

type ClientRegistryInterface interface {
	GetUser() clients.UserClientInterface
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
