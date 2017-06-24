package httpd

import (
	"errors"
	"fmt"
)

var (
	ErrRouteNotFound    = errors.New("router: Route not found")
	ErrMethodNotAllowed = errors.New("router: Method not allowed")
)

type BadRequestError struct {
	Text string
}

func NewBadRequestError(text string) BadRequestError {
	return BadRequestError{fmt.Sprintf("Bad request: %s", text)}
}

func (r BadRequestError) Error() string {
	return r.Text
}
