package fieldschedule

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

type FieldScheduleController struct {
	service services.ServiceRegistryInterface
}

type FieldScheduleControllerInterface interface {
	GetAllWithPagination(ctx *gin.Context)
	GetAllByFieldIDAndDate(ctx *gin.Context)
	GetByUUID(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	UpdateStatus(ctx *gin.Context)
	GenerateFieldScheduleForOneMonth(ctx *gin.Context)
}

func NewFieldScheduleController(service services.ServiceRegistryInterface) FieldScheduleControllerInterface {
	return &FieldScheduleController{service: service}
}

func (f *FieldScheduleController) GetAllWithPagination(ctx *gin.Context) {
	var params dto.FieldScheduleRequestParam

	err := ctx.ShouldBindQuery(&params)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err, nil, nil)
		return
	}

	validate := validator.New()
	if err := validate.Struct(&params); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.ErrValidationResponse(err)
		response.Error(ctx, http.StatusUnprocessableEntity, err, &errMessage, errResponse)
		return
	}

	result, err := f.service.GetFieldSchedule().GetAllWithPagination(ctx.Request.Context(), &params)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, result, nil)
}

func (f *FieldScheduleController) GetAllByFieldIDAndDate(ctx *gin.Context) {
	id := ctx.Param("uuid")
	if _, err := uuid.Parse(id); err != nil {
		errMessage := "invalid uuid format"
		errResponse := errWrap.ErrValidationResponse(err)
		response.Error(ctx, http.StatusBadRequest, err, &errMessage, errResponse)
		return
	}

	var params dto.FieldScheduleByFieldIDAndDateRequestParam
	err := ctx.ShouldBindQuery(&params)
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err, nil, nil)
		return
	}

	validate := validator.New()
	if err := validate.Struct(&params); err != nil {
		errMessage := http.StatusText(http.StatusUnprocessableEntity)
		errResponse := errWrap.ErrValidationResponse(err)
		response.Error(ctx, http.StatusUnprocessableEntity, err, &errMessage, errResponse)
		return
	}

	result, err := f.service.GetFieldSchedule().GetAllByFieldIDAndDate(ctx.Request.Context(), id, params.Date)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, result, nil)
}

func (f *FieldScheduleController) GetByUUID(ctx *gin.Context) {
	id := ctx.Param("uuid")
	if _, err := uuid.Parse(id); err != nil {
		errMessage := "invalid uuid format"
		errResponse := errWrap.ErrValidationResponse(err)
		response.Error(ctx, http.StatusBadRequest, err, &errMessage, errResponse)
		return
	}

	result, err := f.service.GetFieldSchedule().GetByUUID(ctx.Request.Context(), id)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, result, nil)
}

func (f *FieldScheduleController) Create(ctx *gin.Context) {
	var request dto.FieldScheduleRequest
	if err := ctx.ShouldBind(&request); err != nil {
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

	err := f.service.GetFieldSchedule().Create(ctx.Request.Context(), &request)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusCreated, nil, nil)
}

func (f *FieldScheduleController) Update(ctx *gin.Context) {
	id := ctx.Param("uuid")
	if _, err := uuid.Parse(id); err != nil {
		errMessage := "invalid uuid format"
		errResponse := errWrap.ErrValidationResponse(err)
		response.Error(ctx, http.StatusBadRequest, err, &errMessage, errResponse)
		return
	}

	var request dto.UpdateFieldScheduleRequest
	if err := ctx.ShouldBind(&request); err != nil {
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

	result, err := f.service.GetFieldSchedule().Update(ctx.Request.Context(), id, &request)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, result, nil)
}

func (f *FieldScheduleController) Delete(ctx *gin.Context) {
	id := ctx.Param("uuid")
	if _, err := uuid.Parse(id); err != nil {
		errMessage := "invalid uuid format"
		errResponse := errWrap.ErrValidationResponse(err)
		response.Error(ctx, http.StatusBadRequest, err, &errMessage, errResponse)
		return
	}

	err := f.service.GetFieldSchedule().Delete(ctx.Request.Context(), id)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, nil, nil)
}

func (f *FieldScheduleController) UpdateStatus(ctx *gin.Context) {
	var request dto.UpdateStatusFieldScheduleRequest
	if err := ctx.ShouldBind(&request); err != nil {
		response.Error(ctx, http.StatusBadRequest, err, nil, nil)
		return
	}

	err := f.service.GetFieldSchedule().UpdateStatus(ctx.Request.Context(), &request)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, nil, nil)
}

func (f *FieldScheduleController) GenerateFieldScheduleForOneMonth(ctx *gin.Context) {
	var request dto.GenerateFieldScheduleForOneMonthRequest
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

	err := f.service.GetFieldSchedule().GenerateFieldScheduleForOneMonth(ctx.Request.Context(), &request)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, nil, nil)
}
