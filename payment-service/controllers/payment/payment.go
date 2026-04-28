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

// GetAllWithPagination handles listing payments with pagination.
// @Summary List payments
// @Description Get a paginated list of payment transactions
// @Tags Payments
// @Accept json
// @Produce json
// @Param page query int true "Page number"
// @Param limit query int true "Page limit"
// @Param sort_column query string false "Column to sort by" Enums(id, order_id, expired_at, amount, status)
// @Param sort_order query string false "Sort order" Enums(asc, desc)
// @Success 200 {object} response.Response{data=util.PaginationResult}
// @Failure 400 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /payments [get]
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

// GetByUUID handles getting a single payment by UUID.
// @Summary Get payment detail
// @Description Get detailed information about a specific transaction by its UUID
// @Tags Payments
// @Accept json
// @Produce json
// @Param uuid path string true "Payment UUID"
// @Success 200 {object} response.Response{data=dto.PaymentResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /payments/{uuid} [get]
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

// Create handles creating a new payment.
// @Summary Create payment
// @Description Initiate a new payment transaction and get a payment link
// @Tags Payments
// @Accept json
// @Produce json
// @Param request body dto.PaymentRequest true "Payment Request"
// @Success 201 {object} response.Response{data=dto.PaymentResponse}
// @Failure 400 {object} response.Response
// @Failure 422 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /payments [post]
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

// WebHook handles payment notifications from Midtrans.
// @Summary Payment Webhook
// @Description Callback endpoint for Midtrans to update payment status
// @Tags Webhooks
// @Accept json
// @Produce json
// @Param request body dto.WebHook true "Webhook Data"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 500 {object} response.Response
// @Router /payments/webhook [post]
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
