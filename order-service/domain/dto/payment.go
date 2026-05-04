package dto

import (
	"time"

	"github.com/google/uuid"
)

type PaymentRequest struct {
	OrderID        uuid.UUID      `json:"order_id"`
	Amount         float64        `json:"amount"`
	Description    string         `json:"description"`
	ExpiredAt      time.Time      `json:"expired_at"`
	CustomerDetail CustomerDetail `json:"customer_detail"`
	ItemDetails    []ItemDetails  `json:"item_details"`
}

type CustomerDetail struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	PhoneNumber string `json:"phone_number"`
}

type ItemDetails struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	Amount   float64   `json:"amount"`
	Quantity int       `json:"quantity"`
}
