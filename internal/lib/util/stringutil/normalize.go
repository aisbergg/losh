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

package stringutil

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

// NormalizeName normalizes a name by performing a NFKD normalization, removing
// whitespaces and some other characters and lowercasing the result.
func NormalizeName(name string) string {
	transformer := transform.Chain(norm.NFKD, runes.Remove(runes.In(unicode.Mn)), runes.Remove(runes.In(unicode.White_Space)))
	normalized, _, _ := transform.String(transformer, name)
	normalized = strings.ToLower(normalized)
	return normalized
}
