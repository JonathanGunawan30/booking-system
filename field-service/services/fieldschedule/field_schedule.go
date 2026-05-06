package services

import (
	"context"
	"errors"
	"field-service/common/util"
	"field-service/constants"
	errField "field-service/constants/error/field"
	errFieldSchedule "field-service/constants/error/fieldschedule"
	"field-service/domain/dto"
	"field-service/domain/models"
	"field-service/repositories"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type FieldScheduleService struct {
	repository repositories.RepositoryRegistryInterface
}

type FieldScheduleServiceInterface interface {
	GetAllWithPagination(ctx context.Context, param *dto.FieldScheduleRequestParam) (*util.PaginationResult, error)
	GetAllByFieldIDAndDate(ctx context.Context, uuid string, date string) ([]dto.FieldScheduleBookingResponse, error)
	GetByUUID(ctx context.Context, uuid string) (*dto.FieldScheduleResponse, error)
	GenerateFieldScheduleForOneMonth(ctx context.Context, req *dto.GenerateFieldScheduleForOneMonthRequest) error
	Create(ctx context.Context, req *dto.FieldScheduleRequest) error
	Update(ctx context.Context, uuid string, req *dto.UpdateFieldScheduleRequest) (*dto.FieldScheduleResponse, error)
	UpdateStatus(ctx context.Context, req *dto.UpdateStatusFieldScheduleRequest) error
	Delete(ctx context.Context, uuid string) error
}

func NewFieldScheduleService(repository repositories.RepositoryRegistryInterface) FieldScheduleServiceInterface {
	return &FieldScheduleService{repository: repository}
}

func (f *FieldScheduleService) GetAllWithPagination(ctx context.Context, param *dto.FieldScheduleRequestParam) (*util.PaginationResult, error) {
	fieldSchedules, total, err := f.repository.GetFieldSchedule().FindAllWithPagination(ctx, param)
	if err != nil {
		return nil, err
	}

	fieldScheduleResults := make([]dto.FieldScheduleResponse, 0, len(fieldSchedules))
	for _, schedule := range fieldSchedules {
		fieldScheduleResults = append(fieldScheduleResults, dto.FieldScheduleResponse{
			UUID:         schedule.UUID,
			FieldName:    schedule.Field.Name,
			Date:         schedule.Date.Format(time.DateOnly),
			PricePerHour: schedule.Field.PricePerHour,
			Status:       schedule.Status.GetStatusString(),
			Time:         fmt.Sprintf("%s - %s", schedule.Time.StartTime, schedule.Time.EndTime),
			CreatedAt:    schedule.CreatedAt,
			UpdatedAt:    schedule.UpdatedAt,
		})
	}

	pagination := &util.PaginationParam{
		Count: total,
		Page:  param.Page,
		Limit: param.Limit,
		Data:  fieldScheduleResults,
	}
	return new(util.GeneratePagination(*pagination)), nil
}

func (f *FieldScheduleService) GetAllByFieldIDAndDate(ctx context.Context, uuid, date string) ([]dto.FieldScheduleBookingResponse, error) {
	field, err := f.repository.GetField().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	fieldSchedules, err := f.repository.GetFieldSchedule().FindAllByFieldIDAndDate(ctx, int(field.ID), date)
	if err != nil {
		if errors.Is(err, errField.ErrFieldNotFound) {
			return []dto.FieldScheduleBookingResponse{}, nil
		}
		return nil, err
	}

	fieldSchedulesResults := make([]dto.FieldScheduleBookingResponse, 0, len(fieldSchedules))
	for _, schedule := range fieldSchedules {
		fieldSchedulesResults = append(fieldSchedulesResults, dto.FieldScheduleBookingResponse{
			UUID:         schedule.UUID,
			Date:         schedule.Date.Format(time.DateOnly),
			Time:         fmt.Sprintf("%s - %s", schedule.Time.StartTime, schedule.Time.EndTime),
			Status:       schedule.Status.GetStatusString(),
			PricePerHour: util.RupiahFormat(util.Float64(float64(schedule.Field.PricePerHour))),
		})
	}

	return fieldSchedulesResults, nil

}

func (f *FieldScheduleService) GetByUUID(ctx context.Context, uuid string) (*dto.FieldScheduleResponse, error) {
	fieldSchedule, err := f.repository.GetFieldSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	response := dto.FieldScheduleResponse{
		UUID:         fieldSchedule.UUID,
		FieldName:    fieldSchedule.Field.Name,
		Date:         fieldSchedule.Date.Format(time.DateOnly),
		PricePerHour: fieldSchedule.Field.PricePerHour,
		Status:       fieldSchedule.Status.GetStatusString(),
		Time:         fmt.Sprintf("%s - %s", fieldSchedule.Time.StartTime, fieldSchedule.Time.EndTime),
		CreatedAt:    fieldSchedule.CreatedAt,
		UpdatedAt:    fieldSchedule.UpdatedAt,
	}

	return &response, nil
}

func (f *FieldScheduleService) GenerateFieldScheduleForOneMonth(ctx context.Context, req *dto.GenerateFieldScheduleForOneMonthRequest) error {
	if req.Year < time.Now().Year() {
		return errFieldSchedule.ErrFieldScheduleYearPast
	}

	if req.Year == time.Now().Year() && req.Month < int(time.Now().Month()) {
		return errFieldSchedule.ErrFieldScheduleMonthPast
	}

	field, err := f.repository.GetField().FindByUUID(ctx, req.FieldID)
	if err != nil {
		return err
	}

	times, err := f.repository.GetTime().FindAll(ctx)
	if err != nil {
		return err
	}

	firstDay := time.Date(req.Year, time.Month(req.Month), 1, 0, 0, 0, 0, time.UTC)
	lastDay := firstDay.AddDate(0, 1, -1)
	numberOfDays := lastDay.Day()

	existingSchedules, err := f.repository.GetFieldSchedule().FindAllByFieldIDAndDate(ctx, int(field.ID), firstDay.Format(time.DateOnly))
	if err != nil {
		return err
	}

	existingMap := make(map[string]bool)
	for _, s := range existingSchedules {
		key := fmt.Sprintf("%s-%d", s.Date.Format(time.DateOnly), s.TimeID)
		existingMap[key] = true
	}

	fieldSchedules := make([]models.FieldSchedule, 0, numberOfDays*len(times))
	for i := 0; i < numberOfDays; i++ {
		date := firstDay.AddDate(0, 0, i)
		for _, item := range times {
			key := fmt.Sprintf("%s-%d", date.Format(time.DateOnly), item.ID)
			if existingMap[key] {
				return errFieldSchedule.ErrFieldScheduleExist
			}

			fieldSchedules = append(fieldSchedules, models.FieldSchedule{
				UUID:    uuid.New(),
				FieldID: field.ID,
				TimeID:  item.ID,
				Date:    date,
				Status:  constants.Available,
			})
		}
	}

	if err = f.repository.GetFieldSchedule().Create(ctx, fieldSchedules); err != nil {
		return err
	}

	return nil
}

func (f *FieldScheduleService) Create(ctx context.Context, req *dto.FieldScheduleRequest) error {
	field, err := f.repository.GetField().FindByUUID(ctx, req.FieldID)
	if err != nil {
		return err
	}

	fieldSchedules := make([]models.FieldSchedule, 0, len(req.TimeIDs))
	dateParsed, err := time.Parse(time.DateOnly, req.Date)
	if err != nil {
		return err
	}

	for _, timeID := range req.TimeIDs {
		scheduleTime, err := f.repository.GetTime().FindByUUID(ctx, timeID)
		if err != nil {
			return err
		}

		schedule, err := f.repository.GetFieldSchedule().FindByDateAndTimeID(ctx, req.Date, int(scheduleTime.ID), int(field.ID))
		if err != nil {
			if !errors.Is(err, errFieldSchedule.ErrFieldScheduleNotFound) {
				return err
			}
		}

		if schedule != nil {
			return errFieldSchedule.ErrFieldScheduleExist
		}

		fieldSchedules = append(fieldSchedules, models.FieldSchedule{
			UUID:    uuid.New(),
			FieldID: field.ID,
			TimeID:  scheduleTime.ID,
			Date:    dateParsed,
			Status:  constants.Available,
		})
	}

	err = f.repository.GetFieldSchedule().Create(ctx, fieldSchedules)
	if err != nil {
		return err
	}
	return nil
}

func (f *FieldScheduleService) Update(ctx context.Context, uuid string, req *dto.UpdateFieldScheduleRequest) (*dto.FieldScheduleResponse, error) {
	fieldSchedule, err := f.repository.GetFieldSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	scheduleTime, err := f.repository.GetTime().FindByUUID(ctx, req.TimeID)
	if err != nil {
		return nil, err
	}

	isTimeExists, err := f.repository.GetFieldSchedule().FindByDateAndTimeID(ctx, req.Date, int(scheduleTime.ID), int(fieldSchedule.FieldID))
	if err != nil {
		return nil, err
	}

	if isTimeExists != nil && isTimeExists.UUID != fieldSchedule.UUID {
		return nil, errFieldSchedule.ErrFieldScheduleExist
	}

	dateParsed, err := time.Parse(time.DateOnly, req.Date)
	if err != nil {
		return nil, err
	}

	fieldUpdated, err := f.repository.GetFieldSchedule().Update(ctx, uuid, &models.FieldSchedule{
		Date:   dateParsed,
		TimeID: scheduleTime.ID,
	})

	if err != nil {
		return nil, err
	}

	return &dto.FieldScheduleResponse{
		UUID:         fieldUpdated.UUID,
		FieldName:    fieldSchedule.Field.Name,
		Date:         fieldUpdated.Date.Format(time.DateOnly),
		PricePerHour: fieldSchedule.Field.PricePerHour,
		Status:       fieldUpdated.Status.GetStatusString(),
		Time:         fmt.Sprintf("%s - %s", fieldUpdated.Time.StartTime, fieldUpdated.Time.EndTime),
		CreatedAt:    fieldUpdated.CreatedAt,
		UpdatedAt:    fieldUpdated.UpdatedAt,
	}, nil
}

func (f *FieldScheduleService) UpdateStatus(ctx context.Context, req *dto.UpdateStatusFieldScheduleRequest) error {
	for _, item := range req.FieldScheduleIDs {
		logrus.Infof("[FieldScheduleService] Updating status for UUID: %s to %d", item, req.Status)
		_, err := f.repository.GetFieldSchedule().FindByUUID(ctx, item)
		if err != nil {
			logrus.Errorf("[FieldScheduleService] Error finding UUID %s: %v", item, err)
			return err
		}

		err = f.repository.GetFieldSchedule().UpdateStatus(ctx, item, constants.FieldScheduleStatus(req.Status))
		if err != nil {
			logrus.Errorf("[FieldScheduleService] Error updating status for UUID %s: %v", item, err)
			return err
		}
		logrus.Infof("[FieldScheduleService] Successfully updated status for UUID: %s", item)
	}
	return nil
}

func (f *FieldScheduleService) Delete(ctx context.Context, uuid string) error {
	_, err := f.repository.GetFieldSchedule().FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	err = f.repository.GetFieldSchedule().Delete(ctx, uuid)
	if err != nil {
		return err
	}
	return nil
}
