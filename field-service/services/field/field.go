package services

import (
	"bytes"
	"context"
	"field-service/common/cloudflare"
	"field-service/common/util"
	"field-service/config"
	"field-service/constants/error"
	"field-service/domain/dto"
	"field-service/domain/models"
	"field-service/repositories"
	"fmt"
	"io"
	"mime/multipart"
	"path"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type FieldService struct {
	repository repositories.RepositoryRegistryInterface
	r2         cloudflare.R2Client
}

type FieldServiceInterface interface {
	GetAllWithPagination(ctx context.Context, param *dto.FieldRequestParam) (*util.PaginationResult, error)
	GetAllWithoutPagination(ctx context.Context) ([]dto.FieldResponse, error)
	GetByUUID(ctx context.Context, uuid string) (*dto.FieldResponse, error)
	Create(ctx context.Context, req *dto.FieldRequest, images []multipart.FileHeader) (*dto.FieldResponse, error)
	Update(ctx context.Context, uuid string, req *dto.UpdateFieldRequest, images []multipart.FileHeader) (*dto.FieldResponse, error)
	Delete(ctx context.Context, uuid string) error
}

func NewFieldService(repository repositories.RepositoryRegistryInterface, r2 cloudflare.R2Client) FieldServiceInterface {
	return &FieldService{repository: repository, r2: r2}
}

func (f *FieldService) GetAllWithPagination(ctx context.Context, param *dto.FieldRequestParam) (*util.PaginationResult, error) {
	fields, total, err := f.repository.GetField().FindAllWithPagination(ctx, param)
	if err != nil {
		return nil, err
	}

	fieldResults := make([]dto.FieldResponse, 0, len(fields))
	for _, field := range fields {
		fieldResults = append(fieldResults, dto.FieldResponse{
			UUID:         field.UUID,
			Code:         field.Code,
			Name:         field.Name,
			PricePerHour: field.PricePerHour,
			Images:       field.Image,
			CreatedAt:    field.CreatedAt,
			UpdatedAt:    field.UpdatedAt,
		})
	}

	pagination := &util.PaginationParam{
		Count: total,
		Page:  param.Page,
		Limit: param.Limit,
		Data:  fieldResults,
	}

	response := util.GeneratePagination(*pagination)
	return &response, nil
}

func (f *FieldService) GetAllWithoutPagination(ctx context.Context) ([]dto.FieldResponse, error) {
	fields, err := f.repository.GetField().FindAllWithoutPagination(ctx)
	if err != nil {
		return nil, err
	}

	fieldResults := make([]dto.FieldResponse, 0, len(fields))
	for _, field := range fields {
		fieldResults = append(fieldResults, dto.FieldResponse{
			UUID:         field.UUID,
			Code:         field.Code,
			Name:         field.Name,
			PricePerHour: field.PricePerHour,
			Images:       field.Image,
			CreatedAt:    field.CreatedAt,
			UpdatedAt:    field.UpdatedAt,
		})
	}

	return fieldResults, nil
}

func (f *FieldService) GetByUUID(ctx context.Context, uuid string) (*dto.FieldResponse, error) {
	fields, err := f.repository.GetField().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	return toFieldResponse(fields), nil
}

func (f *FieldService) validateUpload(images []multipart.FileHeader) error {
	if len(images) == 0 {
		return constants.ErrInvalidUploadFile
	}

	allowedExt := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true}
	for _, image := range images {
		if image.Size > int64(config.AppConfig.MaxUploadSize)*1024*1024 {
			return constants.ErrFileTooLarge
		}
		ext := strings.ToLower(path.Ext(image.Filename))
		if !allowedExt[ext] {
			return constants.ErrInvalidUploadFile
		}
	}

	return nil
}

func (f *FieldService) processAndUploadImage(image multipart.FileHeader) (string, error) {
	file, err := image.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, file)
	if err != nil {
		return "", err
	}

	ext := strings.ToLower(path.Ext(image.Filename))

	filename := fmt.Sprintf("fields/%s-%s%s",
		time.Now().Format("2006-01-02"),
		uuid.New().String(),
		ext,
	)

	url, err := f.r2.Upload(filename, buffer, image.Header.Get("Content-Type"))
	if err != nil {
		return "", err
	}
	return url, nil
}

func (f *FieldService) Create(ctx context.Context, req *dto.FieldRequest, images []multipart.FileHeader) (*dto.FieldResponse, error) {
	if err := f.validateUpload(images); err != nil {
		return nil, err
	}

	var imageURLs []string
	for _, image := range images {
		url, err := f.processAndUploadImage(image)
		if err != nil {
			return nil, err
		}
		imageURLs = append(imageURLs, url)
	}

	field := models.Field{
		Code:         req.Code,
		Name:         req.Name,
		PricePerHour: req.PricePerHour,
		Image:        pq.StringArray(imageURLs),
	}

	result, err := f.repository.GetField().Create(ctx, &field)
	if err != nil {
		return nil, err
	}

	return &dto.FieldResponse{
		UUID:         result.UUID,
		Code:         result.Code,
		Name:         result.Name,
		PricePerHour: result.PricePerHour,
		Images:       imageURLs,
	}, nil
}

func (f *FieldService) Update(ctx context.Context, uuid string, req *dto.UpdateFieldRequest, images []multipart.FileHeader) (*dto.FieldResponse, error) {
	field, err := f.repository.GetField().FindByUUID(ctx, uuid)
	if err != nil {
		return nil, err
	}

	field.Code = req.Code
	field.Name = req.Name
	field.PricePerHour = req.PricePerHour

	if len(images) > 0 {
		if err := f.validateUpload(images); err != nil {
			return nil, err
		}
		oldImages := make([]string, len(field.Image))
		copy(oldImages, field.Image)

		var imageURLs []string
		for _, image := range images {
			url, err := f.processAndUploadImage(image)
			if err != nil {
				for _, uploadedURL := range imageURLs {
					key := strings.Replace(uploadedURL, config.AppConfig.R2PublicURL+"/", "", 1)
					f.r2.Delete(key)
				}
				return nil, err
			}
			imageURLs = append(imageURLs, url)
		}

		field.Image = imageURLs

		updated, err := f.repository.GetField().Update(ctx, uuid, field)
		if err != nil {
			for _, uploadedURL := range imageURLs {
				key := strings.Replace(uploadedURL, config.AppConfig.R2PublicURL+"/", "", 1)
				f.r2.Delete(key)
			}
			return nil, err
		}

		for _, oldImage := range oldImages {
			key := strings.Replace(oldImage, config.AppConfig.R2PublicURL+"/", "", 1)
			f.r2.Delete(key)
		}

		return toFieldResponse(updated), nil
	}

	updated, err := f.repository.GetField().Update(ctx, uuid, field)
	if err != nil {
		return nil, err
	}

	return toFieldResponse(updated), nil
}

func (f *FieldService) Delete(ctx context.Context, uuid string) error {
	field, err := f.repository.GetField().FindByUUID(ctx, uuid)
	if err != nil {
		return err
	}

	for _, image := range field.Image {
		key := strings.Replace(image, config.AppConfig.R2PublicURL+"/", "", 1)
		if err := f.r2.Delete(key); err != nil {
			return err
		}
	}

	return f.repository.GetField().Delete(ctx, uuid)
}

func toFieldResponse(field *models.Field) *dto.FieldResponse {
	return &dto.FieldResponse{
		UUID:         field.UUID,
		Code:         field.Code,
		Name:         field.Name,
		PricePerHour: field.PricePerHour,
		Images:       field.Image,
		CreatedAt:    field.CreatedAt,
		UpdatedAt:    field.UpdatedAt,
	}
}
