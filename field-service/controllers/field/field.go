package fieldcontroller

import (
	errWrap "field-service/common/error"
	"field-service/common/response"
	"field-service/domain/dto"
	"field-service/services"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type FieldController struct {
	service services.ServiceRegistryInterface
}

type FieldControllerInterface interface {
	GetAllWithPagination(ctx *gin.Context)
	GetAllWithoutPagination(ctx *gin.Context)
	GetByUUID(ctx *gin.Context)
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

func NewFieldController(service services.ServiceRegistryInterface) FieldControllerInterface {
	return &FieldController{service: service}
}

// GetAllWithPagination godoc
// @Summary Get fields with pagination
// @Description Get a paginated list of fields
// @Tags fields
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
// @Router /field/pagination [get]
func (f *FieldController) GetAllWithPagination(ctx *gin.Context) {
	var params dto.FieldRequestParam
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

	result, err := f.service.GetField().GetAllWithPagination(ctx.Request.Context(), &params)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, result, nil)
}

// GetAllWithoutPagination godoc
// @Summary Get all fields without pagination
// @Description Get a list of all fields
// @Tags fields
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=[]dto.FieldResponse}
// @Failure 500 {object} response.Response
// @Router /field [get]
func (f *FieldController) GetAllWithoutPagination(ctx *gin.Context) {
	result, err := f.service.GetField().GetAllWithoutPagination(ctx.Request.Context())
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, result, nil)
}

// GetByUUID godoc
// @Summary Get field by UUID
// @Description Get detailed information about a field
// @Tags fields
// @Accept json
// @Produce json
// @Param uuid path string true "Field UUID"
// @Success 200 {object} response.Response{data=dto.FieldResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /field/{uuid} [get]
func (f *FieldController) GetByUUID(ctx *gin.Context) {
	id := ctx.Param("uuid")
	if _, err := uuid.Parse(id); err != nil {
		errMessage := "invalid uuid format"
		errResponse := errWrap.ErrValidationResponse(err)
		response.Error(ctx, http.StatusBadRequest, err, &errMessage, errResponse)
		return
	}

	result, err := f.service.GetField().GetByUUID(ctx.Request.Context(), id)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}
	response.Success(ctx, http.StatusOK, result, nil)
}

// Create godoc
// @Summary Create a new field
// @Description Create a new field with images
// @Tags fields
// @Accept multipart/form-data
// @Produce json
// @Param name formData string true "Field Name"
// @Param code formData string true "Field Code"
// @Param price_per_hour formData int true "Price Per Hour"
// @Param images formData file true "Field Images"
// @Security ApiKeyAuth
// @Success 201 {object} response.Response{data=dto.FieldResponse}
// @Failure 400 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /field [post]
func (f *FieldController) Create(ctx *gin.Context) {
	var request dto.FieldRequest

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

	form, _ := ctx.MultipartForm()
	images := form.File["images"]

	var files []multipart.FileHeader
	for _, image := range images {
		files = append(files, *image)
	}

	result, err := f.service.GetField().Create(ctx.Request.Context(), &request, files)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusCreated, result, nil)
}

// Update godoc
// @Summary Update an existing field
// @Description Update field details and images
// @Tags fields
// @Accept multipart/form-data
// @Produce json
// @Param uuid path string true "Field UUID"
// @Param name formData string true "Field Name"
// @Param code formData string true "Field Code"
// @Param price_per_hour formData int true "Price Per Hour"
// @Param images formData file false "Field Images"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=dto.FieldResponse}
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Failure 422 {object} response.Response
// @Router /field/{uuid} [put]
func (f *FieldController) Update(ctx *gin.Context) {
	id := ctx.Param("uuid")
	if _, err := uuid.Parse(id); err != nil {
		errMessage := "invalid uuid format"
		errResponse := errWrap.ErrValidationResponse(err)
		response.Error(ctx, http.StatusBadRequest, err, &errMessage, errResponse)
		return
	}

	var request dto.UpdateFieldRequest

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

	form, _ := ctx.MultipartForm()
	images := form.File["images"]

	var files []multipart.FileHeader
	for _, image := range images {
		files = append(files, *image)
	}

	result, err := f.service.GetField().Update(ctx.Request.Context(), id, &request, files)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, result, nil)
}

// Delete godoc
// @Summary Delete a field
// @Description Delete a field by UUID
// @Tags fields
// @Accept json
// @Produce json
// @Param uuid path string true "Field UUID"
// @Security ApiKeyAuth
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Failure 404 {object} response.Response
// @Router /field/{uuid} [delete]
func (f *FieldController) Delete(ctx *gin.Context) {
	id := ctx.Param("uuid")
	if _, err := uuid.Parse(id); err != nil {
		errMessage := "invalid uuid format"
		errResponse := errWrap.ErrValidationResponse(err)
		response.Error(ctx, http.StatusBadRequest, err, &errMessage, errResponse)
		return
	}

	err := f.service.GetField().Delete(ctx.Request.Context(), id)
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, nil, nil)
}
