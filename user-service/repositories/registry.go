package repositories

import (
	"sync"
	repositories "user-service/repositories/user"

	"gorm.io/gorm"
)

type Registry struct {
	db             *gorm.DB
	userRepository repositories.UserRepositoryInterface
	userOnce       sync.Once
}

type RepositoryRegistryInterface interface {
	GetUser() repositories.UserRepositoryInterface
}

func NewRepositoryRegistry(db *gorm.DB) RepositoryRegistryInterface {
	return &Registry{db: db}
}

func (r *Registry) GetUser() repositories.UserRepositoryInterface {
	r.userOnce.Do(func() {
		r.userRepository = repositories.NewUserRepository(r.db)
	})
	return r.userRepository
}
