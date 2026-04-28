package services

import (
	"context"
	"field-service/domain/dto"
	"field-service/domain/models"
	fieldRepo "field-service/repositories/field"
	fieldScheduleRepo "field-service/repositories/fieldschedule"
	timeRepo "field-service/repositories/time"
	"testing"

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

func TestTimeService_GetAll(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		mockRepoRegistry := new(MockRepositoryRegistry)
		mockTimeRepo := new(MockTimeRepository)
		service := NewTimeService(mockRepoRegistry)

		times := []models.Time{
			{UUID: uuid.New(), StartTime: "08:00", EndTime: "09:00"},
		}

		mockRepoRegistry.On("GetTime").Return(mockTimeRepo)
		mockTimeRepo.On("FindAll", ctx).Return(times, nil)

		result, err := service.GetAll(ctx)

		assert.NoError(t, err)
		assert.Len(t, result, 1)
		assert.Equal(t, "08:00", result[0].StartTime)
	})

	t.Run("error", func(t *testing.T) {
		mockRepoRegistry := new(MockRepositoryRegistry)
		mockTimeRepo := new(MockTimeRepository)
		service := NewTimeService(mockRepoRegistry)

		mockRepoRegistry.On("GetTime").Return(mockTimeRepo)
		mockTimeRepo.On("FindAll", ctx).Return([]models.Time{}, assert.AnError)

		result, err := service.GetAll(ctx)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestTimeService_Create(t *testing.T) {
	ctx := context.Background()
	req := &dto.TimeRequest{StartTime: "08:00", EndTime: "09:00"}

	t.Run("success", func(t *testing.T) {
		mockRepoRegistry := new(MockRepositoryRegistry)
		mockTimeRepo := new(MockTimeRepository)
		service := NewTimeService(mockRepoRegistry)

		timeModel := &models.Time{UUID: uuid.New(), StartTime: "08:00", EndTime: "09:00"}

		mockRepoRegistry.On("GetTime").Return(mockTimeRepo)
		mockTimeRepo.On("Create", ctx, mock.Anything).Return(timeModel, nil)

		result, err := service.Create(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "08:00", result.StartTime)
	})
}
