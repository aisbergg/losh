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

package dgraph

import losherrors "losh/internal/lib/errors"

// RepoErrorer is the error type for the repository errors.
type RepoErrorer interface {
	IsRepoError()
}

type RepoError struct {
	losherrors.AppError
}

// NewRepoError creates a new repository error.
func NewRepoError(format string, args ...interface{}) *RepoError {
	return WrapRepoError(nil, format, args...)
}

// WrapRepoError wraps an error into BaseError.
func WrapRepoError(err error, format string, args ...interface{}) *RepoError {
	return &RepoError{
		AppError: *losherrors.NewAppErrorWrap(err, format, args...),
	}
}

// IsRepoError implements the RepoError interface.
func (*RepoError) IsRepoError() {}

// IsRepoError indicates if the error is a repository error.
func IsRepoError(err error) bool {
	_, ok := err.(RepoErrorer)
	return ok
}
