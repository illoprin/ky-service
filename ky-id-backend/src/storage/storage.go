package storage

import (
	"errors"
)

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidField    = errors.New("invalid field")
	ErrInvalidPassword = errors.New("invalid password")
	ErrSameLoginExist  = errors.New("user with same login exists")
)
