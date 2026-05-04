package repositories

import (
	"context"
	error2 "order-service/common/error"
	constants "order-service/constants/error"
	"order-service/domain/dto"

	"gorm.io/gorm"
)

type OrderHistoryRepository struct {
	db *gorm.DB
}

type OrderHistoryRepositoryInterface interface {
	Create(ctx context.Context, db *gorm.DB, request *dto.OrderHistoryRequest) error
}

func NewOrderHistoryRepository(db *gorm.DB) OrderHistoryRepositoryInterface {
	return &OrderHistoryRepository{db: db}
}

func (o *OrderHistoryRepository) Create(ctx context.Context, db *gorm.DB, request *dto.OrderHistoryRequest) error {
	result := db.WithContext(ctx).Create(request)
	if result.Error != nil {
		return error2.WrapError(constants.ErrSQLError)
	}
	return nil
}
