package httperr

import (
	"errors"
	"net/http"
	"sumni-finance-backend/internal/common/logs"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

func InternalError(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, err.Error(), http.StatusInternalServerError)
}

func Unauthorised(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, err.Error(), http.StatusUnauthorized)
}

func BadRequest(slug string, err error, w http.ResponseWriter, r *http.Request) {
	httpRespondWithError(err, slug, w, r, err.Error(), http.StatusBadRequest)
}

func RespondWithSlugError(err error, w http.ResponseWriter, r *http.Request) {
	slugError, ok := err.(SlugError)
	if !ok {
		InternalError("internal-server-error", err, w, r)
		return
	}

	switch slugError.ErrorType() {
	case ErrorTypeAuthorization:
		Unauthorised(slugError.Slug(), slugError, w, r)
	case ErrorTypeIncorrectInput:
		BadRequest(slugError.Slug(), slugError, w, r)
	default:
		InternalError(slugError.Slug(), slugError, w, r)
	}
}

func httpRespondWithError(err error, slug string, w http.ResponseWriter, r *http.Request, logMSg string, status int) {
	wrappedErr := errors.Unwrap(err)
	if wrappedErr == nil {
		wrappedErr = err // If no wrapping occurred, use the error itself
	}

	logger := logs.FromContext(r.Context()).With(
		"error", wrappedErr,
		"error_slug", slug,
	)

	requestID := middleware.GetReqID(r.Context())

	// Correct severity based on HTTP status
	if status >= 500 {
		logger.Error(logMSg)
	} else {
		logger.Warn(logMSg)
	}

	resp := ErrorResponse{
		RequestID: requestID,
		Error: ErrorDetail{
			Slug:    slug,
			Message: logMSg,
		},
		httpStatus: status,
	}

	if err := render.Render(w, r, resp); err != nil {
		logger.Error("failed to render resp: " + err.Error())
	}
}

type ErrorDetail struct {
	Slug    string `json:"slug"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	RequestID  string      `json:"request_id"`
	Error      ErrorDetail `json:"error"`
	httpStatus int
}

func (e ErrorResponse) Render(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(e.httpStatus)
	return nil
}
