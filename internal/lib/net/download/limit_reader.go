package download

import (
	"errors"
	"io"
	"losh/internal/lib/unit"
)

var errLimitExceeded = errors.New("limit exceeded")

// A limitedReader is just like io.LimitedReader but it returns a ErrTooLarge
// instead of EOF, if the max size is exceeded.
type limitedReader struct {
	R io.Reader     // underlying reader
	N unit.ByteSize // max bytes remaining
}

// Read reads data from the limitedReader.
func (l *limitedReader) Read(p []byte) (n int, err error) {
	if l.N <= 0 {
		return 0, errLimitExceeded
	}
	if unit.ByteSize(len(p)) > l.N {
		p = p[0:l.N]
	}
	n, err = l.R.Read(p)
	l.N -= unit.ByteSize(n)
	return
}
