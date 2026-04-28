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

// GetAllWithPagination godoc
// @Summary Get field schedules with pagination
// @Description Get a paginated list of field schedules
// @Tags field-schedules
// @Accept json
// @Produce json
// @Param page query int true "Page number" minimum(1)
// @Param limit query int true "Number of items per page" minimum(1)
// @Param sort_column query string true "Column to sort by" Enums(id, date, time, status, created_at, updated_at)
// @Param sort_order query string true "Sort order" Enums(asc, desc)
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=util.PaginationResult}
// @Failure 400 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /field/schedule/pagination [get]
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

// GetAllByFieldIDAndDate godoc
// @Summary Get schedules by Field ID and Date
// @Description Get available schedules for a field on a specific date
// @Tags field-schedules
// @Accept json
// @Produce json
// @Param uuid path string true "Field UUID"
// @Param date query string true "Date (YYYY-MM-DD)"
// @Success 200 {object} response.Response{data=[]dto.FieldScheduleBookingResponse}
// @Failure 400 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /field/schedule/lists/{uuid} [get]
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

// GetByUUID godoc
// @Summary Get field schedule by UUID
// @Description Get detailed information about a field schedule
// @Tags field-schedules
// @Accept json
// @Produce json
// @Param uuid path string true "Schedule UUID"
// @Success 200 {object} response.Response{data=dto.FieldScheduleResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /field/schedule/{uuid} [get]
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

// Create godoc
// @Summary Create a field schedule
// @Description Create a specific field schedule
// @Tags field-schedules
// @Accept json
// @Produce json
// @Param request body dto.FieldScheduleRequest true "Create Schedule Request"
// @Security ApiKeyAuth
// @Success 201 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /field/schedule [post]
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

// Update godoc
// @Summary Update field schedule
// @Description Update schedule details
// @Tags field-schedules
// @Accept json
// @Produce json
// @Param uuid path string true "Schedule UUID"
// @Param request body dto.UpdateFieldScheduleRequest true "Update Schedule Request"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=dto.FieldScheduleResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /field/schedule/{uuid} [put]
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

// Delete godoc
// @Summary Delete field schedule
// @Description Delete a schedule by UUID
// @Tags field-schedules
// @Accept json
// @Produce json
// @Param uuid path string true "Schedule UUID"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /field/schedule/{uuid} [delete]
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

// UpdateStatus godoc
// @Summary Update field schedule status
// @Description Update schedule status (e.g., Booked)
// @Tags field-schedules
// @Accept json
// @Produce json
// @Param request body dto.UpdateStatusFieldScheduleRequest true "Update Status Request"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /field/schedule [patch]
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

// GenerateFieldScheduleForOneMonth godoc
// @Summary Generate field schedules for one month
// @Description Generate daily schedules for a specific month and year
// @Tags field-schedules
// @Accept json
// @Produce json
// @Param request body dto.GenerateFieldScheduleForOneMonthRequest true "Generate Monthly Schedule Request"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /field/schedule/one-month [post]
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

