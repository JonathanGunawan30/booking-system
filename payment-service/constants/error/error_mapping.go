package constants

import (
	"errors"
	"net/http"
	errPayment "payment-service/constants/error/payment"
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

	errPayment.ErrPaymentNotFound:  http.StatusNotFound,
	errPayment.ErrExpiredAtInvalid: http.StatusBadRequest,
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
	allErrors = append(allErrors, errPayment.ErrPaymentNotFound)

	for _, item := range allErrors {
		if errors.Is(err, item) {
			return true
		}
	}
	return false
}
