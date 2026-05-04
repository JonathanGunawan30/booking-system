package clients

import (
	"order-service/constants"
	"time"

	"github.com/google/uuid"
)

type FieldResponse struct {
	Code    int       `json:"code"`
	Status  string    `json:"status"`
	Message string    `json:"message"`
	Data    FieldData `json:"data"`
}

type FieldData struct {
	UUID         uuid.UUID                   `json:"uuid"`
	FieldName    string                      `json:"field_name"`
	PricePerHour float64                     `json:"price_per_hour"`
	Date         string                      `json:"date"`
	StartTime    string                      `json:"start_time"`
	EndTime      string                      `json:"end_time"`
	Status       constants.FieldStatusString `json:"status"`
	CreatedAt    *time.Time                  `json:"created_at"`
	UpdatedAt    *time.Time                  `json:"updated_at"`
}
