package constatns

import (
	"errors"
	"net/http"
)

var errorCodeMap = map[error]int{
	ErrInternalServerError: http.StatusInternalServerError,
	ErrSQLError:            http.StatusInternalServerError,
	ErrTooManyRequest:      http.StatusTooManyRequests,
	ErrInvalidToken:        http.StatusUnauthorized,
	ErrForbidden:           http.StatusForbidden,
	ErrUnauthorized:        http.StatusUnauthorized,

	ErrUserNotFound:                http.StatusNotFound,
	ErrPasswordIncorrect:           http.StatusUnauthorized,
	ErrUsernameExists:              http.StatusConflict,
	ErrEmailExists:                 http.StatusConflict,
	ErrPasswordDoesNotMatch:        http.StatusBadRequest,
	ErrUsernameOrPasswordIncorrect: http.StatusUnauthorized,
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
	allErrors = append(append(GeneralErrors[:], UserErrors[:]...))

	for _, item := range allErrors {
		if errors.Is(err, item) {
			return true
		}
	}
	return false
}
