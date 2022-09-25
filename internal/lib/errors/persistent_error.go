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

// IsPersistentError returns true if err is of type PersistentError.
func IsPersistentError(err error) bool {
	_, ok := err.(interface{ IsPersistentError() })
	return ok
}

// PersistentError is an error that is temporary and can be retried.
type PersistentError struct {
	*AppError
}

// NewPersistentError wraps an error into PersistentError.
func NewPersistentError(err error, format string, args ...interface{}) error {
	terr := &PersistentError{
		AppError: NewAppErrorWrap(err, format, args...),
	}
	return terr
}

// IsPersistentError indicates the type of the error.
func (e *PersistentError) IsPersistentError() {}
