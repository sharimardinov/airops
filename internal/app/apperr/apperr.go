package apperr

import (
	"errors"
	"fmt"
)

type Kind string

const (
	KindValidation Kind = "validation"
	KindNotFound   Kind = "not_found"
	KindInternal   Kind = "internal"
)

type Error struct {
	Kind    Kind
	Message string
	Details map[string]any
	Err     error
}

func (e *Error) Error() string {
	if e.Message != "" {
		return e.Message
	}
	if e.Err != nil {
		return e.Err.Error()
	}
	return "error"
}

func (e *Error) Unwrap() error { return e.Err }

func Validation(message string, details map[string]any) *Error {
	return &Error{Kind: KindValidation, Message: message, Details: details}
}

func NotFound(message string, err error) *Error {
	return &Error{Kind: KindNotFound, Message: message, Err: err}
}

func Internal(message string, err error) *Error {
	return &Error{Kind: KindInternal, Message: message, Err: err}
}

func As(err error) (*Error, bool) {
	var ae *Error
	if errors.As(err, &ae) {
		return ae, true
	}
	return nil, false
}

func Wrap(kind Kind, message string, err error) *Error {
	return &Error{Kind: kind, Message: message, Err: fmt.Errorf("%w", err)}
}
