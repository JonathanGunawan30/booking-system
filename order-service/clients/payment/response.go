package clients

import (
	"time"

	"github.com/google/uuid"
)

type PaymentResponse struct {
	Code    int         `json:"code"`
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Data    PaymentData `json:"data"`
}

type PaymentData struct {
	UUID          uuid.UUID  `json:"uuid"`
	OrderID       string     `json:"order_id"`
	Amount        float64    `json:"amount"`
	Status        string     `json:"status"`
	PaymentLink   string     `json:"payment_link"`
	InvoiceLink   *string    `json:"invoice_link,omitempty"`
	Description   *string    `json:"description"`
	VANumber      *string    `json:"va_number,omitempty"`
	BankName      *string    `json:"bank_name,omitempty"`
	TransactionID *string    `json:"transaction_id,omitempty"`
	Acquirer      *string    `json:"acquirer,omitempty"`
	PaidAt        *string    `json:"paid_at,omitempty"`
	ExpiredAt     string     `json:"expired_at"`
	CreatedAt     *time.Time `json:"created_at"`
	UpdatedAt     *time.Time `json:"updated_at"`
}
