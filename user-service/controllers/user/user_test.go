package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"user-service/domain/dto"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockUserService struct {
	mock.Mock
}

func (m *MockUserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.LoginResponse), args.Error(1)
}

func (m *MockUserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.RegisterResponse), args.Error(1)
}

func (m *MockUserService) Update(ctx context.Context, req *dto.UpdateRequest, uuid string) (*dto.UserResponse, error) {
	args := m.Called(ctx, req, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserService) GetUserLogin(ctx context.Context) (*dto.UserResponse, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func (m *MockUserService) GetUserByUUID(ctx context.Context, uuid string) (*dto.UserResponse, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*dto.UserResponse), args.Error(1)
}

func TestLoginController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("Success Login", func(t *testing.T) {
		mockService := new(MockUserService)
		ctrl := &UserController{userService: mockService}
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := dto.LoginRequest{Username: "testuser", Password: "password123"}
		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		mockService.On("Login", mock.Anything, mock.Anything).Return(&dto.LoginResponse{}, nil)
		ctrl.Login(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestRegisterController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("Success Register", func(t *testing.T) {
		mockService := new(MockUserService)
		ctrl := &UserController{userService: mockService}
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		reqBody := dto.RegisterRequest{
			Name: "Test User", Username: "testuser", Email: "test@mail.com", 
			Password: "password123", ConfirmPassword: "password123", PhoneNumber: "081234567890",
		}
		body, _ := json.Marshal(reqBody)
		c.Request, _ = http.NewRequest(http.MethodPost, "/", bytes.NewBuffer(body))
		c.Request.Header.Set("Content-Type", "application/json")
		mockService.On("Register", mock.Anything, mock.Anything).Return(&dto.RegisterResponse{}, nil)
		ctrl.Register(c)
		assert.Equal(t, http.StatusCreated, w.Code)
	})
}

func TestGetUserController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	t.Run("Success Get User By UUID", func(t *testing.T) {
		mockService := new(MockUserService)
		ctrl := &UserController{userService: mockService}
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Params = []gin.Param{{Key: "uuid", Value: "uuid"}}
		c.Request, _ = http.NewRequest(http.MethodGet, "/", nil)
		mockService.On("GetUserByUUID", mock.Anything, "uuid").Return(&dto.UserResponse{}, nil)
		ctrl.GetUserByUUID(c)
		assert.Equal(t, http.StatusOK, w.Code)
	})
}
