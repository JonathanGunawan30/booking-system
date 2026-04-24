package repositories

import (
	payment "payment-service/repositories/payment"
	"payment-service/repositories/paymenthistory"

	"gorm.io/gorm"
)

type Registry struct {
	db *gorm.DB
}

type RepositoryRegistryInterface interface {
	GetPayment() payment.PaymentRepositoryInterface
	GetPaymentHistory() paymenthistory.PaymentHistoryRepositoryInterface
	GetTx() *gorm.DB
}

func NewRepositoryRegistry(db *gorm.DB) RepositoryRegistryInterface {
	return &Registry{db: db}
}

func (r *Registry) GetPayment() payment.PaymentRepositoryInterface {
	return payment.NewPaymentRepository(r.db)
}

func (r *Registry) GetPaymentHistory() paymenthistory.PaymentHistoryRepositoryInterface {
	return paymenthistory.NewPaymentHistory(r.db)
}

func (r *Registry) GetTx() *gorm.DB {
	return r.db
}
