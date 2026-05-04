package controllers

import (
	controllers "order-service/controllers/http/order"
	"order-service/services"
)

type Registry struct {
	service services.ServiceRegistryInterface
}

type ControllerRegistryInterface interface {
	GetOrder() controllers.OrderControllerInterface
}

func NewControllerRegistry(service services.ServiceRegistryInterface) ControllerRegistryInterface {
	return &Registry{service: service}
}

func (r *Registry) GetOrder() controllers.OrderControllerInterface {
	return controllers.NewOrderController(r.service)
}
