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
