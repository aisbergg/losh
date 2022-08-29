package search

type ErrorType int

const (
	ErrorTypeInvalid ErrorType = iota
	ErrorLimitExceeded
	ErrorInvalidQuery
)

type Error struct {
	Msg  string
	Type ErrorType
}

func (e *Error) Error() string {
	return e.Msg
}
