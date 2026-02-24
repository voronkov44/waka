package users

import "errors"

var (
	ErrNotFound        = errors.New("user not found")
	ErrInvalidArgument = errors.New("invalid argument")
)
