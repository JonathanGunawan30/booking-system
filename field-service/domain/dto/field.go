package dto

import (
	"time"

	"github.com/google/uuid"
)

type FieldRequest struct {
	Name         string `form:"name" json:"name" validate:"required,min=3,max=100"`
	Code         string `form:"code" json:"code" validate:"required,min=3,max=15"`
	PricePerHour int    `form:"price_per_hour" json:"price_per_hour" validate:"required"`
}

type UpdateFieldRequest struct {
	Name         string `form:"name" validate:"required,min=3,max=100"`
	Code         string `form:"code" validate:"required,min=3,max=15"`
	PricePerHour int    `form:"price_per_hour" validate:"required"`
}

type FieldResponse struct {
	UUID         uuid.UUID  `json:"uuid"`
	Code         string     `json:"code"`
	Name         string     `json:"name"`
	PricePerHour int        `json:"price_per_hour"`
	Images       []string   `json:"images"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

type FieldDetailResponse struct {
	Code         string     `json:"code"`
	Name         string     `json:"name"`
	PricePerHour int        `json:"price_per_hour"`
	Images       []string   `json:"images"`
	CreatedAt    *time.Time `json:"created_at"`
	UpdatedAt    *time.Time `json:"updated_at"`
}

type FieldRequestParam struct {
	Page       int     `form:"page" validate:"required,gte=1"`
	Limit      int     `form:"limit" validate:"required,gte=1"`
	SortColumn *string `form:"sort_column" validate:"required,oneof=id date time status created_at updated_at"`
	SortOrder  *string `form:"sort_order" validate:"required,oneof=asc desc"`
}
