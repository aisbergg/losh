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
	"bytes"
	"flag"
	"fmt"
	"go/format"
	"html"
	"os"
	"text/template"

	"github.com/aisbergg/go-errors/pkg/errors"

	"github.com/Masterminds/sprig"
	"github.com/aisbergg/go-pathlib/pkg/pathlib"
	"gopkg.in/yaml.v3"
)

func main() {
	force := flag.Bool("f", false, "Force overwrite of rendered files")
	noFormat := flag.Bool("n", false, "Disable formatting")
	flag.Usage = func() {
		fmt.Printf("Usage: %s [OPTIONS] CONTEXT_FILE\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(2)
	}
	flag.Parse()

	ctxFle := flag.Arg(0)
	if ctxFle == "" {
		flag.Usage()
	}

	err := run(ctxFle, *force, *noFormat)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("Done")
}

func run(contextFile string, force, noFormat bool) error {
	ctxFlePth := pathlib.NewPath(contextFile)
	ctx, err := loadContext(contextFile)
	if err != nil {
		return errors.Wrapf(err, "failed to read context file '%s'", contextFile)
	}

	for _, file := range ctx.Generate {
		srcStr := file.Src
		if srcStr == "" {
			srcStr = ctx.Defaults.Src
		}
		if srcStr == "" {
			return errors.New("missing template src")
		}
		srcPth := pathlib.NewPath(srcStr)
		// make source file path relative to context file path
		if !srcPth.IsAbsolute() {
			srcPth = ctxFlePth.Parent().JoinPath(srcPth)
		}

		var dstPth pathlib.Path
		if file.Dest == "" {
			dstPth, err = srcPth.WithSuffix(".go")
			if err != nil {
				return errors.Wrapf(err, "failed to create destination path for template file '%s'", srcStr)
			}
		} else {
			dstPth = pathlib.NewPath(file.Dest)
		}
		// make template file path relative to context file path
		if !dstPth.IsAbsolute() {
			dstPth = ctxFlePth.Parent().JoinPath(dstPth)
		}

		vars := merge(ctx.Defaults.Vars, file.Vars)
		err = generateFile(srcPth, dstPth, vars, force, noFormat)
		if err != nil {
			return errors.Wrapf(err, "failed to generate file '%s'", dstPth.String())
		}
	}

	return nil
}

func merge(x, y map[string]interface{}) map[string]interface{} {
	merged := make(map[string]interface{})
	for k, v := range x {
		merged[k] = v
	}
	for k, v := range y {
		ev, ok := merged[k]
		if ok {
			if ev, ok := ev.(map[string]interface{}); ok {
				if v, ok := v.(map[string]interface{}); ok {
					merged[k] = merge(ev, v)
					continue
				}
			}
		}
		// set/replace value
		merged[k] = v
	}
	return merged
}

type Context struct {
	Generate []File   `json:"generate"`
	Defaults Defaults `json:"defaults"`
}

type Defaults struct {
	Src  string                 `json:"src"`
	Vars map[string]interface{} `json:"vars"`
}

type File struct {
	Src  string                 `json:"src"`
	Dest string                 `json:"dest"`
	Vars map[string]interface{} `json:"vars"`
}

func loadContext(pathStr string) (*Context, error) {
	path := pathlib.NewPath(pathStr)
	cnt, err := path.ReadFile()
	if err != nil {
		return nil, err
	}
	ctx := &Context{}
	err = yaml.Unmarshal(cnt, ctx)
	return ctx, err
}

func generateFile(tplPth, dstPth pathlib.Path, ctx map[string]interface{}, force, noFormat bool) error {
	fmt.Printf("Generating %s --> %s\n", tplPth.String(), dstPth.String())

	// create output directory if it doesn't exist
	exists, err := dstPth.Exists()
	if err != nil {
		return err
	}
	if exists && !force {
		return errors.Errorf("file '%s' already exists. Use -f to forcefully overwrite it", dstPth.String())
	}
	err = dstPth.Parent().MkdirAll()
	if err != nil {
		return errors.Wrapf(err, "failed to create directory '%s'", dstPth.Parent().String())
	}

	// render tpl
	tpl, err := template.New(tplPth.Name()).
		Funcs(template.FuncMap(additionalTemplateFuncs)).
		Funcs(sprig.TxtFuncMap()).
		ParseFiles(tplPth.String())
	if err != nil {
		return errors.Wrapf(err, "failed to load template file '%s'", tplPth.String())
	}
	var buf bytes.Buffer
	err = tpl.Execute(&buf, ctx)
	if err != nil {
		return errors.Wrapf(err, "failed to render template file '%s'", tplPth.String())
	}
	rendered := buf.Bytes()

	// format code
	if !noFormat {
		rendered, err = format.Source(rendered)
		if err != nil {
			return errors.Wrapf(err, "failed to format file '%s'", dstPth.String())
		}
	}

	// write file
	dstFile, err := dstPth.OpenFile(os.O_CREATE | os.O_TRUNC | os.O_RDWR)
	defer dstFile.Close()
	_, err = dstFile.Write(rendered)
	if err != nil {
		return errors.Wrapf(err, "failed to write file '%s'", dstPth.String())
	}

	return nil
}

var additionalTemplateFuncs = map[string]interface{}{
	"escapeHTML": escapeHTML,
}

func escapeHTML(name string) string {
	return html.EscapeString(name)
}
