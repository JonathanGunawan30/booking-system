package error

import (
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
)

type ValidationResponse struct {
	Field   string `json:"field,omitempty"`
	Message string `json:"message,omitempty"`
}

var ErrValidator = map[string]string{
	"required": "%s is required",
	"email":    "%s is not a valid email",
	"min":      "%s must be at least %s characters",
	"max":      "%s must be at most %s characters",
}

func ErrValidationResponse(err error) []ValidationResponse {
	var validationResponse []ValidationResponse

	var fieldErrors validator.ValidationErrors
	if errors.As(err, &fieldErrors) {
		for _, fieldErr := range fieldErrors {
			msg, ok := ErrValidator[fieldErr.Tag()]
			if ok {
				var message string
				if fieldErr.Param() != "" {
					message = fmt.Sprintf(msg, fieldErr.Field(), fieldErr.Param())
				} else {
					message = fmt.Sprintf(msg, fieldErr.Field())
				}
				validationResponse = append(validationResponse, ValidationResponse{
					Field:   fieldErr.Field(),
					Message: message,
				})
			} else {
				validationResponse = append(validationResponse, ValidationResponse{
					Field:   fieldErr.Field(),
					Message: fmt.Sprintf("%s is invalid", fieldErr.Field()),
				})
			}
		}
	}

	return validationResponse
}

func WrapError(err error) error {
	logrus.Errorf("error: %v", err)
	return err
}
