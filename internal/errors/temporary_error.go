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
