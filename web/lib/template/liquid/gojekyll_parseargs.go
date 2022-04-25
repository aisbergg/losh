// MIT License
//
// Copyright (c) 2017 Oliver Steele
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package liquid

import (
	"fmt"
	"regexp"

	"github.com/osteele/liquid/render"
)

// -----------------------------------------------------------------------------
//
// https://github.com/osteele/gojekyll/blob/5ff68e140697bcb30660f083ad42dc328578b16f/tags/parseargs.go
//
// -----------------------------------------------------------------------------

var argPattern = regexp.MustCompile(`^([^=\s]+)(?:\s+|$)`)
var optionPattern = regexp.MustCompile(`^([\w-]+)=("[^"]*"|'[^']*'|[^'"\s]*)(?:\s+|$)`)

// ParsedArgs holds the parsed arguments from ParseArgs.
type ParsedArgs struct {
	Args    []string
	Options map[string]optionRecord
}

type optionRecord struct {
	value  string
	quoted bool
}

// ParseArgs parses a tag argument line {% include arg1 arg2 opt=a opt2='b' %}
func ParseArgs(argsline string) (*ParsedArgs, error) {
	args := ParsedArgs{
		[]string{},
		map[string]optionRecord{},
	}
	// Ranging over FindAllStringSubmatch would be better golf but got out of hand
	// maintenance-wise.
	for r, i := argsline, 0; len(r) > 0; r = r[i:] {
		am := argPattern.FindStringSubmatch(r)
		om := optionPattern.FindStringSubmatch(r)
		switch {
		case am != nil:
			args.Args = append(args.Args, am[1])
			i = len(am[0])
		case om != nil:
			k, v, quoted := om[1], om[2], false
			if v[0] == '\'' || v[0] == '"' {
				v, quoted = v[1:len(v)-1], true
			}
			args.Options[k] = optionRecord{v, quoted}
			i = len(om[0])
		default:
			return nil, fmt.Errorf("parse error in tag parameters %q", argsline)
		}
	}
	return &args, nil
}

// EvalOptions evaluates unquoted options.
func (r *ParsedArgs) EvalOptions(ctx render.Context) (map[string]interface{}, error) {
	options := map[string]interface{}{}
	for k, v := range r.Options {
		if v.quoted {
			options[k] = v.value
		} else {
			value, err := ctx.EvaluateString(v.value)
			if err != nil {
				return nil, err
			}
			options[k] = value
		}
	}
	return options, nil
}
