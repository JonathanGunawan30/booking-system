package fieldschedule

import (
	"context"
	"encoding/json"
	"field-service/common/util"
	"field-service/domain/dto"
	fieldService "field-service/services/field"
	fieldScheduleService "field-service/services/fieldschedule"
	timeService "field-service/services/time"
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

type MockFieldScheduleService struct {
	mock.Mock
}

func (m *MockFieldScheduleService) GetAllWithPagination(ctx context.Context, param *dto.FieldScheduleRequestParam) (*util.PaginationResult, error) {
	args := m.Called(ctx, param)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*util.PaginationResult), args.Error(1)
}

func (m *MockFieldScheduleService) GetAllByFieldIDAndDate(ctx context.Context, uuid string, date string) ([]dto.FieldScheduleBookingResponse, error) {
	args := m.Called(ctx, uuid, date)
	return args.Get(0).([]dto.FieldScheduleBookingResponse), args.Error(1)
}

func (m *MockFieldScheduleService) GetByUUID(ctx context.Context, uuid string) (*dto.FieldScheduleResponse, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.FieldScheduleResponse), args.Error(1)
}

func (m *MockFieldScheduleService) GenerateFieldScheduleForOneMonth(ctx context.Context, req *dto.GenerateFieldScheduleForOneMonthRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockFieldScheduleService) Create(ctx context.Context, req *dto.FieldScheduleRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockFieldScheduleService) Update(ctx context.Context, uuid string, req *dto.UpdateFieldScheduleRequest) (*dto.FieldScheduleResponse, error) {
	args := m.Called(ctx, uuid, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.FieldScheduleResponse), args.Error(1)
}

func (m *MockFieldScheduleService) UpdateStatus(ctx context.Context, req *dto.UpdateStatusFieldScheduleRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

func (m *MockFieldScheduleService) Delete(ctx context.Context, uuid string) error {
	args := m.Called(ctx, uuid)
	return args.Error(0)
}

func TestFieldScheduleController_GetByUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockServiceRegistry := new(MockServiceRegistry)
	mockFieldScheduleService := new(MockFieldScheduleService)
	controller := NewFieldScheduleController(mockServiceRegistry)

	id := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		mockServiceRegistry.On("GetFieldSchedule").Return(mockFieldScheduleService)
		mockFieldScheduleService.On("GetByUUID", mock.Anything, id).Return(&dto.FieldScheduleResponse{UUID: uuid.MustParse(id)}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "uuid", Value: id}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/field-schedule/"+id, nil)

		controller.GetByUUID(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "success", response["status"])
	})

	t.Run("invalid uuid", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "uuid", Value: "invalid"}}

		controller.GetByUUID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestFieldScheduleController_GetAllByFieldIDAndDate(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockServiceRegistry := new(MockServiceRegistry)
	mockFieldScheduleService := new(MockFieldScheduleService)
	controller := NewFieldScheduleController(mockServiceRegistry)

	id := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		mockServiceRegistry.On("GetFieldSchedule").Return(mockFieldScheduleService)
		mockFieldScheduleService.On("GetAllByFieldIDAndDate", mock.Anything, id, "2024-01-01").Return([]dto.FieldScheduleBookingResponse{}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "uuid", Value: id}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/field-schedule/field/"+id+"?date=2024-01-01", nil)

		controller.GetAllByFieldIDAndDate(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("missing date", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "uuid", Value: id}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/field-schedule/field/"+id, nil)

		controller.GetAllByFieldIDAndDate(c)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
	})
}
