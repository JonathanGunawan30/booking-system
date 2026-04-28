package controllers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"payment-service/common/util"
	"payment-service/domain/dto"
	payService "payment-service/services/payment"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockPaymentService struct {
	mock.Mock
}

func (m *MockPaymentService) GetAllWithPagination(ctx context.Context, param *dto.PaymentRequestParam) (*util.PaginationResult, error) {
	args := m.Called(ctx, param)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*util.PaginationResult), args.Error(1)
}

func (m *MockPaymentService) GetByUUID(ctx context.Context, uuid string) (*dto.PaymentResponse, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PaymentResponse), args.Error(1)
}

func (m *MockPaymentService) Create(ctx context.Context, request *dto.PaymentRequest) (*dto.PaymentResponse, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.PaymentResponse), args.Error(1)
}

func (m *MockPaymentService) WebHook(ctx context.Context, hook *dto.WebHook) error {
	args := m.Called(ctx, hook)
	return args.Error(0)
}

type MockServiceRegistry struct {
	mock.Mock
}

func (m *MockServiceRegistry) GetPayment() payService.PaymentServiceInterface {
	return m.Called().Get(0).(payService.PaymentServiceInterface)
}

func TestPaymentController_Create(t *testing.T) {
	gin.SetMode(gin.TestMode)
	
	t.Run("Success 201", func(t *testing.T) {
		mockServiceRegistry := new(MockServiceRegistry)
		mockPaymentService := new(MockPaymentService)
		mockServiceRegistry.On("GetPayment").Return(mockPaymentService)
		controller := NewPaymentController(mockServiceRegistry)

		reqBody := dto.PaymentRequest{
			OrderID: "550e8400-e29b-41d4-a716-446655440000",
			Amount:  10000,
			ExpiredAt: time.Now().Add(time.Hour),
			CustomerDetail: &dto.CustomerDetail{
				Name: "Test",
				Email: "test@example.com",
				Phone: "08123456789",
			},
			ItemDetails: []dto.ItemDetail{
				{ID: "1", Name: "Item 1", Amount: 10000, Quantity: 1},
			},
		}
		mockPaymentService.On("Create", mock.Anything, mock.Anything).Return(&dto.PaymentResponse{Amount: 10000}, nil)

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest("POST", "/payments", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		
		controller.Create(c)

		assert.Equal(t, http.StatusCreated, w.Code)
		mockPaymentService.AssertExpectations(t)
	})

	t.Run("Internal Server Error 500", func(t *testing.T) {
		mockServiceRegistry := new(MockServiceRegistry)
		mockPaymentService := new(MockPaymentService)
		mockServiceRegistry.On("GetPayment").Return(mockPaymentService)
		controller := NewPaymentController(mockServiceRegistry)

		reqBody := dto.PaymentRequest{
			OrderID: "550e8400-e29b-41d4-a716-446655440000",
			Amount:  10000,
			ExpiredAt: time.Now().Add(time.Hour),
			CustomerDetail: &dto.CustomerDetail{
				Name: "Test",
				Email: "test@example.com",
				Phone: "08123456789",
			},
			ItemDetails: []dto.ItemDetail{
				{ID: "1", Name: "Item 1", Amount: 10000, Quantity: 1},
			},
		}
		mockPaymentService.On("Create", mock.Anything, mock.Anything).Return(nil, errors.New("some error"))

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		
		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest("POST", "/payments", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		
		controller.Create(c)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		mockPaymentService.AssertExpectations(t)
	})
}
