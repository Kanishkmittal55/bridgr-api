package apierrors

import (
	"errors"
	"net/http"
)

var (
	ErrBadRequest         = errors.New("bad request")
	ErrForbidden          = errors.New("forbidden")
	ErrNotFound           = errors.New("not found")
	ErrTooManyRequests    = errors.New("too many requests")
	ErrInternal           = errors.New("internal error")
	ErrNotImplemented     = errors.New("not implemented")
	ErrServiceUnavailable = errors.New("service unavailable")
)

const InternalErrorString = "internal server error"

func HTTPErrorFromStatusCode(statusCode int) error {
	switch {
	case statusCode == http.StatusBadRequest:
		return ErrBadRequest
	case statusCode == http.StatusNotFound:
		return ErrNotFound
	case statusCode == http.StatusForbidden:
		return ErrForbidden
	case statusCode == http.StatusTooManyRequests:
		return ErrTooManyRequests
	case statusCode == http.StatusInternalServerError:
		return ErrInternal
	case statusCode == http.StatusNotImplemented:
		return ErrNotImplemented
	case statusCode == http.StatusServiceUnavailable:
		return ErrServiceUnavailable
	default:
		return ErrInternal
	}
}
