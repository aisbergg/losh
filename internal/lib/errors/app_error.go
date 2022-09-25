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

import (
	"time"

	"github.com/aisbergg/go-errors/pkg/errors"
)

// IsAppError indicates if the error is a LOSH App Error.
func IsAppError(err error) bool {
	_, ok := err.(interface{ IsAppError() })
	return ok
}

// AppError is an error that is temporary and can be retried.
type AppError struct {
	errors.TraceableError
	timestamp time.Time
	context   map[string]interface{}
}

// NewAppError creates a new AppError.
func NewAppError(format string, args ...interface{}) *AppError {
	timestamp := time.Now()
	return &AppError{
		TraceableError: *errors.ErrorfSkip(1, format, args...),
		timestamp:      timestamp,
		context:        map[string]interface{}{"timestamp": timestamp},
	}
}

// NewAppErrorWrap wraps an error into AppError.
func NewAppErrorWrap(err error, format string, args ...interface{}) *AppError {
	timestamp := time.Now()
	return &AppError{
		TraceableError: *errors.WrapfSkip(err, 1, format, args...),
		timestamp:      timestamp,
		context:        map[string]interface{}{"timestamp": timestamp},
	}
}

// IsAppError implements the AppErrorer interface.
func (*AppError) IsAppError() {}

// Timestamp returns the craetion timestamp of the error.
func (e *AppError) Timestamp() time.Time {
	return e.timestamp
}

// Add adds a key-value pair to the error context. If the key already exists, it will be overwritten. If the value is nil, the key will be removed.
func (e *AppError) Add(key string, value interface{}) *AppError {
	if value == nil {
		delete(e.context, key)
		return e
	}
	e.context[key] = value
	return e
}

// AddAll adds all key-value pairs to the error context.
func (e *AppError) AddAll(context map[string]interface{}) *AppError {
	for key, value := range context {
		e.Add(key, value)
	}
	return e
}

// Context returns the context for the error.
func (e *AppError) Context() map[string]interface{} {
	return e.context
}

// FullContext returns the context of the whole error chain.
func (e *AppError) FullContext() map[string]interface{} {
	return GetFullContext(e)
}
