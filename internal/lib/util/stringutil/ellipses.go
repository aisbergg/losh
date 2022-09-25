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
