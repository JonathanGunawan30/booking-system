package repositories

import (
	"context"
	error2 "order-service/common/error"
	constants "order-service/constants/error"
	"order-service/domain/dto"
	"order-service/domain/models"

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
	history := models.OrderHistory{
		OrderID: request.OrderID,
		Status:  request.Status,
	}
	result := db.WithContext(ctx).Create(&history)
	if result.Error != nil {
		return error2.WrapError(constants.ErrSQLError)
	}
	return nil
}
