package controllers

import "payment-service/services"

type Registry struct {
	service services.ServiceRegistryInterface
}

type RegistryControllerInterface interface {
	GetPayment() PaymentControllerInterface
}

func NewRegistryController(service services.ServiceRegistryInterface) RegistryControllerInterface {
	return &Registry{service: service}
}

func (r *Registry) GetPayment() PaymentControllerInterface {
	return NewPaymentController(r.service)
}
