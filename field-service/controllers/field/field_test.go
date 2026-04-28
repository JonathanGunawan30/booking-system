package fieldcontroller

import (
	"context"
	"encoding/json"
	"field-service/common/util"
	"field-service/domain/dto"
	fieldService "field-service/services/field"
	fieldScheduleService "field-service/services/fieldschedule"
	timeService "field-service/services/time"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockServiceRegistry struct {
	mock.Mock
}

func (m *MockServiceRegistry) GetField() fieldService.FieldServiceInterface {
	args := m.Called()
	return args.Get(0).(fieldService.FieldServiceInterface)
}

func (m *MockServiceRegistry) GetFieldSchedule() fieldScheduleService.FieldScheduleServiceInterface {
	args := m.Called()
	return args.Get(0).(fieldScheduleService.FieldScheduleServiceInterface)
}

func (m *MockServiceRegistry) GetTime() timeService.TimeServiceInterface {
	args := m.Called()
	return args.Get(0).(timeService.TimeServiceInterface)
}

type MockFieldService struct {
	mock.Mock
}

func (m *MockFieldService) GetAllWithPagination(ctx context.Context, param *dto.FieldRequestParam) (*util.PaginationResult, error) {
	args := m.Called(ctx, param)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*util.PaginationResult), args.Error(1)
}

func (m *MockFieldService) GetAllWithoutPagination(ctx context.Context) ([]dto.FieldResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dto.FieldResponse), args.Error(1)
}

func (m *MockFieldService) GetByUUID(ctx context.Context, uuid string) (*dto.FieldResponse, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.FieldResponse), args.Error(1)
}

func (m *MockFieldService) Create(ctx context.Context, req *dto.FieldRequest, images []multipart.FileHeader) (*dto.FieldResponse, error) {
	args := m.Called(ctx, req, images)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.FieldResponse), args.Error(1)
}

func (m *MockFieldService) Update(ctx context.Context, uuid string, req *dto.UpdateFieldRequest, images []multipart.FileHeader) (*dto.FieldResponse, error) {
	args := m.Called(ctx, uuid, req, images)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.FieldResponse), args.Error(1)
}

func (m *MockFieldService) Delete(ctx context.Context, uuid string) error {
	args := m.Called(ctx, uuid)
	return args.Error(0)
}

func TestFieldController_GetByUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockServiceRegistry := new(MockServiceRegistry)
	mockFieldService := new(MockFieldService)
	controller := NewFieldController(mockServiceRegistry)

	t.Run("success", func(t *testing.T) {
		id := uuid.New().String()
		fieldResponse := &dto.FieldResponse{UUID: uuid.MustParse(id), Name: "Field 1"}

		mockServiceRegistry.On("GetField").Return(mockFieldService)
		mockFieldService.On("GetByUUID", mock.Anything, id).Return(fieldResponse, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "uuid", Value: id}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/field/"+id, nil)

		controller.GetByUUID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "success", response["status"])
	})

	t.Run("invalid uuid", func(t *testing.T) {
		id := "invalid-uuid"
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "uuid", Value: id}}

		controller.GetByUUID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		id := uuid.New().String()
		mockServiceRegistry.On("GetField").Return(mockFieldService)
		mockFieldService.On("GetByUUID", mock.Anything, id).Return((*dto.FieldResponse)(nil), assert.AnError)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "uuid", Value: id}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/field/"+id, nil)

		controller.GetByUUID(c)

		assert.NotEqual(t, http.StatusOK, w.Code)
	})
}

func TestFieldController_GetAllWithPagination(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockServiceRegistry := new(MockServiceRegistry)
	mockFieldService := new(MockFieldService)
	controller := NewFieldController(mockServiceRegistry)

	t.Run("success", func(t *testing.T) {
		mockServiceRegistry.On("GetField").Return(mockFieldService)
		mockFieldService.On("GetAllWithPagination", mock.Anything, mock.Anything).Return(&util.PaginationResult{}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/field?page=1&limit=10&sort_column=id&sort_order=asc", nil)

		controller.GetAllWithPagination(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/field?page=0", nil) // page < 1 fails validation

		controller.GetAllWithPagination(c)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
}
