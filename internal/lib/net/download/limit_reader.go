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
