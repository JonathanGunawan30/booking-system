package services

import (
	"context"
	"field-service/domain/dto"
	"field-service/domain/models"
	repositories "field-service/repositories"

	"github.com/google/uuid"
)

type TimeService struct {
	repository repositories.RepositoryRegistryInterface
}

type TimeServiceInterface interface {
	GetAll(ctx context.Context) ([]dto.TimeResponse, error)
	GetByUUID(ctx context.Context, uuid string) (*dto.TimeResponse, error)
	Create(ctx context.Context, req *dto.TimeRequest) (*dto.TimeResponse, error)
}

func NewTimeService(repository repositories.RepositoryRegistryInterface) TimeServiceInterface {
	return &TimeService{repository: repository}
}

func (t *TimeService) GetAll(ctx context.Context) ([]dto.TimeResponse, error) {
	times, err := t.repository.GetTime().FindAll(ctx)
	if err != nil {
		return nil, err
	}

	timeResults := make([]dto.TimeResponse, 0, len(times))
	for _, time := range times {
		timeResults = append(timeResults, dto.TimeResponse{
			UUID:      time.UUID,
			StartTime: time.StartTime,
			EndTime:   time.EndTime,
			CreatedAt: time.CreatedAt,
			UpdatedAt: time.UpdatedAt,
		})
	}

	return timeResults, nil
}

func (t *TimeService) GetByUUID(ctx context.Context, uuid string) (*dto.TimeResponse, error) {
	time, err := t.repository.GetTime().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	response := &dto.TimeResponse{
		UUID:      time.UUID,
		StartTime: time.StartTime,
		EndTime:   time.EndTime,
		CreatedAt: time.CreatedAt,
		UpdatedAt: time.UpdatedAt,
	}

	return response, nil
}

func (t *TimeService) Create(ctx context.Context, req *dto.TimeRequest) (*dto.TimeResponse, error) {
	time, err := t.repository.GetTime().Create(ctx, &models.Time{
		UUID:      uuid.New(),
		StartTime: req.StartTime,
		EndTime:   req.EndTime,
	})

	if err != nil {
		return nil, err
	}

	response := &dto.TimeResponse{
		UUID:      time.UUID,
		StartTime: time.StartTime,
		EndTime:   time.EndTime,
		CreatedAt: time.CreatedAt,
		UpdatedAt: time.UpdatedAt,
	}

	return response, nil
}
