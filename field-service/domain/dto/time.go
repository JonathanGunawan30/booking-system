package dto

import (
	"time"

	"github.com/google/uuid"
)

type TimeRequest struct {
	StartTime string `json:"start_time" validate:"required"`
	EndTime   string `json:"end_time" validate:"required"`
}

type TimeResponse struct {
	UUID      uuid.UUID  `json:"uuid"`
	StartTime string     `json:"start_time"`
	EndTime   string     `json:"end_time"`
	CreatedAt *time.Time `json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
}
