package stringutil

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// NormalizeName normalizes a license name by performing a NFD normalization,
// removing whitespaces and some other characters and lowercasing the result.
func NormalizeName(name string) string {
	transformer := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), runes.Remove(runes.In(unicode.White_Space)))
	normalized, _, _ := transform.String(transformer, name)
	normalized = strings.ToLower(normalized)
	return normalized
}
