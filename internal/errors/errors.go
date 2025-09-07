package errors

import "errors"

var (
	ErrUserAlreadyExists  = errors.New("user already exists")
	ErrorUserNotFound     = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrActivityNotFound   = errors.New("activity not found")
)
