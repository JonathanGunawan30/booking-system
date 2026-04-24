package controllers

import (
	"net/http"
	errWrap "payment-service/common/error"
	"payment-service/common/response"
	"payment-service/domain/dto"
	"payment-service/services"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type PaymentController struct {
	service services.ServiceRegistryInterface
}

type PaymentControllerInterface interface {
	GetAllWithPagination(*gin.Context)
	GetByUUID(*gin.Context)
	Create(*gin.Context)
	WebHook(*gin.Context)
}

func NewPaymentController(service services.ServiceRegistryInterface) PaymentControllerInterface {
	return &PaymentController{service: service}
}

func (p *PaymentController) GetAllWithPagination(ctx *gin.Context) {
	var param dto.PaymentRequestParam
	err := ctx.ShouldBindQuery(&param)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err, nil)
	}

	validate := validator.New()
	if err := validate.Struct(&param); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.ErrValidationResponse(err)
		response.Error(ctx, http.StatusUnprocessableEntity, err, &errMessage, errResponse)
		return
	}

	pagination, err := p.service.GetPayment().GetAllWithPagination(ctx, &param)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, pagination, nil)
}

func (p *PaymentController) GetByUUID(ctx *gin.Context) {
	id := ctx.Param("uuid")
	if _, err := uuid.Parse(id); err != nil {
		errMessage := "invalid uuid format"
		errResponse := errWrap.ErrValidationResponse(err)
		response.Error(ctx, http.StatusBadRequest, err, &errMessage, errResponse)
		return
	}

	result, err := p.service.GetPayment().GetByUUID(ctx, id)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, result, nil)
}

func (p *PaymentController) Create(ctx *gin.Context) {
	var request dto.PaymentRequest
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err, nil)
	}

	validate := validator.New()
	if err := validate.Struct(&request); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.ErrValidationResponse(err)
		response.Error(ctx, http.StatusUnprocessableEntity, err, &errMessage, errResponse)
		return
	}

	result, err := p.service.GetPayment().Create(ctx, &request)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}
	response.Success(ctx, http.StatusCreated, result, nil)
}

func (p *PaymentController) WebHook(ctx *gin.Context) {
	var request dto.WebHook
	err := ctx.ShouldBindJSON(&request)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err, nil)
	}

	err = p.service.GetPayment().WebHook(ctx, &request)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}
	response.Success(ctx, http.StatusOK, nil, nil)
}
