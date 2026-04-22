package repositories

import (
	"context"
	"errors"
	errWrap "field-service/common/error"
	"field-service/constants"
	errConstant "field-service/constants/error"
	errFieldSchedule "field-service/constants/error/fieldschedule"

	"field-service/domain/dto"
	"field-service/domain/models"
	"fmt"

	"gorm.io/gorm"
)

type FieldScheduleRepositories struct {
	db *gorm.DB
}

type FieldScheduleRepositoriesInterface interface {
	FindAllWithPagination(ctx context.Context, param *dto.FieldScheduleRequestParam) ([]models.FieldSchedule, int64, error)
	FindAllByFieldIDAndDate(ctx context.Context, id int, date string) ([]models.FieldSchedule, error)
	FindByUUID(ctx context.Context, uuid string) (*models.FieldSchedule, error)
	FindByDateAndTimeID(ctx context.Context, date string, timeID, fieldID int) (*models.FieldSchedule, error)
	Create(ctx context.Context, req []models.FieldSchedule) error
	Update(ctx context.Context, uuid string, req *models.FieldSchedule) (*models.FieldSchedule, error)
	UpdateStatus(ctx context.Context, uuid string, status constants.FieldScheduleStatus) error
	Delete(ctx context.Context, uuid string) error
}

func NewFieldScheduleRepositories(db *gorm.DB) FieldScheduleRepositoriesInterface {
	return &FieldScheduleRepositories{db: db}
}

func (f *FieldScheduleRepositories) FindAllWithPagination(ctx context.Context, param *dto.FieldScheduleRequestParam) ([]models.FieldSchedule, int64, error) {
	var (
		fieldSchedule []models.FieldSchedule
		sort          string
		total         int64
	)

	if param.SortColumn != nil {
		sort = fmt.Sprintf("%s %s", *param.SortColumn, *param.SortOrder)
	} else {
		sort = "created_at desc"
	}

	limit := param.Limit
	offset := (param.Page - 1) * limit
	err := f.db.WithContext(ctx).Preload("Field").Preload("Time").Limit(limit).Offset(offset).Order(sort).Find(&fieldSchedule).Error
	if err != nil {
		return nil, 0, err
	}

	err = f.db.WithContext(ctx).Model(&models.FieldSchedule{}).Count(&total).Error
	if err != nil {
		return nil, 0, errWrap.WrapError(errConstant.ErrSQLError)
	}

	return fieldSchedule, total, nil
}

func (f *FieldScheduleRepositories) FindAllByFieldIDAndDate(ctx context.Context, fieldID int, date string) ([]models.FieldSchedule, error) {
	var fieldSchedules []models.FieldSchedule
	err := f.db.WithContext(ctx).
		Preload("Field").
		Preload("Time").
		Joins("LEFT JOIN times ON field_schedules.time_id = times.id").
		Where("field_schedules.field_id = ? AND field_schedules.date = ?", fieldID, date).
		Order("times.start_time asc").
		Find(&fieldSchedules).Error

	if err != nil {
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return fieldSchedules, nil
}

func (f *FieldScheduleRepositories) FindByUUID(ctx context.Context, uuid string) (*models.FieldSchedule, error) {
	var fieldSchedule models.FieldSchedule
	if err := f.db.WithContext(ctx).Preload("Field").Preload("Time").Where("uuid = ?", uuid).First(&fieldSchedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errFieldSchedule.ErrFieldScheduleNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &fieldSchedule, nil
}

func (f *FieldScheduleRepositories) FindByDateAndTimeID(ctx context.Context, date string, timeID, fieldID int) (*models.FieldSchedule, error) {
	var fieldSchedule models.FieldSchedule
	if err := f.db.WithContext(ctx).Where("date = ? AND time_id = ? AND field_id = ?", date, timeID, fieldID).First(&fieldSchedule).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return &fieldSchedule, nil

}

func (f *FieldScheduleRepositories) Create(ctx context.Context, req []models.FieldSchedule) error {
	if err := f.db.WithContext(ctx).Create(&req).Error; err != nil {
		return errWrap.WrapError(errConstant.ErrSQLError)
	}
	return nil
}

func (f *FieldScheduleRepositories) Update(ctx context.Context, uuid string, req *models.FieldSchedule) (*models.FieldSchedule, error) {
	if err := f.db.WithContext(ctx).Where("uuid = ?", uuid).Updates(&req).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errWrap.WrapError(errFieldSchedule.ErrFieldScheduleNotFound)
		}
		return nil, errWrap.WrapError(errConstant.ErrSQLError)
	}
	return f.FindByUUID(ctx, uuid)
}

func (f *FieldScheduleRepositories) UpdateStatus(ctx context.Context, uuid string, status constants.FieldScheduleStatus) error {
	if err := f.db.WithContext(ctx).Model(&models.FieldSchedule{}).Where("uuid = ?", uuid).Update("status", status).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errWrap.WrapError(errFieldSchedule.ErrFieldScheduleNotFound)
		}
		return errWrap.WrapError(errConstant.ErrSQLError)
	}
	return nil
}

func (f *FieldScheduleRepositories) Delete(ctx context.Context, uuid string) error {
	result := f.db.WithContext(ctx).Where("uuid = ?", uuid).Delete(&models.FieldSchedule{})
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return errWrap.WrapError(errFieldSchedule.ErrFieldScheduleNotFound)
		}
		return errWrap.WrapError(errConstant.ErrSQLError)
	}
	return nil
}
