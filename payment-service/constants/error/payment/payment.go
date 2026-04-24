package error

import "errors"

var (
	ErrPaymentNotFound  = errors.New("payment not found")
	ErrExpiredAtInvalid = errors.New("expired time must be greater than current time")
)

var PaymentErrors = []error{
	ErrPaymentNotFound,
	ErrExpiredAtInvalid,
}
