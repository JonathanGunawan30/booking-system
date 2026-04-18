package controllers

import (
	controller "user-service/controllers/user"
	"user-service/services"
)

type Registry struct {
	service services.ServiceRegistryInterface
}

type UserControllerRegistryInterface interface {
	GetUserController() controller.UserControllerInterface
}

func NewControllerRegistry(service services.ServiceRegistryInterface) UserControllerRegistryInterface {
	return &Registry{service: service}
}

func (r *Registry) GetUserController() controller.UserControllerInterface {
	return controller.NewUserController(r.service.GetUser())
}
