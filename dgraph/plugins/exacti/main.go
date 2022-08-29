package main

import (
	"strings"
)

func Tokenizer() interface{} { return ExactCaseInsensitiveTokenizer{} }

// ExactCaseInsensitiveTokenizer is a tokenizer for Dgraph to enable case
// insensitive text matching.
//
// For more information, see:
// https://dgraph.io/docs/query-language/indexing-custom-tokenizers/
type ExactCaseInsensitiveTokenizer struct{}

func (ExactCaseInsensitiveTokenizer) Name() string     { return "exacti" }
func (ExactCaseInsensitiveTokenizer) Type() string     { return "string" }
func (ExactCaseInsensitiveTokenizer) Identifier() byte { return 0xe1 }

func (t ExactCaseInsensitiveTokenizer) Tokens(value interface{}) ([]string, error) {
	return []string{strings.ToLower(strings.TrimSpace(value.(string)))}, nil
}
