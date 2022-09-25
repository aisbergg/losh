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
)

// IsTemporaryError returns true if err is of type TemporaryError.
func IsTemporaryError(err error) bool {
	_, ok := err.(interface{ IsTemporaryError() })
	return ok
}

// TemporaryError is an error that is temporary and can be retried.
type TemporaryError struct {
	AppError
	RetryAfter time.Time
	retries    int
}

// NewTemporaryError wraps an error into TemporaryError.
func NewTemporaryError(err error, retryAfter time.Time, format string, args ...interface{}) error {
	tmpErr := &TemporaryError{
		AppError:   *NewAppErrorWrap(err, format, args...),
		RetryAfter: retryAfter,
	}
	tmpErr.Add("retry_after", retryAfter)
	return tmpErr
}

// IsTemporaryError indicates the type of the error.
func (e *TemporaryError) IsTemporaryError() {}
