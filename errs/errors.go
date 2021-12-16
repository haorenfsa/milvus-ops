package errs

import (
	"errors"
	"net/http"
)

// common errs
var (
	ErrNotFound = errors.New("resource not found")
	ErrBadRequest = errors.New("err bad request")
	ErrStorage    = errors.New("storage error")
)

// ErrorToHTTPCode ..
func ErrorToHTTPCode(err error) int {
	if errors.Is(err, ErrNotFound) {
		return http.StatusNotFound
	}
	if errors.Is(err, ErrBadRequest) {
		return http.StatusBadRequest
	}
	return http.StatusInternalServerError
}
