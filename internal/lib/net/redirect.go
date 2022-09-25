// Copyright 2022 André Lehmann
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
