package dto

import (
	"order-service/constants"
	"time"

	"github.com/google/uuid"
)

type OrderRequest struct {
	FieldScheduleIDs []string `json:"field_schedule_ids" validate:"required"`
}

type OrderRequestParam struct {
	Page       int     `form:"page" json:"page" validate:"required,gte=1"`
	Limit      int     `form:"limit" json:"limit" validate:"required,gte=1,lte=100"`
	SortColumn *string `form:"sort_column" json:"sort_column" validate:"omitempty,oneof=id order_id expired_at amount status created_at"`
	SortOrder  *string `form:"sort_order" json:"sort_order" validate:"omitempty,oneof=asc desc"`
}

type OrderResponse struct {
	UUID        uuid.UUID                   `json:"uuid"`
	Code        string                      `json:"code"`
	UserName    string                      `json:"user_name"`
	Amount      float64                     `json:"amount"`
	Status      constants.OrderStatusString `json:"status"`
	PaymentLink string                      `json:"payment_link"`
	InvoiceLink *string                     `json:"invoice_link,omitempty"`
	Schedules   []FieldData                 `json:"schedules"`
	OrderDate   time.Time                   `json:"order_date"`
	CreatedAt   time.Time                   `json:"created_at"`
	UpdatedAt   time.Time                   `json:"updated_at"`
}

type OrderByUserIDResponse struct {
	Code        string                      `json:"code"`
	Amount      string                      `json:"amount"`
	Status      constants.OrderStatusString `json:"status"`
	OrderDate   time.Time                   `json:"order_date"`
	PaymentLink string                      `json:"payment_link"`
	InvoiceLink *string                     `json:"invoice_link,omitempty"`
	Schedules   []FieldData                 `json:"schedules"`
}
