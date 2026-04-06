package errors

import (
	stderrors "errors"
	"net/http"
)

var (
	ErrBadRequest         = stderrors.New("bad request")
	ErrConflict           = stderrors.New("conflict")
	ErrForbidden          = stderrors.New("forbidden")
	ErrNotFound           = stderrors.New("not found")
	ErrTooManyRequests    = stderrors.New("too many requests")
	ErrInternal           = stderrors.New("internal error")
	ErrNotImplemented     = stderrors.New("not implemented")
	ErrServiceUnavailable = stderrors.New("service unavailable")
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
