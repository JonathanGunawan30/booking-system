package services

import (
	"context"
	"field-service/common/cloudflare"
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

func TestFieldService_GetAllWithPagination(t *testing.T) {
	ctx := context.Background()
	param := &dto.FieldRequestParam{Page: 1, Limit: 10}

	t.Run("success", func(t *testing.T) {
		mockRepoRegistry := new(MockRepositoryRegistry)
		mockFieldRepo := new(MockFieldRepository)
		service := NewFieldService(mockRepoRegistry, cloudflare.R2Client{})

		fields := []models.Field{
			{UUID: uuid.New(), Name: "Field 1"},
		}
		mockRepoRegistry.On("GetField").Return(mockFieldRepo)
		mockFieldRepo.On("FindAllWithPagination", ctx, param).Return(fields, int64(1), nil)

		result, err := service.GetAllWithPagination(ctx, param)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, int64(1), result.TotalData)
	})

	t.Run("repository error", func(t *testing.T) {
		mockRepoRegistry := new(MockRepositoryRegistry)
		mockFieldRepo := new(MockFieldRepository)
		service := NewFieldService(mockRepoRegistry, cloudflare.R2Client{})

		mockRepoRegistry.On("GetField").Return(mockFieldRepo)
		mockFieldRepo.On("FindAllWithPagination", ctx, param).Return([]models.Field{}, int64(0), assert.AnError)

		result, err := service.GetAllWithPagination(ctx, param)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestFieldService_GetByUUID(t *testing.T) {
	ctx := context.Background()
	id := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		mockRepoRegistry := new(MockRepositoryRegistry)
		mockFieldRepo := new(MockFieldRepository)
		service := NewFieldService(mockRepoRegistry, cloudflare.R2Client{})

		field := &models.Field{UUID: uuid.New(), Name: "Field 1"}
		mockRepoRegistry.On("GetField").Return(mockFieldRepo)
		mockFieldRepo.On("FindByUUID", ctx, id).Return(field, nil)

		result, err := service.GetByUUID(ctx, id)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, field.Name, result.Name)
	})

	t.Run("not found", func(t *testing.T) {
		mockRepoRegistry := new(MockRepositoryRegistry)
		mockFieldRepo := new(MockFieldRepository)
		service := NewFieldService(mockRepoRegistry, cloudflare.R2Client{})

		mockRepoRegistry.On("GetField").Return(mockFieldRepo)
		mockFieldRepo.On("FindByUUID", ctx, id).Return((*models.Field)(nil), assert.AnError)

		result, err := service.GetByUUID(ctx, id)

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}
