package error

import "errors"

var (
	ErrOrderNotFound           = errors.New("order not found")
	ErrFieldAlreadyBooked      = errors.New("field already booked")
	ErrFieldNotFound           = errors.New("field not found")
	ErrOrderFailed             = errors.New("failed to create order")
	ErrPaymentLinkCreationFailed = errors.New("failed to create payment link")
)

var OrderErrors = []error{
	ErrOrderNotFound,
	ErrFieldAlreadyBooked,
	ErrFieldNotFound,
	ErrOrderFailed,
	ErrPaymentLinkCreationFailed,
}
