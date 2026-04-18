package constatns

import "errors"

var (
	ErrUserNotFound                = errors.New("user not found")
	ErrUsernameOrPasswordIncorrect = errors.New("username or password is incorrect")
	ErrPasswordIncorrect           = errors.New("password incorrect")
	ErrUsernameExists              = errors.New("username already exists")
	ErrEmailExists                 = errors.New("email already exists")
	ErrPasswordDoesNotMatch        = errors.New("password doest not match")
)

var UserErrors = []error{
	ErrUserNotFound,
	ErrPasswordIncorrect,
	ErrUsernameExists,
	ErrPasswordDoesNotMatch,
	ErrEmailExists,
	ErrUsernameOrPasswordIncorrect,
}
