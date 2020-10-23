package httperror

import (
	"net/http"
	"runtime/debug"

	"github.com/go-errors/errors"
	"github.com/go-openapi/runtime"
	"github.com/go-openapi/runtime/middleware"
)

// HTTPError represents an error that occurred while handling a request.
type HTTPError struct {
	Code        int    `json:"-"`
	Description string `json:"description"`
	Inner       error  `json:"inner"`
	Stacktrace  string `json:"stacktrace"`
}

// Error makes it compatible with `error` interface.
func (he *HTTPError) Error() string {
	return he.Description
}

// NewHTTPError creates a new HTTPError instance.
func NewHTTPError(code int, description string, inner ...error) *HTTPError {
	he := &HTTPError{
		Code:        code,
		Description: description,
	}
	if len(inner) > 0 {
		he.Inner = inner[0]
		he.Description = he.Description + ": " + he.Inner.Error()
		if withTrace, ok := inner[0].(*errors.Error); ok {
			he.Stacktrace = withTrace.ErrorStack()
		}
	}
	if len(he.Stacktrace) == 0 {
		he.Stacktrace = string(debug.Stack())
	}

	return he
}

func (he *HTTPError) WriteResponse(rw http.ResponseWriter, producer runtime.Producer) {
	if he.Code > 0 {
		rw.WriteHeader(he.Code)
	} else {
		rw.WriteHeader(http.StatusInternalServerError)
	}
	if err := producer.Produce(rw, he); err != nil {
		panic(err)
	}
}

func NewNotFoundError(description string) *HTTPError {
	return NewHTTPError(http.StatusNotFound, description)
}

func NewConflictError(description string) *HTTPError {
	return NewHTTPError(http.StatusConflict, description)
}

func NewForbiddenError(description string) *HTTPError {
	return NewHTTPError(http.StatusForbidden, description)
}

func NewBadRequestError(description string) *HTTPError {
	return NewHTTPError(http.StatusBadRequest, description)
}

func NewRequestEntityTooLargeError(description string) *HTTPError {
	return NewHTTPError(http.StatusRequestEntityTooLarge, description)
}

func NewServerError(description string) *HTTPError {
	return NewHTTPError(http.StatusInternalServerError, description)
}

func ConvertHTTPErrorToResponse(err error) middleware.Responder {
	if httpError, ok := err.(*HTTPError); ok {
		return httpError

	}
	return NewServerError(err.Error())
}
