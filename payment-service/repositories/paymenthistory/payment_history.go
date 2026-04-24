package paymenthistory

import (
	"context"
	errWrap "payment-service/common/error"
	constants "payment-service/constants/error"
	"payment-service/domain/dto"
	"payment-service/domain/models"

	"gorm.io/gorm"
)

type PaymentHistoryRepository struct {
	db *gorm.DB
}

type PaymentHistoryRepositoryInterface interface {
	Create(ctx context.Context, db *gorm.DB, request *dto.PaymentHistoryRequest) error
}

func NewPaymentHistory(db *gorm.DB) PaymentHistoryRepositoryInterface {
	return &PaymentHistoryRepository{db: db}
}

func (p *PaymentHistoryRepository) Create(ctx context.Context, db *gorm.DB, request *dto.PaymentHistoryRequest) error {
	paymentHistory := models.PaymentHistory{
		PaymentID: request.PaymentID,
		Status:    request.Status,
	}

	err := db.WithContext(ctx).Create(&paymentHistory).Error
	if err != nil {
		return errWrap.WrapError(constants.ErrSQLError)
	}
	return nil
}
