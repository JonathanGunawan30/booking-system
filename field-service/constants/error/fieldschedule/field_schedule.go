package error

import "errors"

var (
	ErrFieldScheduleNotFound  = errors.New("field schedule not found")
	ErrFieldScheduleExist     = errors.New("field schedule already exist")
	ErrFieldScheduleYearPast  = errors.New("field schedule year must not in the past")
	ErrFieldScheduleMonthPast = errors.New("field schedule month must not in the past")
)

var FieldScheduleErrors = []error{
	ErrFieldScheduleNotFound,
	ErrFieldScheduleExist,
	ErrFieldScheduleYearPast,
	ErrFieldScheduleMonthPast,
}
