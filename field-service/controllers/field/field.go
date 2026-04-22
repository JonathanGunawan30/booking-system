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

func (f *FieldController) GetAllWithoutPagination(ctx *gin.Context) {
	result, err := f.service.GetField().GetAllWithoutPagination(ctx.Request.Context())
	if err != nil {
		response.ErrorFromApp(ctx, err)
		return
	}

	response.Success(ctx, http.StatusOK, result, nil)
}

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
