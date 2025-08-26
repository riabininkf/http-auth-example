package auth

import "errors"

var (
	ErrEmailBusy = errors.New("email is busy")

	ErrUserNotFound = errors.New("user not found")
)
