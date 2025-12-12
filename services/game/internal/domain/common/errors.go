package common

import "errors"

var (
	ErrNotFound    = errors.New("not found")
	ErrInvalidData = errors.New("invalid data")
	ErrUnauthorized = errors.New("unauthorized")
)

