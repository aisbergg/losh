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
