package dto

import (
	"payment-service/constants"
	"time"

	"github.com/google/uuid"
)

type PaymentRequest struct {
	PaymentLink    string          `json:"payment_link"`
	OrderID        string          `json:"order_id" validate:"required"`
	ExpiredAt      time.Time       `json:"expired_at" validate:"required"`
	Amount         float64         `json:"amount" validate:"required,gt=0"`
	Description    *string         `json:"description,omitempty"`
	CustomerDetail *CustomerDetail `json:"customer_detail,omitempty" validate:"required"`
	ItemDetails    []ItemDetail    `json:"item_details" validate:"required,gt=0,dive"`
}

type CustomerDetail struct {
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type ItemDetail struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Amount   float64 `json:"amount"`
	Quantity int     `json:"quantity"`
}

type PaymentRequestParam struct {
	Page       int     `form:"page" json:"page" validate:"required,gte=1"`
	Limit      int     `form:"limit" json:"limit" validate:"required,gte=1,lte=100"`
	SortColumn *string `form:"sort_column" json:"sort_column" validate:"omitempty,oneof=id order_id expired_at amount status"`
	SortOrder  *string `form:"sort_order" json:"sort_order" validate:"omitempty,oneof=asc desc"`
}

type UpdatePaymentRequest struct {
	TransactionID *string                  `json:"transaction_id" validate:"omitempty"`
	Status        *constants.PaymentStatus `json:"status" validate:"omitempty"`
	PaidAt        *time.Time               `json:"paid_at" validate:"omitempty"`
	VANumber      *string                  `json:"va_number" validate:"omitempty,min=10,max=20"`
	Bank          *string                  `json:"bank" validate:"omitempty,oneof=bca,bni,bri,mandiri,permata,cimb,danamon"`
	InvoiceLink   *string                  `json:"invoice_link" validate:"omitempty,url"`
	Acquirer      *string                  `json:"acquirer" validate:"omitempty,oneof=bca,bni,bri,mandiri,permata,cimb,danamon"`
}

type PaymentResponse struct {
	UUID          uuid.UUID                     `json:"uuid"`
	OrderID       uuid.UUID                     `json:"order_id"`
	Amount        float64                       `json:"amount"`
	Status        constants.PaymentStatusString `json:"status"`
	PaymentLink   string                        `json:"payment_link"`
	InvoiceLink   *string                       `json:"invoice_link,omitempty"`
	TransactionID *string                       `json:"transaction_id,omitempty"`
	VANumber      *string                       `json:"va_number,omitempty"`
	Bank          *string                       `json:"bank,omitempty"`
	Acquirer      *string                       `json:"acquirer,omitempty"`
	Description   *string                       `json:"description"`
	PaidAt        *time.Time                    `json:"paid_at,omitempty"`
	ExpiredAt     *time.Time                    `json:"expired_at"`
	CreatedAt     *time.Time                    `json:"created_at"`
	UpdatedAt     *time.Time                    `json:"updated_at"`
}

type WebHook struct {
	VANumbers         []VANumnber                   `json:"va_numbers"`
	TransactionTime   string                        `json:"transaction_time"`
	TransactionStatus constants.PaymentStatusString `json:"transaction_status"`
	TransactionID     string                        `json:"transaction_id"`
	StatusMessage     string                        `json:"status_message"`
	StatusCode        string                        `json:"status_code"`
	SignatureKey      string                        `json:"signature_key"`
	SettlementTime    string                        `json:"settlement_time"`
	PaymentType       string                        `json:"payment_type"`
	PaymentAmount     []PaymentAmount               `json:"payment_amount"`
	OrderID           uuid.UUID                     `json:"order_id"`
	MerchantID        string                        `json:"merchant_id"`
	GrossAmount       string                        `json:"gross_amount"`
	FraudStatus       string                        `json:"fraud_status"`
	Currency          string                        `json:"currency"`
	Acquirer          *string                       `json:"acquirer"`
}

type VANumnber struct {
	Bank     string `json:"bank"`
	VANumber string `json:"va_number"`
}

type PaymentAmount struct {
	PaidAt *string `json:"paid_at"`
	Amount *string `json:"amount"`
}
