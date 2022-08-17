package stringutil

import (
	"regexp"
	"strings"

	"github.com/aisbergg/go-unidecode/pkg/unidecode"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

var (
	invalidCharPattern  = regexp.MustCompile(`[^\w^-]`)
	dupSeparatorPattern = regexp.MustCompile(`-{2,}`)
)

// Slugify returns a slugified version of the given string.
//
// Example:
//   Slugify("kožušček hello world") // Output: kozuscek-hello-world
func Slugify(s string) string {
	s, _, _ = transform.String(norm.NFKD, s)
	s, _ = unidecode.Unidecode(s, unidecode.Ignore)
	s = invalidCharPattern.ReplaceAllString(s, "-")
	s = dupSeparatorPattern.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	s = strings.ToLower(s)
	return s
}
