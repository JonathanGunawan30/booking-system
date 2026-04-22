package constants

import (
	"errors"
	errField "field-service/constants/error/field"
	errFieldSchedule "field-service/constants/error/fieldschedule"
	errTime "field-service/constants/error/time"
	"net/http"
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

	errField.ErrFieldNotFound:                  http.StatusNotFound,
	errFieldSchedule.ErrFieldScheduleNotFound:  http.StatusNotFound,
	errFieldSchedule.ErrFieldScheduleExist:     http.StatusConflict,
	errFieldSchedule.ErrFieldScheduleYearPast:  http.StatusUnprocessableEntity,
	errFieldSchedule.ErrFieldScheduleMonthPast: http.StatusUnprocessableEntity,

	errTime.ErrTimeNotFound: http.StatusNotFound,
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
	allErrors = append(allErrors, errField.FieldErrors...)
	allErrors = append(allErrors, errFieldSchedule.FieldScheduleErrors...)
	allErrors = append(allErrors, errTime.TimeErrors...)

	for _, item := range allErrors {
		if errors.Is(err, item) {
			return true
		}
	}
	return false
}
