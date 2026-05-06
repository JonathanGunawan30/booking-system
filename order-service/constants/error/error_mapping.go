package constants

import (
	"errors"
	"net/http"
	errOrder "order-service/constants/error/order"
)

var errorCodeMap = map[error]int{
	ErrInternalServerError: http.StatusInternalServerError,
	ErrSQLError:            http.StatusInternalServerError,
	ErrTooManyRequest:      http.StatusTooManyRequests,
	ErrInvalidToken:        http.StatusUnauthorized,
	ErrForbidden:           http.StatusForbidden,
	ErrUnauthorized:        http.StatusUnauthorized,
	ErrInvalidUploadFile:   http.StatusBadRequest,
	ErrFileTooLarge:        http.StatusRequestEntityTooLarge,
	ErrInvalidUUIDFormat:   http.StatusBadRequest,
	ErrUserNotFound:        http.StatusNotFound,

	errOrder.ErrOrderNotFound:           http.StatusNotFound,
	errOrder.ErrFieldNotFound:           http.StatusNotFound,
	errOrder.ErrFieldAlreadyBooked:      http.StatusConflict,
	errOrder.ErrOrderFailed:             http.StatusInternalServerError,
	errOrder.ErrPaymentLinkCreationFailed: http.StatusInternalServerError,
}

func GetErrorCode(err error) int {
	for target, code := range errorCodeMap {
		if errors.Is(err, target) {
			return code
		}
	}
	return http.StatusInternalServerError
}

func ErrMapping(err error) bool {
	allErrors := make([]error, 0)
	allErrors = append(allErrors, GeneralErrors...)
	allErrors = append(allErrors, errOrder.OrderErrors...)

	for _, item := range allErrors {
		if errors.Is(err, item) {
			return true
		}
	}
	return false
}
