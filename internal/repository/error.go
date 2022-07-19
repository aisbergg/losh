package repository

import losherrors "losh/internal/errors"

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
