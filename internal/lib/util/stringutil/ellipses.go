package stringutil

// Ellipses shortens a string up to a given length and adds an ellipsis at the
// end.
func Ellipses(s string, l int) string {
	if l <= 0 {
		return ""
	}
	sr := []rune(s)
	if len(sr) <= l {
		return s
	}
	if len(sr) >= l {
		sr = sr[:l-1]
	}
	// trim right
	for i := len(sr) - 1; i >= 0; i-- {
		if sr[i] == ' ' || sr[i] == '\n' || sr[i] == '\t' {
			sr = sr[:i]
		} else {
			break
		}
	}
	if len(sr) == l {
		sr = sr[:l-1]
	}
	sr = append(sr, '…')
	return string(sr)
}
