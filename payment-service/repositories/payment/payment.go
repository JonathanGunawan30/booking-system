package repositories

import (
	"context"
	"errors"
	"fmt"
	errWrap "payment-service/common/error"
	"payment-service/constants"
	errConstant "payment-service/constants/error"
	errPayment "payment-service/constants/error/payment"
	"payment-service/domain/dto"
	"payment-service/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PaymentRepository struct {
	db *gorm.DB
}

type PaymentRepositoryInterface interface {
	FindAllWithPagination(ctx context.Context, param *dto.PaymentRequestParam) ([]models.Payment, int64, error)
	FindByUUID(ctx context.Context, uuid string) (*models.Payment, error)
	FindByOrderID(ctx context.Context, orderID string) (*models.Payment, error)
	Create(ctx context.Context, db *gorm.DB, request *dto.PaymentRequest) (*models.Payment, error)
	Update(ctx context.Context, db *gorm.DB, request *dto.UpdatePaymentRequest, orderID string) (*models.Payment, error)
}

func NewPaymentRepository(db *gorm.DB) PaymentRepositoryInterface {
	return &PaymentRepository{db: db}
}

func (p *PaymentRepository) FindAllWithPagination(ctx context.Context, param *dto.PaymentRequestParam) ([]models.Payment, int64, error) {
	var (
		fields []models.Payment
		sort   string
		total  int64
	)

	if param.SortColumn != nil {
		sort = fmt.Sprintf("%s %s", *param.SortColumn, *param.SortOrder)
	} else {
		sort = "created_at desc"
	}

	limit := param.Limit
	offset := (param.Page - 1) * limit
	err := p.db.WithContext(ctx).Limit(limit).Offset(offset).Order(sort).Find(&fields).Error
	if err != nil {
		return nil, 0, err
	}

	err = p.db.WithContext(ctx).Model(&models.Payment{}).Count(&total).Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return fields, total, nil
}

func (p *PaymentRepository) FindByUUID(ctx context.Context, uuid string) (*models.Payment, error) {
	var payment models.Payment
	result := p.db.WithContext(ctx).Where("uuid = ?", uuid).First(&payment)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errPayment.ErrPaymentNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &payment, nil
}

func (p *PaymentRepository) FindByOrderID(ctx context.Context, orderID string) (*models.Payment, error) {
	var payment models.Payment
	err := p.db.WithContext(ctx).Where("order_id = ?", orderID).First(&payment).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errPayment.ErrPaymentNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &payment, nil
}

func (p *PaymentRepository) Create(ctx context.Context, db *gorm.DB, request *dto.PaymentRequest) (*models.Payment, error) {
	status := constants.Initial
	orderID := uuid.MustParse(request.OrderID)
	payment := models.Payment{
		UUID:        uuid.New(),
		OrderID:     orderID,
		Amount:      request.Amount,
		PaymentLink: request.PaymentLink,
		ExpiredAt:   &request.ExpiredAt,
		Description: request.Description,
		Status:      &status,
	}

	err := db.WithContext(ctx).Create(&payment).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &payment, nil
}

func (p *PaymentRepository) Update(ctx context.Context, db *gorm.DB, request *dto.UpdatePaymentRequest, orderID string) (*models.Payment, error) {
	payment := models.Payment{
		Status:        request.Status,
		TransactionID: request.TransactionID,
		InvoiceLink:   request.InvoiceLink,
		PaidAt:        request.PaidAt,
		VANumber:      request.VANumber,
		Bank:          request.Bank,
		Acquirer:      request.Acquirer,
	}

	if err := db.WithContext(ctx).Where("order_id = ?", orderID).Updates(&payment).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errPayment.ErrPaymentNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &payment, nil
}
