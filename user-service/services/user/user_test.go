package services

import (
	"context"
	"testing"
	constants "user-service/constants"
	errConstants "user-service/constants/error"
	"user-service/domain/dto"
	"user-service/domain/models"
	repoUser "user-service/repositories/user"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Register(ctx context.Context, request *dto.RegisterRequest) (*models.User, error) {
	args := m.Called(ctx, request)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) Update(ctx context.Context, request *dto.UpdateRequest, uuid string) (*models.User, error) {
	args := m.Called(ctx, request, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	args := m.Called(ctx, username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	args := m.Called(ctx, email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) FindByUUID(ctx context.Context, uuid string) (*models.User, error) {
	args := m.Called(ctx, uuid)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

type MockRepositoryRegistry struct {
	mock.Mock
}

func (m *MockRepositoryRegistry) GetUser() repoUser.UserRepositoryInterface {
	args := m.Called()
	return args.Get(0).(repoUser.UserRepositoryInterface)
}

func TestLogin(t *testing.T) {
	ctx := context.Background()
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	user := &models.User{
		Username: "testuser",
		Password: string(hashedPassword),
		Role: models.Role{
			Name: "Customer",
		},
	}

	t.Run("Success Login", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockRegistry := new(MockRepositoryRegistry)
		userService := NewUserService(mockRegistry)

		mockRegistry.On("GetUser").Return(mockUserRepo)
		mockUserRepo.On("FindByUsername", ctx, "testuser").Return(user, nil)

		req := &dto.LoginRequest{
			Username: "testuser",
			Password: password,
		}

		resp, err := userService.Login(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "testuser", resp.User.Username)
	})

	t.Run("User Not Found", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockRegistry := new(MockRepositoryRegistry)
		userService := NewUserService(mockRegistry)

		mockRegistry.On("GetUser").Return(mockUserRepo)
		mockUserRepo.On("FindByUsername", ctx, "unknown").Return(nil, errConstants.ErrUserNotFound)

		req := &dto.LoginRequest{
			Username: "unknown",
			Password: password,
		}

		resp, err := userService.Login(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, errConstants.ErrUsernameOrPasswordIncorrect, err)
	})
}

func TestRegister(t *testing.T) {
	ctx := context.Background()
	req := &dto.RegisterRequest{
		Name:     "New User",
		Username: "newuser",
		Email:    "new@mail.com",
		Password: "password123",
	}

	t.Run("Success Register", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockRegistry := new(MockRepositoryRegistry)
		userService := NewUserService(mockRegistry)

		mockRegistry.On("GetUser").Return(mockUserRepo)
		mockUserRepo.On("FindByUsername", ctx, req.Username).Return(nil, errConstants.ErrUserNotFound)
		mockUserRepo.On("FindByEmail", ctx, req.Email).Return(nil, errConstants.ErrUserNotFound)
		
		mockUserRepo.On("Register", ctx, mock.Anything).Return(&models.User{
			UUID:     uuid.New(),
			Username: req.Username,
			Email:    req.Email,
		}, nil)

		resp, err := userService.Register(ctx, req)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, req.Username, resp.User.Username)
	})

	t.Run("Username Exists", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockRegistry := new(MockRepositoryRegistry)
		userService := NewUserService(mockRegistry)

		mockRegistry.On("GetUser").Return(mockUserRepo)
		mockUserRepo.On("FindByUsername", ctx, req.Username).Return(&models.User{Username: req.Username}, nil)

		resp, err := userService.Register(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, errConstants.ErrUsernameExists, err)
	})
}

func TestGetUserByUUID(t *testing.T) {
	ctx := context.Background()
	userUUID := uuid.New().String()

	t.Run("Success Get User", func(t *testing.T) {
		mockUserRepo := new(MockUserRepository)
		mockRegistry := new(MockRepositoryRegistry)
		userService := NewUserService(mockRegistry)

		mockRegistry.On("GetUser").Return(mockUserRepo)
		mockUserRepo.On("FindByUUID", ctx, userUUID).Return(&models.User{
			UUID:     uuid.MustParse(userUUID),
			Username: "testuser",
			Role:     models.Role{Name: "Customer"},
		}, nil)

		resp, err := userService.GetUserByUUID(ctx, userUUID)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "testuser", resp.Username)
	})
}

func TestGetUserLogin(t *testing.T) {
	t.Run("Success Get User Login", func(t *testing.T) {
		mockRegistry := new(MockRepositoryRegistry)
		userService := NewUserService(mockRegistry)

		userData := &dto.UserResponse{
			Username: "testuser",
			Role:     "customer",
		}
		ctx := context.WithValue(context.Background(), constants.UserLogin, userData)

		resp, err := userService.GetUserLogin(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, resp)
		assert.Equal(t, "testuser", resp.Username)
	})
}
