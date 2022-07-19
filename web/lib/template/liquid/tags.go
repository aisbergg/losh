package liquid

import (
	"errors"
	"regexp"
	"strings"

	"github.com/osteele/liquid"
	"github.com/osteele/liquid/render"
)

func addTags(e *liquid.Engine, templates map[string]*liquid.Template) {
	addJekyllTags(e, templates)

	// tabler.io tags
	e.RegisterBlock("removeemptylines", removeEmptyLinesBlock)
	e.RegisterBlock("hide", hideBlock)
	e.RegisterBlock("card", cardBlock)
	e.RegisterBlock("capture_global", captureGlobalBlock)
}

// -----------------------------------------------------------------------------
//
// required for tabler.io
//
// -----------------------------------------------------------------------------

var removeEmptyLinesPattern = regexp.MustCompile(`^\s*$\n`)

func removeEmptyLinesBlock(rc render.Context) (string, error) {
	var content string
	var err error
	if content, err = rc.InnerString(); err != nil {
		return "", err
	}
	content = strings.TrimSpace(content)
	content = removeEmptyLinesPattern.ReplaceAllString(content, "")
	return content, nil
}

func hideBlock(rc render.Context) (string, error) {
	var content string
	var err error
	if content, err = rc.InnerString(); err != nil {
		return "", err
	}
	return "{% hide %}" + content + "{% endhide %}", nil
}

func cardBlock(rc render.Context) (string, error) {
	var content string
	var title string
	var err error
	if content, err = rc.InnerString(); err != nil {
		return "", err
	}
	if title, err = rc.ExpandTagArg(); err != nil {
		return "", err
	}
	title = strings.TrimSpace(title)
	builder := strings.Builder{}
	builder.Grow(len(content) + len(title) + 120) // fit all the contents
	builder.WriteString(`<div class="card">`)
	if len(title) > 0 {
		builder.WriteString(`<div class="card-header"><div class="card-title">`)
		builder.WriteString(title)
		builder.WriteString(`</div></div>`)
	}
	builder.WriteString(`<div class="card-body">`)
	builder.WriteString(content)
	builder.WriteString(`</div></div>`)
	return builder.String(), nil
}

func captureGlobalBlock(rc render.Context) (string, error) {
	var name string
	var content string
	var err error
	if name, err = rc.ExpandTagArg(); err != nil {
		return "", err
	}
	if content, err = rc.InnerString(); err != nil {
		return "", err
	}
	captureGlobal, ok := rc.Get("capture_global").(map[string][]string)
	if !ok {
		return "", errors.New("capture_global has a wrong type, did you overwrite it?")
	}
	captured, ok := captureGlobal[name]
	if ok {
		captured = append(captured, content)
	} else {
		captured = []string{content}
	}
	captureGlobal[name] = captured

	// block returns nothing
	return "", nil
}
