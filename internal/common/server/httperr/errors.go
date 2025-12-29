package httperr

import (
	"errors"
	"sumni-finance-backend/internal/common/validator"
)

type ErrorType struct {
	t string
}

var (
	ErrorTypeUnknown        = ErrorType{"unknown"}
	ErrorTypeAuthorization  = ErrorType{"authorization"}
	ErrorTypeIncorrectInput = ErrorType{"incorrect-input"}
)

type SlugError struct {
	logMsg     string
	slug       string
	errorType  ErrorType
	wrappedErr error
}

func (s SlugError) Error() string {
	return s.logMsg
}

func (s SlugError) Slug() string {
	return s.slug
}

func (s SlugError) ErrorType() ErrorType {
	return s.errorType
}

func (s SlugError) Unwrap() error {
	return s.wrappedErr
}

func NewError(wrappedErr error, slug string) SlugError {
	if wrappedErr == nil {
		return NewUnknowError(errors.New("NewError called with nil error"), "internal-error")
	}

	var validErrs *validator.ErrorList
	if errors.As(wrappedErr, &validErrs) {
		return NewIncorrectInputError(wrappedErr, "validate-failed")
	}

	return NewUnknowError(wrappedErr, "unknown-error")
}

func NewUnknowError(wrappedErr error, slug string) SlugError {
	return SlugError{
		logMsg:     wrappedErr.Error(),
		slug:       slug,
		errorType:  ErrorTypeUnknown,
		wrappedErr: wrappedErr,
	}
}

func NewAuthorizationError(wrappedErr error, slug string) SlugError {
	return SlugError{
		logMsg:     wrappedErr.Error(),
		slug:       slug,
		errorType:  ErrorTypeAuthorization,
		wrappedErr: wrappedErr,
	}
}

func NewIncorrectInputError(wrappedErr error, slug string) SlugError {
	return SlugError{
		logMsg:     wrappedErr.Error(),
		slug:       slug,
		errorType:  ErrorTypeIncorrectInput,
		wrappedErr: wrappedErr,
	}
}
