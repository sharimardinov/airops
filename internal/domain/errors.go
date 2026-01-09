package domain

import "errors"

var (
	ErrNotFound           = errors.New("not found")
	ErrInvalidArgument    = errors.New("invalid argument")
	ErrFlightNotAvailable = errors.New("flight not available")
	ErrSeatNotAvailable   = errors.New("seat not available")
)
