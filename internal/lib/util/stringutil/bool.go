package stringutil

import "strings"

// ParseBool returns the boolean value represented by the string.
func ParseBool(str string) (b bool, ok bool) {
	str = strings.ToLower(str)
	switch str {
	case "1", "t", "true", "yes", "y":
		return true, true
	case "0", "f", "false", "no", "n":
		return false, true
	}
	return false, false
}
