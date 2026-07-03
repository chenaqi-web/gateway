package user

import "errors"

var (
	// ErrNotFound is returned when a user aggregate does not exist.
	ErrNotFound = errors.New("user not found")
)
