package dto

import (
	"field-service/constants"
	"time"

	"github.com/google/uuid"
)

type FieldScheduleRequest struct {
	FieldID string   `json:"field_id" validate:"required,uuid4"`
	Date    string   `json:"date" validate:"required"`
	TimeIDs []string `json:"time_ids" validate:"required,dive,uuid4"`
}

type GenerateFieldScheduleForOneMonthRequest struct {
	FieldID string `json:"field_id" validate:"required,uuid4"`
	Month   int    `json:"month" validate:"required,gte=1,lte=12"`
	Year    int    `json:"year" validate:"required"`
}

type UpdateFieldScheduleRequest struct {
	Date   string `form:"date" validate:"required,datetime=2006-01-02"`
	TimeID string `json:"time_id" validate:"required,uuid4"`
}
type UpdateStatusFieldScheduleRequest struct {
	FieldScheduleIDs []string `json:"field_schedule_ids" validate:"required,dive,uuid4"`
	Status           int      `json:"status" validate:"oneof=100 200"`
}

type FieldScheduleResponse struct {
	UUID         uuid.UUID                         `json:"uuid"`
	FieldName    string                            `json:"field_name"`
	PricePerHour int                               `json:"price_per_hour"`
	Date         string                            `json:"date"`
	Status       constants.FieldScheduleStatusName `json:"status"`
	Time         string                            `json:"time"`
	CreatedAt    *time.Time                        `json:"created_at"`
	UpdatedAt    *time.Time                        `json:"updated_at"`
}

type FieldScheduleBookingResponse struct {
	UUID         uuid.UUID                         `json:"uuid"`
	PricePerHour string                            `json:"price_per_hour"`
	Date         string                            `json:"date"`
	Status       constants.FieldScheduleStatusName `json:"status"`
	Time         string                            `json:"time"`
}

type FieldScheduleRequestParam struct {
	Page       int     `form:"page" validate:"required,gte=1"`
	Limit      int     `form:"limit" validate:"required,gte=1"`
	SortColumn *string `form:"sort_column" validate:"required,oneof=id date time status created_at updated_at"`
	SortOrder  *string `form:"sort_order" validate:"required,oneof=asc desc"`
}

type FieldScheduleByFieldIDAndDateRequestParam struct {
	Date string `form:"date" validate:"required"`
}
