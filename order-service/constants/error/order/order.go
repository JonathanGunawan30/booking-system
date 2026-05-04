package error

import "errors"

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrFieldAlreadyBooked = errors.New("field already booked")
)

var OrderErrors = []error{
	ErrOrderNotFound,
	ErrFieldAlreadyBooked,
}
