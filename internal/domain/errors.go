package domain

import "errors"

var (
	ErrNotFound        = errors.New("not found")
	ErrInvalidArgument = errors.New("invalid argument")
	ErrBadRequest      = errors.New("bad request")
)
