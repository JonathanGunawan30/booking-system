package services

import (
	"context"
	"errors"
	"strings"
	"time"
	"user-service/config"
	constants "user-service/constants"
	errConstants "user-service/constants/error"
	"user-service/domain/dto"
	"user-service/repositories"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repository repositories.RepositoryRegistryInterface
}

type UserServiceInterface interface {
	Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error)
	Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error)
	Update(ctx context.Context, rep *dto.UpdateRequest, uuid string) (*dto.UserResponse, error)
	GetUserLogin(ctx context.Context) (*dto.UserResponse, error)
	GetUserByUUID(ctx context.Context, uuid string) (*dto.UserResponse, error)
}

type Claims struct {
	User *dto.UserResponse
	jwt.RegisteredClaims
}

func NewUserService(repository repositories.RepositoryRegistryInterface) UserServiceInterface {
	return &UserService{repository: repository}
}

func (u *UserService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := u.repository.GetUser().FindByUsername(ctx, req.Username)
	if err != nil {
		if errors.Is(err, errConstants.ErrUserNotFound) {
			return nil, errConstants.ErrUsernameOrPasswordIncorrect
		}
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, errConstants.ErrUsernameOrPasswordIncorrect
	}

	expirationTime := time.Now().Add(time.Duration(config.AppConfig.JwtExpirationTime) * time.Minute).Unix()

	data := &dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		PhoneNumber: user.PhoneNumber,
		Email:       user.Email,
		Role:        strings.ToLower(user.Role.Name),
	}

	claims := &Claims{
		User: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Unix(expirationTime, 0)),
			Issuer:    config.AppConfig.JwtIssuer,
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(config.AppConfig.JwtSecretKey))

	if err != nil {
		return nil, err
	}

	response := &dto.LoginResponse{
		User:  *data,
		Token: tokenString,
	}

	return response, nil
}

func (u *UserService) isUsernameExists(ctx context.Context, username string) bool {
	user, err := u.repository.GetUser().FindByUsername(ctx, username)
	if err != nil {
		return false
	}
	return user != nil
}

func (u *UserService) isEmailExists(ctx context.Context, email string) bool {
	user, err := u.repository.GetUser().FindByEmail(ctx, email)
	if err != nil {
		return false
	}
	return user != nil
}

func (u *UserService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.RegisterResponse, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	if u.isUsernameExists(ctx, req.Username) {
		return nil, errConstants.ErrUsernameExists
	}

	if u.isEmailExists(ctx, req.Email) {
		return nil, errConstants.ErrEmailExists
	}

	register, err := u.repository.GetUser().Register(ctx, &dto.RegisterRequest{
		Name:        req.Name,
		Username:    req.Username,
		Email:       req.Email,
		Password:    string(hashedPassword),
		PhoneNumber: req.PhoneNumber,
		RoleID:      constants.Customer,
	})

	if err != nil {
		return nil, err
	}

	response := &dto.RegisterResponse{
		User: dto.UserResponse{
			UUID:        register.UUID,
			Name:        register.Name,
			Username:    register.Username,
			PhoneNumber: register.PhoneNumber,
			Email:       register.Email,
		},
	}

	return response, nil

}

func (u *UserService) Update(ctx context.Context, req *dto.UpdateRequest, uuid string) (*dto.UserResponse, error) {
	user, err := u.repository.GetUser().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	isUsernameExists := u.isUsernameExists(ctx, req.Username)
	if isUsernameExists && user.Username != req.Username {
		return nil, errConstants.ErrUsernameExists
	}

	isEmailExists := u.isEmailExists(ctx, req.Email)
	if isEmailExists && user.Email != req.Email {
		return nil, errConstants.ErrEmailExists
	}

	var password string
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		password = string(hashedPassword)
	}

	result, err := u.repository.GetUser().Update(ctx, &dto.UpdateRequest{
		Name:        req.Name,
		Username:    req.Username,
		Password:    password,
		Email:       req.Email,
		PhoneNumber: req.PhoneNumber,
	}, uuid)

	if err != nil {
		return nil, err
	}

	data := dto.UserResponse{
		UUID:        result.UUID,
		Name:        result.Name,
		Username:    result.Username,
		Email:       result.Email,
		PhoneNumber: result.PhoneNumber,
	}

	return &data, nil

}

func (u *UserService) GetUserLogin(ctx context.Context) (*dto.UserResponse, error) {
	userLogin, ok := ctx.Value(constants.UserLogin).(*dto.UserResponse)
	if !ok || userLogin == nil {
		return nil, errConstants.ErrUnauthorized
	}

	data := dto.UserResponse{
		UUID:        userLogin.UUID,
		Name:        userLogin.Name,
		Username:    userLogin.Username,
		Email:       userLogin.Email,
		PhoneNumber: userLogin.PhoneNumber,
		Role:        strings.ToLower(userLogin.Role),
	}
	return &data, nil
}

func (u *UserService) GetUserByUUID(ctx context.Context, uuid string) (*dto.UserResponse, error) {
	logrus.Infof("[UserService] GetUserByUUID: searching for UUID: %s", uuid)
	user, err := u.repository.GetUser().FindByUUID(ctx, uuid)
	if err != nil {
		logrus.Errorf("[UserService] GetUserByUUID error: %v", err)
		return nil, err
	}
	data := dto.UserResponse{
		UUID:        user.UUID,
		Name:        user.Name,
		Username:    user.Username,
		Email:       user.Email,
		PhoneNumber: user.PhoneNumber,
		Role:        strings.ToLower(user.Role.Name),
	}
	return &data, nil
}
