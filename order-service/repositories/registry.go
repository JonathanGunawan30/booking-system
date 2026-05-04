package repositories

import (
	repositories "order-service/repositories/order"
	repositories2 "order-service/repositories/orderfield"
	repositories3 "order-service/repositories/orderhistory"

	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

type RegistryRepositoryInterface interface {
	GetOrder() repositories.OrderRepositoryInterface
	GetOrderField() repositories2.OrderFieldRepositoryInterface
	GetOrderHistory() repositories3.OrderHistoryRepositoryInterface
	GetTx() *gorm.DB
}

func NewRepositoryRegistry(db *gorm.DB) RegistryRepositoryInterface {
	return &Registry{db: db}
}

func (r *Registry) GetOrder() repositories.OrderRepositoryInterface {
	return repositories.NewOrderRepository(r.db)
}

func (r *Registry) GetOrderField() repositories2.OrderFieldRepositoryInterface {
	return repositories2.NewOrderFieldRepository(r.db)
}

func (r *Registry) GetOrderHistory() repositories3.OrderHistoryRepositoryInterface {
	return repositories3.NewOrderHistoryRepository(r.db)
}

func (r *Registry) GetTx() *gorm.DB {
	return r.db
}
