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
	"bytes"
	"errors"
	"fmt"
	"path"
	"path/filepath"
	"regexp"

	"github.com/alecthomas/chroma"
	"github.com/alecthomas/chroma/formatters/html"
	"github.com/alecthomas/chroma/lexers"
	"github.com/alecthomas/chroma/styles"
	"github.com/osteele/liquid"
	"github.com/osteele/liquid/render"
)

// -----------------------------------------------------------------------------
//
// https://github.com/osteele/gojekyll/blob/5ff68e140697bcb30660f083ad42dc328578b16f/tags/tags.go
//
// -----------------------------------------------------------------------------

// A LinkTagHandler given an include tag file name returns a URL.
type LinkTagHandler func(string) (string, bool)

// AddJekyllTags adds the Jekyll tags to the Liquid engine.
func addJekyllTags(e *liquid.Engine, templates map[string]*liquid.Template) {
	tc := tagContext{templates}
	e.RegisterBlock("highlight", highlightTag)
	e.RegisterTag("include", tc.includeTag)
	e.RegisterTag("include_relative", tc.includeRelativeTag)
}

// tagContext provides the context to a tag renderer.
type tagContext struct {
	templates map[string]*liquid.Template
}

// -----------------------------------------------------------------------------
//
// https://github.com/osteele/gojekyll/blob/5ff68e140697bcb30660f083ad42dc328578b16f/tags/highlight.go
//
// -----------------------------------------------------------------------------

var highlightArgsRE = regexp.MustCompile(`^\s*(\S+)(\s+linenos)?\s*$`)

func highlightTag(rc render.Context) (string, error) {
	argStr, err := rc.ExpandTagArg()
	if err != nil {
		return "", err
	}
	args := highlightArgsRE.FindStringSubmatch(argStr)
	if args == nil {
		return "", fmt.Errorf("syntax error")
	}
	source, err := rc.InnerString()
	if err != nil {
		return "", err
	}

	// Determine lexer.
	l := lexers.Get(args[1])
	if l == nil {
		l = lexers.Analyse(source)
	}
	if l == nil {
		l = lexers.Fallback
	}
	l = chroma.Coalesce(l)

	lineNum := args[2] != ""

	// Determine formatter.
	f := html.New(
		html.WithClasses(true),
		html.WithLineNumbers(lineNum),
		html.LineNumbersInTable(true),
	)

	// Determine style.
	s := styles.Get("")
	if s == nil {
		s = styles.Fallback
	}

	it, err := l.Tokenise(nil, source)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = f.Format(buf, s, it); err != nil {
		return "", err
	}
	return buf.String(), nil
}

// -----------------------------------------------------------------------------
//
// reimplemented
//
// -----------------------------------------------------------------------------

func (tc tagContext) includeTag(rc render.Context) (s string, err error) {
	path.Join()
	return includeFromDir("includes", tc.templates, rc)
}

func (tc tagContext) includeRelativeTag(rc render.Context) (string, error) {
	return includeFromDir(path.Dir(rc.SourceFile()), tc.templates, rc)
}

func includeFromDir(dir string, templates map[string]*liquid.Template, rc render.Context) (string, error) {
	// parse args and options
	argsline, err := rc.ExpandTagArg()
	if err != nil {
		return "", err
	}
	args, err := ParseArgs(argsline)
	if err != nil {
		return "", err
	}
	if len(args.Args) != 1 {
		return "", fmt.Errorf("failed to parse args")
	}
	include, err := args.EvalOptions(rc)
	if err != nil {
		return "", err
	}
	filename := filepath.Clean("/" + args.Args[0]) // remove any '../'
	filename = filepath.Join(dir, filename)

	// get template
	template, ok := templates[filename]
	if !ok {
		return "", NewRenderError(filename, errors.New("cannot locate file"), rc.Bindings())
	}

	// copy lexical environment and add 'include'
	bindings := map[string]interface{}{}
	for k, v := range rc.Bindings() {
		bindings[k] = v
	}
	bindings["include"] = include

	// render template
	rendered, err := template.Render(bindings)
	if err != nil {
		return "", NewRenderError(filename, err, bindings)
	}
	return string(rendered), nil
}
