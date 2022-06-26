package net

import (
	"errors"
	"net/http"
)

// ErrMaxRedirectsExceeded indicates that the maximum number of redirects has
// been exceeded.
var ErrMaxRedirectsExceeded = errors.New("redirect loop")

// NewRedirectHandler creates a handler that redirects for the given number
// of redirects.
func NewRedirectHandler(maxRedirects int) func(req *http.Request, via []*http.Request) error {
	return func(req *http.Request, via []*http.Request) error {
		if len(via) >= maxRedirects {
			return ErrMaxRedirectsExceeded
		}
		return nil
	}
}
