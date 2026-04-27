package controller

import (
	"net/http"
	"user-service/common/error"
	"user-service/common/response"
	"user-service/domain/dto"
	services "user-service/services/user"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type UserController struct {
	userService services.UserServiceInterface
}

type UserControllerInterface interface {
	Login(ctx *gin.Context)
	Register(ctx *gin.Context)
	Update(ctx *gin.Context)
	GetUserLogin(ctx *gin.Context)
	GetUserByUUID(ctx *gin.Context)
}

func NewUserController(userService services.UserServiceInterface) UserControllerInterface {
	return &UserController{userService: userService}
}

// Login handles user authentication
// @Summary User Login
// @Description Authenticate user with username and password to get JWT token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login Credentials"
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 401 {object} response.Response
// @Router /auth/login [post]
func (u *UserController) Login(ctx *gin.Context) {
	request := &dto.LoginRequest{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		response.Error(ctx, http.StatusBadRequest, err, nil)
		return
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := error.ErrValidationResponse(err)
		response.Error(ctx, http.StatusUnprocessableEntity, nil, &errMessage, errResponse)
		return
	}

	result, err := u.userService.Login(ctx.Request.Context(), request)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, result.User, &result.Token)

}

// Register handles user registration
// @Summary User Registration
// @Description Register a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration Details"
// @Success 201 {object} response.Response{data=dto.RegisterResponse}
// @Failure 400 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /auth/register [post]
func (u *UserController) Register(ctx *gin.Context) {
	request := &dto.RegisterRequest{}
	if err := ctx.ShouldBindJSON(request); err != nil {
		response.Error(ctx, http.StatusBadRequest, err, nil)
		return
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := error.ErrValidationResponse(err)
		response.Error(ctx, http.StatusUnprocessableEntity, nil, &errMessage, errResponse)
		return
	}

	result, err := u.userService.Register(ctx.Request.Context(), request)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusCreated, result.User, nil)
}

// Update handles user profile update
// @Summary Update User Profile
// @Description Update the profile of a user by UUID
// @Tags User
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param uuid path string true "User UUID"
// @Param request body dto.UpdateRequest true "Update Details"
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 400 {object} response.Response
// @Router /auth/{uuid} [put]
func (u *UserController) Update(ctx *gin.Context) {
	request := &dto.UpdateRequest{}
	uuid := ctx.Param("uuid")

	if err := ctx.ShouldBindJSON(request); err != nil {
		response.Error(ctx, http.StatusBadRequest, err, nil)
		return
	}

	validate := validator.New()
	if err := validate.Struct(request); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := error.ErrValidationResponse(err)
		response.Error(ctx, http.StatusUnprocessableEntity, nil, &errMessage, errResponse)
		return
	}

	result, err := u.userService.Update(ctx.Request.Context(), request, uuid)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, result, nil)
}

// GetUserLogin gets the currently logged-in user's profile
// @Summary Get Current User
// @Description Retrieve the profile of the authenticated user
// @Tags User
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Router /auth/user [get]
func (u *UserController) GetUserLogin(ctx *gin.Context) {
	user, err := u.userService.GetUserLogin(ctx.Request.Context())
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, user, nil)
}

// GetUserByUUID gets a user's profile by UUID
// @Summary Get User By UUID
// @Description Retrieve a specific user's profile using their UUID
// @Tags User
// @Produce json
// @Security ApiKeyAuth
// @Param uuid path string true "User UUID"
// @Success 200 {object} response.Response{data=dto.UserResponse}
// @Failure 400 {object} response.Response
// @Router /auth/user/{uuid} [get]
func (u *UserController) GetUserByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	result, err := u.userService.GetUserByUUID(ctx.Request.Context(), uuid)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, result, nil)
}
