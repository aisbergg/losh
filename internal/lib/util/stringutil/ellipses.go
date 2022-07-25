package stringutil

import (
	"strings"
)

// Ellipses shortens a string up to a given length and adds an ellipsis at the
// end.
func Ellipses(s string, l int) string {
	if l <= 0 {
		return ""
	}
	if len(s) <= l {
		return s
	}
	s = strings.TrimRight(s[0:(l-1)], " \n\t")
	return s[0:(l-1)] + "…"
}
