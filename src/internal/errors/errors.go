package errors

import "errors"

var (
	ErrUserNotFound = errors.New("user not found")
	ErrInvalidInput = errors.New("invalid input")
	ErrInvalidEmail = errors.New("invalid email format")
)
