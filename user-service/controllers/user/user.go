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

func (u *UserController) GetUserLogin(ctx *gin.Context) {
	user, err := u.userService.GetUserLogin(ctx.Request.Context())
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, user, nil)
}

func (u *UserController) GetUserByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")

	result, err := u.userService.GetUserByUUID(ctx.Request.Context(), uuid)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, result, nil)
}
