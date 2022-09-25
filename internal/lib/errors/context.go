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

package errors

import "github.com/aisbergg/go-errors/pkg/errors"

// Contexter is an error that includes additional context in form of a map.
type Contexter interface {
	error
	Context() map[string]interface{}
}

// ContextAdder is an error that can add context to itself.
type ContextAdder interface {
	error
	Add(key string, value interface{}) ContextAdder
	AddAll(context map[string]interface{}) ContextAdder
}

// GetContext returns all the context of the error.
func GetContext(err error) map[string]interface{} {
	ctx := make(map[string]interface{})
	if e, ok := err.(Contexter); ok {
		ctx = e.Context()
	}
	return ctx
}

// GetFullContext returns the context of the whole error chain.
func GetFullContext(err error) map[string]interface{} {
	if err == nil {
		return make(map[string]interface{})
	}
	ctx := GetFullContext(errors.Unwrap(err))
	if e, ok := err.(Contexter); ok {
		for k, v := range e.Context() {
			ctx[k] = v
		}
	}
	return ctx
}
