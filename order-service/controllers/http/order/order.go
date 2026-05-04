package controllers

import (
	"net/http"
	error2 "order-service/common/error"
	"order-service/common/response"
	"order-service/domain/dto"
	"order-service/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type OrderController struct {
	service services.ServiceRegistryInterface
}

type OrderControllerInterface interface {
	GetAllWithPagination(ctx *gin.Context)
	GetByUUID(ctx *gin.Context)
	GetOrderByUserID(ctx *gin.Context)
	Create(ctx *gin.Context)
}

func NewOrderController(service services.ServiceRegistryInterface) OrderControllerInterface {
	return &OrderController{service: service}
}

func (o *OrderController) GetAllWithPagination(ctx *gin.Context) {
	var params dto.OrderRequestParam
	err := ctx.ShouldBindQuery(&params)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err, nil)
		return
	}

	validate := validator.New()
	if err = validate.Struct(&params); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errorResponse := error2.ErrValidationResponse(err)
		response.Error(ctx, http.StatusUnprocessableEntity, err, &errMessage, errorResponse)
		return
	}

	result, err := o.service.GetOrder().GetAllWithPagination(ctx.Request.Context(), &params)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, result, nil)
}

func (o *OrderController) GetByUUID(ctx *gin.Context) {
	uuid := ctx.Param("uuid")
	if uuid == "" {
		response.Error(ctx, http.StatusBadRequest, nil, nil)
		return
	}

	result, err := o.service.GetOrder().GetOrderByUUID(ctx.Request.Context(), uuid)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, result, nil)
}

func (o *OrderController) GetOrderByUserID(ctx *gin.Context) {
	result, err := o.service.GetOrder().GetOrderByUserID(ctx.Request.Context())
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, result, nil)
}

func (o *OrderController) Create(ctx *gin.Context) {
	var req dto.OrderRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.Error(ctx, http.StatusBadRequest, err, nil)
		return
	}

	validate := validator.New()
	if err := validate.Struct(&req); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errorResponse := error2.ErrValidationResponse(err)
		response.Error(ctx, http.StatusUnprocessableEntity, err, &errMessage, errorResponse)
		return
	}

	result, err := o.service.GetOrder().Create(ctx.Request.Context(), &req)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusCreated, result, nil)
}
