package services

import (
	"context"
	"field-service/constants"
	"field-service/domain/dto"
	"field-service/domain/models"
	fieldRepo "field-service/repositories/field"
	fieldScheduleRepo "field-service/repositories/fieldschedule"
	timeRepo "field-service/repositories/time"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockRepositoryRegistry struct {
	mock.Mock
}

func (m *MockRepositoryRegistry) GetField() fieldRepo.FieldRepositoriesInterface {
	args := m.Called()
	return args.Get(0).(fieldRepo.FieldRepositoriesInterface)
}

func (m *MockRepositoryRegistry) GetFieldSchedule() fieldScheduleRepo.FieldScheduleRepositoriesInterface {
	args := m.Called()
	return args.Get(0).(fieldScheduleRepo.FieldScheduleRepositoriesInterface)
}

func (m *MockRepositoryRegistry) GetTime() timeRepo.TimeRepositoryInterface {
	args := m.Called()
	return args.Get(0).(timeRepo.TimeRepositoryInterface)
}

type MockFieldScheduleRepository struct {
	mock.Mock
}

func (m *MockFieldScheduleRepository) FindAllWithPagination(ctx context.Context, param *dto.FieldScheduleRequestParam) ([]models.FieldSchedule, int64, error) {
	args := m.Called(ctx, param)
	return args.Get(0).([]models.FieldSchedule), args.Get(1).(int64), args.Error(2)
}

func (m *MockFieldScheduleRepository) FindAllByFieldIDAndDate(ctx context.Context, id int, date string) ([]models.FieldSchedule, error) {
	args := m.Called(ctx, id, date)
	return args.Get(0).([]models.FieldSchedule), args.Error(1)
}

func (m *MockFieldScheduleRepository) FindByUUID(ctx context.Context, uuid string) (*models.FieldSchedule, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FieldSchedule), args.Error(1)
}

func (m *MockFieldScheduleRepository) FindByDateAndTimeID(ctx context.Context, date string, timeID, fieldID int) (*models.FieldSchedule, error) {
	args := m.Called(ctx, date, timeID, fieldID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FieldSchedule), args.Error(1)
}

func (m *MockFieldScheduleRepository) Create(ctx context.Context, req []models.FieldSchedule) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockFieldScheduleRepository) Update(ctx context.Context, uuid string, req *models.FieldSchedule) (*models.FieldSchedule, error) {
	args := m.Called(ctx, uuid, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.FieldSchedule), args.Error(1)
}

func (m *MockFieldScheduleRepository) UpdateStatus(ctx context.Context, uuid string, status constants.FieldScheduleStatus) error {
	args := m.Called(ctx, uuid, status)
	return args.Error(0)
}

func (m *MockFieldScheduleRepository) Delete(ctx context.Context, uuid string) error {
	args := m.Called(ctx, uuid)
	return args.Error(0)
}

type MockFieldRepository struct {
	mock.Mock
}

func (m *MockFieldRepository) FindAllWithPagination(ctx context.Context, param *dto.FieldRequestParam) ([]models.Field, int64, error) {
	args := m.Called(ctx, param)
	return args.Get(0).([]models.Field), args.Get(1).(int64), args.Error(2)
}

func (m *MockFieldRepository) FindAllWithoutPagination(ctx context.Context) ([]models.Field, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Field), args.Error(1)
}

func (m *MockFieldRepository) FindByUUID(ctx context.Context, uuid string) (*models.Field, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Field), args.Error(1)
}

func (m *MockFieldRepository) Create(ctx context.Context, req *models.Field) (*models.Field, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*models.Field), args.Error(1)
}

func (m *MockFieldRepository) Update(ctx context.Context, uuid string, req *models.Field) (*models.Field, error) {
	args := m.Called(ctx, uuid, req)
	return args.Get(0).(*models.Field), args.Error(1)
}

func (m *MockFieldRepository) Delete(ctx context.Context, uuid string) error {
	args := m.Called(ctx, uuid)
	return args.Error(0)
}

type MockTimeRepository struct {
	mock.Mock
}

func (m *MockTimeRepository) FindAll(ctx context.Context) ([]models.Time, error) {
	args := m.Called(ctx)
	return args.Get(0).([]models.Time), args.Error(1)
}

func (m *MockTimeRepository) FindByUUID(ctx context.Context, uuid string) (*models.Time, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Time), args.Error(1)
}

func (m *MockTimeRepository) FindByID(ctx context.Context, id int) (*models.Time, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Time), args.Error(1)
}

func (m *MockTimeRepository) Create(ctx context.Context, req *models.Time) (*models.Time, error) {
	args := m.Called(ctx, req)
	return args.Get(0).(*models.Time), args.Error(1)
}

func TestFieldScheduleService_GetByUUID(t *testing.T) {
	ctx := context.Background()
	id := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		mockRepoRegistry := new(MockRepositoryRegistry)
		mockFieldScheduleRepo := new(MockFieldScheduleRepository)
		service := NewFieldScheduleService(mockRepoRegistry)

		schedule := &models.FieldSchedule{
			UUID:   uuid.MustParse(id),
			Date:   time.Now(),
			Status: constants.Available,
			Field:  models.Field{Name: "Field 1"},
			Time:   models.Time{StartTime: "08:00", EndTime: "09:00"},
		}

		mockRepoRegistry.On("GetFieldSchedule").Return(mockFieldScheduleRepo)
		mockFieldScheduleRepo.On("FindByUUID", ctx, id).Return(schedule, nil)

		result, err := service.GetByUUID(ctx, id)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Field 1", result.FieldName)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepoRegistry := new(MockRepositoryRegistry)
		mockFieldScheduleRepo := new(MockFieldScheduleRepository)
		service := NewFieldScheduleService(mockRepoRegistry)

		mockRepoRegistry.On("GetFieldSchedule").Return(mockFieldScheduleRepo)
		mockFieldScheduleRepo.On("FindByUUID", ctx, id).Return((*models.FieldSchedule)(nil), assert.AnError)

		result, err := service.GetByUUID(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestFieldScheduleService_Delete(t *testing.T) {
	ctx := context.Background()
	id := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		mockRepoRegistry := new(MockRepositoryRegistry)
		mockFieldScheduleRepo := new(MockFieldScheduleRepository)
		service := NewFieldScheduleService(mockRepoRegistry)

		mockRepoRegistry.On("GetFieldSchedule").Return(mockFieldScheduleRepo)
		mockFieldScheduleRepo.On("FindByUUID", ctx, id).Return(&models.FieldSchedule{}, nil)
		mockFieldScheduleRepo.On("Delete", ctx, id).Return(nil)

		err := service.Delete(ctx, id)

		assert.NoError(t, err)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepoRegistry := new(MockRepositoryRegistry)
		mockFieldScheduleRepo := new(MockFieldScheduleRepository)
		service := NewFieldScheduleService(mockRepoRegistry)

		mockRepoRegistry.On("GetFieldSchedule").Return(mockFieldScheduleRepo)
		mockFieldScheduleRepo.On("FindByUUID", ctx, id).Return((*models.FieldSchedule)(nil), assert.AnError)

		err := service.Delete(ctx, id)

		assert.Error(t, err)
	})
}
