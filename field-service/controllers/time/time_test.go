package time

import (
	"context"
	"encoding/json"
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

type MockTimeService struct {
	mock.Mock
}

func (m *MockTimeService) GetAll(ctx context.Context) ([]dto.TimeResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).([]dto.TimeResponse), args.Error(1)
}

func (m *MockTimeService) GetByUUID(ctx context.Context, uuid string) (*dto.TimeResponse, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TimeResponse), args.Error(1)
}

func (m *MockTimeService) Create(ctx context.Context, req *dto.TimeRequest) (*dto.TimeResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.TimeResponse), args.Error(1)
}

func TestTimeController_GetAll(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockServiceRegistry := new(MockServiceRegistry)
	mockTimeService := new(MockTimeService)
	controller := NewTimeController(mockServiceRegistry)

	t.Run("success", func(t *testing.T) {
		mockServiceRegistry.On("GetTime").Return(mockTimeService)
		mockTimeService.On("GetAll", mock.Anything).Return([]dto.TimeResponse{}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/time", nil)

		controller.GetAll(c)

		assert.Equal(t, http.StatusOK, w.Code)
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		assert.Equal(t, "success", response["status"])
	})
}

func TestTimeController_GetByUUID(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockServiceRegistry := new(MockServiceRegistry)
	mockTimeService := new(MockTimeService)
	controller := NewTimeController(mockServiceRegistry)

	id := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		mockServiceRegistry.On("GetTime").Return(mockTimeService)
		mockTimeService.On("GetByUUID", mock.Anything, id).Return(&dto.TimeResponse{UUID: uuid.MustParse(id)}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "uuid", Value: id}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/time/"+id, nil)

		controller.GetByUUID(c)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("invalid uuid", func(t *testing.T) {
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "uuid", Value: "invalid"}}

		controller.GetByUUID(c)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}
