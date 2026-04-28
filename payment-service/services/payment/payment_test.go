package services

import (
	"context"
	"errors"
	"payment-service/clients/midtrans"
	"payment-service/constants"
	"payment-service/domain/dto"
	"payment-service/domain/models"
	payRepoPackage "payment-service/repositories/payment"
	payHistoryPackage "payment-service/repositories/paymenthistory"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockPaymentRepository struct {
	mock.Mock
}

func (m *MockPaymentRepository) FindAllWithPagination(ctx context.Context, param *dto.PaymentRequestParam) ([]models.Payment, int64, error) {
	args := m.Called(ctx, param)
	return args.Get(0).([]models.Payment), args.Get(1).(int64), args.Error(2)
}

func (m *MockPaymentRepository) FindByUUID(ctx context.Context, uuid string) (*models.Payment, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Payment), args.Error(1)
}

func (m *MockPaymentRepository) FindByOrderID(ctx context.Context, orderID string) (*models.Payment, error) {
	args := m.Called(ctx, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Payment), args.Error(1)
}

func (m *MockPaymentRepository) Create(ctx context.Context, db *gorm.DB, request *dto.PaymentRequest) (*models.Payment, error) {
	args := m.Called(ctx, db, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Payment), args.Error(1)
}

func (m *MockPaymentRepository) Update(ctx context.Context, db *gorm.DB, request *dto.UpdatePaymentRequest, orderID string) (*models.Payment, error) {
	args := m.Called(ctx, db, request, orderID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Payment), args.Error(1)
}

type MockPaymentHistoryRepository struct {
	mock.Mock
}

func (m *MockPaymentHistoryRepository) Create(ctx context.Context, db *gorm.DB, request *dto.PaymentHistoryRequest) error {
	args := m.Called(ctx, db, request)
	return args.Error(0)
}

type MockRepositoryRegistry struct {
	mock.Mock
}

func (m *MockRepositoryRegistry) GetPayment() payRepoPackage.PaymentRepositoryInterface {
	return m.Called().Get(0).(payRepoPackage.PaymentRepositoryInterface)
}

func (m *MockRepositoryRegistry) GetPaymentHistory() payHistoryPackage.PaymentHistoryRepositoryInterface {
	return m.Called().Get(0).(payHistoryPackage.PaymentHistoryRepositoryInterface)
}

func (m *MockRepositoryRegistry) GetTx() *gorm.DB {
	return m.Called().Get(0).(*gorm.DB)
}

type MockMidtransClient struct {
	mock.Mock
}

func (m *MockMidtransClient) CreatePaymentLink(request *dto.PaymentRequest) (*midtrans.MidtransData, error) {
	args := m.Called(request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*midtrans.MidtransData), args.Error(1)
}

func TestPaymentService_GetByUUID(t *testing.T) {
	ctx := context.Background()
	uid := uuid.New()

	t.Run("Success", func(t *testing.T) {
		mockRepoRegistry := new(MockRepositoryRegistry)
		mockPaymentRepo := new(MockPaymentRepository)
		mockRepoRegistry.On("GetPayment").Return(mockPaymentRepo)
		service := &PaymentService{repository: mockRepoRegistry}

		payment := &models.Payment{
			UUID:   uid,
			Amount: 1000,
			Status: func() *constants.PaymentStatus { s := constants.Pending; return &s }(),
		}
		mockPaymentRepo.On("FindByUUID", ctx, uid.String()).Return(payment, nil)

		res, err := service.GetByUUID(ctx, uid.String())

		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, uid, res.UUID)
		mockPaymentRepo.AssertExpectations(t)
	})

	t.Run("Error", func(t *testing.T) {
		mockRepoRegistry := new(MockRepositoryRegistry)
		mockPaymentRepo := new(MockPaymentRepository)
		mockRepoRegistry.On("GetPayment").Return(mockPaymentRepo)
		service := &PaymentService{repository: mockRepoRegistry}

		mockPaymentRepo.On("FindByUUID", ctx, uid.String()).Return(nil, errors.New("error"))

		res, err := service.GetByUUID(ctx, uid.String())

		assert.Error(t, err)
		assert.Nil(t, res)
		mockPaymentRepo.AssertExpectations(t)
	})
}
