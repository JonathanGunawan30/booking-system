package services

import (
	"order-service/clients"
	"order-service/repositories"
	services "order-service/services/order"
)

type Registry struct {
	repository repositories.RegistryRepositoryInterface
	client     clients.ClientRegistryInterface
}

type ServiceRegistryInterface interface {
	GetOrder() services.OrderServiceInterface
}

func NewServiceRegistry(repository repositories.RegistryRepositoryInterface, client clients.ClientRegistryInterface) ServiceRegistryInterface {
	return &Registry{repository: repository, client: client}
}

func (r *Registry) GetOrder() services.OrderServiceInterface {
	return services.NewOrderService(r.repository, r.client)
}
