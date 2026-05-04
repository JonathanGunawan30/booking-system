package repositories

import (
	"context"
	error2 "order-service/common/error"
	constants "order-service/constants/error"
	"order-service/domain/models"

	"gorm.io/gorm"
)

type OrderFieldRepository struct {
	db *gorm.DB
}

type OrderFieldRepositoryInterface interface {
	FindByOrderID(ctx context.Context, orderID uint) ([]models.OrderField, error)
	Create(ctx context.Context, db *gorm.DB, req []models.OrderField) error
}

func NewOrderFieldRepository(db *gorm.DB) OrderFieldRepositoryInterface {
	return &OrderFieldRepository{db: db}
}

func (o *OrderFieldRepository) FindByOrderID(ctx context.Context, orderID uint) ([]models.OrderField, error) {
	var orderFields []models.OrderField
	err := o.db.WithContext(ctx).Where("order_id = ?", orderID).Find(&orderFields).Error
	if err != nil {
		return nil, error2.WrapError(constants.ErrSQLError)
	}
	return orderFields, nil
}

func (o *OrderFieldRepository) Create(ctx context.Context, db *gorm.DB, req []models.OrderField) error {
	err := db.WithContext(ctx).Create(&req).Error
	if err != nil {
		return error2.WrapError(constants.ErrSQLError)
	}
	return nil
}
