package repositories

import (
	"context"
	"errors"
	"fmt"
	error2 "order-service/common/error"
	constants "order-service/constants/error"
	error3 "order-service/constants/error/order"
	"order-service/domain/dto"
	"order-service/domain/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderRepository struct {
	db *gorm.DB
}

type OrderRepositoryInterface interface {
	FindAllWithPagination(ctx context.Context, param *dto.OrderRequestParam) ([]models.Order, int64, error)
	FindByUserID(ctx context.Context, userID string) ([]models.Order, error)
	FindByUUID(ctx context.Context, uuid string) (*models.Order, error)
	Create(ctx context.Context, db *gorm.DB, req *models.Order) (*models.Order, error)
	Update(ctx context.Context, db *gorm.DB, req *models.Order, userID uuid.UUID) error
}

func NewOrderRepository(db *gorm.DB) OrderRepositoryInterface {
	return &OrderRepository{db: db}
}

func (o *OrderRepository) FindAllWithPagination(ctx context.Context, param *dto.OrderRequestParam) ([]models.Order, int64, error) {
	var (
		orders []models.Order
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
	err := o.db.WithContext(ctx).Limit(limit).Offset(offset).Order(sort).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}
	err = o.db.WithContext(ctx).Model(&models.Order{}).Count(&total).Error
	if err != nil {
		return nil, 0, error2.WrapError(constants.ErrSQLError)
	}
	return orders, total, nil

}

func (o *OrderRepository) FindByUserID(ctx context.Context, userID string) ([]models.Order, error) {
	var orders []models.Order
	err := o.db.WithContext(ctx).Where("user_id = ?", userID).Find(&orders).Error
	if err != nil {
		return nil, error2.WrapError(constants.ErrSQLError)
	}
	return orders, nil
}

func (o *OrderRepository) FindByUUID(ctx context.Context, uuid string) (*models.Order, error) {
	var order models.Order
	err := o.db.WithContext(ctx).Where("uuid = ?", uuid).First(&order).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, error2.WrapError(error3.ErrOrderNotFound)
		}
		return nil, error2.WrapError(constants.ErrSQLError)
	}
	return &order, nil
}

func (o *OrderRepository) incrementCode(ctx context.Context, tx *gorm.DB) (string, error) {
	var (
		seq models.OrderSequence
		now = time.Now().Format(time.DateOnly)
	)

	err := tx.WithContext(ctx).
		Where("date = ?", now).
		Clauses(clause.Locking{Strength: "UPDATE"}).
		First(&seq).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			seq = models.OrderSequence{
				Date:       now,
				LastNumber: 1,
			}
			if err := tx.WithContext(ctx).Create(&seq).Error; err != nil {
				return "", error2.WrapError(constants.ErrSQLError)
			}
		} else {
			return "", error2.WrapError(constants.ErrSQLError)
		}
	} else {
		seq.LastNumber++
		if err := tx.WithContext(ctx).Save(&seq).Error; err != nil {
			return "", error2.WrapError(constants.ErrSQLError)
		}
	}

	result := fmt.Sprintf("ORD-%05d-%s", seq.LastNumber, now)
	return result, nil
}

func (o *OrderRepository) Create(ctx context.Context, db *gorm.DB, req *models.Order) (*models.Order, error) {
	var order models.Order

	err := db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		code, err := o.incrementCode(ctx, tx)
		if err != nil {
			return err
		}

		order = models.Order{
			UUID:   uuid.New(),
			Code:   code,
			UserID: req.UserID,
			Amount: req.Amount,
			Date:   req.Date,
			Status: req.Status,
			IsPaid: req.IsPaid,
		}

		if err := tx.Create(&order).Error; err != nil {
			return error2.WrapError(constants.ErrSQLError)
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &order, nil
}

func (o *OrderRepository) Update(ctx context.Context, db *gorm.DB, req *models.Order, userID uuid.UUID) error {
	err := db.WithContext(ctx).Where("uuid = ?", userID).Updates(req).Error
	if err != nil {
		return error2.WrapError(constants.ErrSQLError)
	}
	return nil
}
