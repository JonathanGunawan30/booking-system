package repositories

import (
	"context"
	"errors"
	errWrap "field-service/common/error"
	errConstant "field-service/constants/error"
	errTime "field-service/constants/error/time"
	"field-service/domain/models"

	"gorm.io/gorm"
)

type TimeRepository struct {
	db *gorm.DB
}

type TimeRepositoryInterface interface {
	FindAll(ctx context.Context) ([]models.Time, error)
	FindByUUID(ctx context.Context, uuid string) (*models.Time, error)
	FindByID(ctx context.Context, id int) (*models.Time, error)
	Create(ctx context.Context, req *models.Time) (*models.Time, error)
}

func NewTimeRepository(db *gorm.DB) TimeRepositoryInterface {
	return &TimeRepository{
		db: db,
	}
}

func (t *TimeRepository) FindAll(ctx context.Context) ([]models.Time, error) {
	var times []models.Time
	err := t.db.WithContext(ctx).Find(&times).Error
	if err != nil {
		return nil, err
	}
	return times, nil
}

func (t *TimeRepository) FindByUUID(ctx context.Context, uuid string) (*models.Time, error) {
	var time models.Time
	if err := t.db.WithContext(ctx).Where("uuid = ?", uuid).First(&time).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errTime.ErrTimeNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &time, nil
}

func (t *TimeRepository) FindByID(ctx context.Context, id int) (*models.Time, error) {
	var time models.Time
	if err := t.db.WithContext(ctx).Where("id = ?", id).First(&time).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errTime.ErrTimeNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &time, nil
}

func (t *TimeRepository) Create(ctx context.Context, req *models.Time) (*models.Time, error) {
	err := t.db.WithContext(ctx).Create(&req).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return req, nil
}
