package httperr

type ErrorType struct {
	t string
}

var (
	ErrorTypeUnknown        = ErrorType{"unknown"}
	ErrorTypeAuthorization  = ErrorType{"authorization"}
	ErrorTypeIncorrectInput = ErrorType{"incorrect-input"}
)

type SlugError struct {
	logMsg    string
	slug      string
	errorType ErrorType
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

func NewUnknowError(logMsg string, slug string) SlugError {
	return SlugError{
		logMsg:    logMsg,
		slug:      slug,
		errorType: ErrorTypeUnknown,
	}
}

func NewAuthorizationError(logMsg string, slug string) SlugError {
	return SlugError{
		logMsg:    logMsg,
		slug:      slug,
		errorType: ErrorTypeAuthorization,
	}
}

func NewIncorrectInputError(logMsg string, slug string) SlugError {
	return SlugError{
		logMsg:    logMsg,
		slug:      slug,
		errorType: ErrorTypeIncorrectInput,
	}
}
