package repositories

import (
	"context"
	"errors"
	wrap "user-service/common/error"
	errConstant "user-service/constants/error"
	"user-service/domain/dto"
	"user-service/domain/models"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

type UserRepositoryInterface interface {
	Register(ctx context.Context, request *dto.RegisterRequest) (*models.User, error)
	Update(ctx context.Context, request *dto.UpdateRequest, uuid string) (*models.User, error)
	FindByUsername(ctx context.Context, username string) (*models.User, error)
	FindByEmail(ctx context.Context, email string) (*models.User, error)
	FindByUUID(ctx context.Context, uuid string) (*models.User, error)
}

func NewUserRepository(db *gorm.DB) UserRepositoryInterface {
	return &UserRepository{db: db}
}

func (u *UserRepository) Register(ctx context.Context, request *dto.RegisterRequest) (*models.User, error) {
	user := models.User{
		UUID:        uuid.New(),
		Name:        request.Name,
		Username:    request.Username,
		Email:       request.Email,
		Password:    request.Password,
		PhoneNumber: request.PhoneNumber,
		RoleID:      request.RoleID,
	}

	if err := u.db.WithContext(ctx).Create(&user).Error; err != nil {
		return nil, wrap.WrapError(errConstant.ErrSQLError)
	}

	return &user, nil
}

func (u *UserRepository) Update(ctx context.Context, request *dto.UpdateRequest, uuid string) (*models.User, error) {
	user := models.User{
		Name:        request.Name,
		Username:    request.Username,
		Email:       request.Email,
		Password:    request.Password,
		PhoneNumber: request.PhoneNumber,
	}

	result := u.db.WithContext(ctx).Where("uuid = ?", uuid).Updates(&user)

	if result.Error != nil {
		return nil, wrap.WrapError(errConstant.ErrSQLError)
	}

	if result.RowsAffected == 0 {
		return nil, wrap.WrapError(errConstant.ErrUserNotFound)
	}

	var updatedUser models.User
	if err := u.db.WithContext(ctx).Where("uuid = ?", uuid).First(&updatedUser).Error; err != nil {
		return nil, wrap.WrapError(errConstant.ErrSQLError)
	}

	return &updatedUser, nil
}

func (u *UserRepository) FindByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	result := u.db.WithContext(ctx).Preload("Role").Where("username = ?", username).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, wrap.WrapError(errConstant.ErrUserNotFound)
		}
		return nil, wrap.WrapError(errConstant.ErrSQLError)
	}

	return &user, nil
}

func (u *UserRepository) FindByEmail(ctx context.Context, email string) (*models.User, error) {
	var user models.User
	result := u.db.WithContext(ctx).Preload("Role").Where("email = ?", email).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, wrap.WrapError(errConstant.ErrUserNotFound)
		}
		return nil, wrap.WrapError(errConstant.ErrSQLError)
	}

	return &user, nil
}

func (u *UserRepository) FindByUUID(ctx context.Context, uuid string) (*models.User, error) {
	var user models.User
	result := u.db.WithContext(ctx).Preload("Role").Where("uuid = ?", uuid).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, wrap.WrapError(errConstant.ErrUserNotFound)
		}
		return nil, wrap.WrapError(errConstant.ErrSQLError)
	}

	return &user, nil
}
