package services

import (
	"user-service/repositories"
	services "user-service/services/user"
)

type Registry struct {
	repository repositories.RepositoryRegistryInterface
}

type ServiceRegistryInterface interface {
	GetUser() services.UserServiceInterface
}

func NewServiceRegistry(registryInterface repositories.RepositoryRegistryInterface) ServiceRegistryInterface {
	return &Registry{repository: registryInterface}
}

func (r *Registry) GetUser() services.UserServiceInterface {
	return services.NewUserService(r.repository)
}
