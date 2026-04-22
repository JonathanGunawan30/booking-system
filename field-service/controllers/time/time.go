package time

import (
	errWrap "field-service/common/error"
	"field-service/common/response"
	"field-service/domain/dto"
	"field-service/services"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type TimeController struct {
	service services.ServiceRegistryInterface
}

type TimeControllerInterface interface {
	GetAll(ctx *gin.Context)
	GetByUUID(ctx *gin.Context)
	Create(ctx *gin.Context)
}

func NewTimeController(service services.ServiceRegistryInterface) TimeControllerInterface {
	return &TimeController{service: service}
}

func (t *TimeController) GetAll(ctx *gin.Context) {
	result, err := t.service.GetTime().GetAll(ctx.Request.Context())
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, result, nil)
}

func (t *TimeController) GetByUUID(ctx *gin.Context) {
	id := ctx.Param("uuid")
	if _, err := uuid.Parse(id); err != nil {
		errMessage := "invalid uuid format"
		errResponse := errWrap.ErrValidationResponse(err)
		response.Error(ctx, http.StatusBadRequest, err, &errMessage, errResponse)
		return
	}

	result, err := t.service.GetTime().GetByUUID(ctx.Request.Context(), id)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, result, nil)
}

func (t *TimeController) Create(ctx *gin.Context) {
	var request dto.TimeRequest
	if err := ctx.ShouldBindJSON(&request); err != nil {
		response.Error(ctx, http.StatusBadRequest, err, nil, nil)
		return
	}

	validate := validator.New()
	if err := validate.Struct(&request); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.ErrValidationResponse(err)
		response.Error(ctx, http.StatusUnprocessableEntity, err, &errMessage, errResponse)
		return
	}

	result, err := t.service.GetTime().Create(ctx.Request.Context(), &request)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusCreated, result, nil)
}
