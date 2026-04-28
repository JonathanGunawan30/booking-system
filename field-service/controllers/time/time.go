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

// GetAll godoc
// @Summary Get all operational times
// @Description Get a list of all operational times
// @Tags times
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=[]dto.TimeResponse}
// @Failure 500 {object} response.Response
// @Router /time [get]
func (t *TimeController) GetAll(ctx *gin.Context) {
	result, err := t.service.GetTime().GetAll(ctx.Request.Context())
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, result, nil)
}

// GetByUUID godoc
// @Summary Get time detail by UUID
// @Description Get detailed information about an operational time
// @Tags times
// @Accept json
// @Produce json
// @Param uuid path string true "Time UUID"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=dto.TimeResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /time/{uuid} [get]
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

// Create godoc
// @Summary Create a new operational time
// @Description Create a new operational time slot
// @Tags times
// @Accept json
// @Produce json
// @Param request body dto.TimeRequest true "Create Time Request"
// @Security ApiKeyAuth
// @Success 201 {object} response.Response{data=dto.TimeResponse}
// @Failure 400 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /time [post]
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

