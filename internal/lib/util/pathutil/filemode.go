package pathutil

import (
	"io/fs"
	"strconv"
	"strings"
)

func ParseFileMode(s string) (fs.FileMode, error) {
	s = strings.TrimSpace(s)
	n, err := strconv.ParseUint(s, 8, 32)
	if err != nil {
		return 0, err
	}

	return fs.FileMode(n), nil
}
