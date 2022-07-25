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
