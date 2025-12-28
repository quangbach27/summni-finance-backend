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

func NewError(err error, slug string) SlugError {
	if err == nil {
		return NewUnknowError(errors.New("NewError called with nil error"), "internal-error")
	}

	var validErrs *validator.ErrorList
	if errors.As(err, &validErrs) {
		return NewIncorrectInputError(err, "validate-failed")
	}

	return NewUnknowError(err, "unknown-error")
}

func NewUnknowError(err error, slug string) SlugError {
	return SlugError{
		logMsg:     err.Error(),
		slug:       slug,
		errorType:  ErrorTypeUnknown,
		wrappedErr: err,
	}
}

func NewAuthorizationError(err error, slug string) SlugError {
	return SlugError{
		logMsg:     err.Error(),
		slug:       slug,
		errorType:  ErrorTypeAuthorization,
		wrappedErr: err,
	}
}

func NewIncorrectInputError(err error, slug string) SlugError {
	return SlugError{
		logMsg:     err.Error(),
		slug:       slug,
		errorType:  ErrorTypeIncorrectInput,
		wrappedErr: err,
	}
}
