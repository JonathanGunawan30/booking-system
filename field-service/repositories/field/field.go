package repositories

import (
	"context"
	"errors"
	errWrap "field-service/common/error"
	errConstant "field-service/constants/error"
	errField "field-service/constants/error/field"
	"field-service/domain/dto"
	"field-service/domain/models"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type FieldRepositories struct {
	db *gorm.DB
}

type FieldRepositoriesInterface interface {
	FindAllWithPagination(ctx context.Context, param *dto.FieldRequestParam) ([]models.Field, int64, error)
	FindAllWithoutPagination(ctx context.Context) ([]models.Field, error)
	FindByUUID(ctx context.Context, uuid string) (*models.Field, error)
	Create(ctx context.Context, req *models.Field) (*models.Field, error)
	Update(ctx context.Context, uuid string, req *models.Field) (*models.Field, error)
	Delete(ctx context.Context, uuid string) error
}

func NewFieldRepositories(db *gorm.DB) FieldRepositoriesInterface {
	return &FieldRepositories{db: db}
}

func (f *FieldRepositories) FindAllWithPagination(ctx context.Context, param *dto.FieldRequestParam) ([]models.Field, int64, error) {
	var (
		fields []models.Field
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
	err := f.db.WithContext(ctx).Limit(limit).Offset(offset).Order(sort).Find(&fields).Error
	if err != nil {
		return nil, 0, err
	}

	err = f.db.WithContext(ctx).Model(&models.Field{}).Count(&total).Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return fields, total, nil
}

func (f *FieldRepositories) FindAllWithoutPagination(ctx context.Context) ([]models.Field, error) {
	var fields []models.Field
	err := f.db.WithContext(ctx).Find(&fields).Error
	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return fields, nil
}

func (f *FieldRepositories) FindByUUID(ctx context.Context, uuid string) (*models.Field, error) {
	var field models.Field
	if err := f.db.WithContext(ctx).Where("uuid = ?", uuid).First(&field).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errField.ErrFieldNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &field, nil
}

func (f *FieldRepositories) Create(ctx context.Context, req *models.Field) (*models.Field, error) {
	field := models.Field{
		UUID:         uuid.New(),
		Code:         req.Code,
		Name:         req.Name,
		Image:        req.Image,
		PricePerHour: req.PricePerHour,
	}
	result := f.db.WithContext(ctx).Create(&field)
	if result.Error != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	if err := f.db.WithContext(ctx).First(&field, field.ID).Error; err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &field, nil
}

func (f *FieldRepositories) Update(ctx context.Context, uuid string, req *models.Field) (*models.Field, error) {
	field := models.Field{
		Code:         req.Code,
		Name:         req.Name,
		Image:        req.Image,
		PricePerHour: req.PricePerHour,
	}
	result := f.db.WithContext(ctx).Where("uuid = ?", uuid).Updates(&field)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errField.ErrFieldNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	var updatedField models.Field
	if err := f.db.WithContext(ctx).Where("uuid = ?", uuid).First(&updatedField).Error; err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return &updatedField, nil
}

func (f *FieldRepositories) Delete(ctx context.Context, uuid string) error {
	result := f.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&models.Field{})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errWrap.WrapError(errField.ErrFieldNotFound)
		}
		return errWrap.WrapError(errConstant.ErrSQLError)
	}
	return nil
}
