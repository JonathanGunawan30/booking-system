package constants

import "errors"

var (
	ErrInternalServerError = errors.New("internal server error")
	ErrSQLError            = errors.New("database server failed to execute query")
	ErrTooManyRequest      = errors.New("too many requests")
	ErrInvalidToken        = errors.New("invalid token")
	ErrForbidden           = errors.New("forbidden")
	ErrUnauthorized        = errors.New("unauthorized")
	ErrInvalidUploadFile   = errors.New("invalid upload file")
	ErrFileTooLarge        = errors.New("file size exceeds maximum allowed size")
	ErrInvalidUUIDFormat   = errors.New("invalid uuid format")
	ErrUserNotFound        = errors.New("user not found")
)

var GeneralErrors = []error{
	ErrInternalServerError,
	ErrSQLError,
	ErrTooManyRequest,
	ErrInvalidToken,
	ErrForbidden,
	ErrUnauthorized,
	ErrInvalidUploadFile,
	ErrFileTooLarge,
	ErrInvalidUUIDFormat,
	ErrUserNotFound,
}
