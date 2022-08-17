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
	"encoding/json"
	"encoding/xml"
	"fmt"
	"math/rand"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"

	"losh/internal/lib/util/stringutil"

	"github.com/osteele/liquid"
	"github.com/osteele/liquid/evaluator"
	"github.com/osteele/liquid/expressions"
	"github.com/russross/blackfriday/v2"
)

// -----------------------------------------------------------------------------
//
// https://github.com/osteele/gojekyll/blob/5ff68e140697bcb30660f083ad42dc328578b16f/filters/filters.go
//
// -----------------------------------------------------------------------------

// AddJekyllFilters adds the Jekyll filters to the Liquid engine.
func addJekyllFilters(e *liquid.Engine) {
	// array filters
	e.RegisterFilter("array_to_sentence_string", arrayToSentenceStringFilter)
	e.RegisterFilter("filter", func(values []map[string]interface{}, key string) []interface{} {
		var result []interface{}
		for _, value := range values {
			if _, ok := value[key]; ok {
				result = append(result, value)
			}
		}
		return result
	})
	e.RegisterFilter("group_by", groupByFilter)
	e.RegisterFilter("group_by_exp", groupByExpFilter)
	e.RegisterFilter("sample", func(array []interface{}) interface{} {
		if len(array) == 0 {
			return nil
		}
		return array[rand.Intn(len(array))]
	})
	// sort overrides the Liquid filter with one that takes parameters
	e.RegisterFilter("sort", sortFilter)
	e.RegisterFilter("where", whereFilter)
	e.RegisterFilter("where_exp", whereExpFilter)
	e.RegisterFilter("xml_escape", xml.Marshal)
	e.RegisterFilter("push", func(array []interface{}, item interface{}) interface{} {
		return append(array, evaluator.MustConvertItem(item, array))
	})
	e.RegisterFilter("pop", requireNonEmptyArray(func(array []interface{}) interface{} {
		return array[0]
	}))
	e.RegisterFilter("shift", requireNonEmptyArray(func(array []interface{}) interface{} {
		return array[len(array)-1]
	}))
	e.RegisterFilter("unshift", func(array []interface{}, item interface{}) interface{} {
		return append([]interface{}{evaluator.MustConvertItem(item, array)}, array...)
	})

	// dates
	e.RegisterFilter("date_to_rfc822", func(date time.Time) string {
		return date.Format(time.RFC822)
		// Out: Mon, 07 Nov 2008 13:07:54 -0800
	})
	e.RegisterFilter("date_to_string", func(date time.Time) string {
		return date.Format("02 Jan 2006")
		// Out: 07 Nov 2008
	})
	e.RegisterFilter("date_to_long_string", func(date time.Time) string {
		return date.Format("02 January 2006")
		// Out: 07 November 2008
	})
	e.RegisterFilter("date_to_xmlschema", func(date time.Time) string {
		return date.Format("2006-01-02T15:04:05-07:00")
		// Out: 2008-11-07T13:07:54-08:00
	})

	e.RegisterFilter("jsonify", json.Marshal)
	e.RegisterFilter("markdownify", blackfriday.Run)
	e.RegisterFilter("normalize_whitespace", func(s string) string {
		// s = strings.Replace(s, "n", "N", -1)
		wsPattern := regexp.MustCompile(`(?s:[\s\n]+)`)
		return wsPattern.ReplaceAllString(s, " ")
	})
	rawPattern := regexp.MustCompile(`\s+`)
	defaultPattern := regexp.MustCompile(`[^[:alnum:]]+`)
	prettyPattern := regexp.MustCompile(`[^[:alnum:]\._~!$&'()+,;=@]+`)
	dupSeparatorPattern := regexp.MustCompile(`-{2,}`)
	e.RegisterFilter("slugify", func(s, mode string) string {
		replFunc := func(s string, pattern *regexp.Regexp) string {
			s = pattern.ReplaceAllString(s, "-")
			s = dupSeparatorPattern.ReplaceAllString(s, "-")
			s = strings.Trim(s, "-")
			s = strings.ToLower(s)
			return s
		}
		switch mode {
		case "none":
			return s
		case "raw":
			return replFunc(s, rawPattern)
		case "pretty":
			return replFunc(s, prettyPattern)
		case "latin":
			return stringutil.Slugify(s)
		case "default":
			fallthrough
		default:
			return replFunc(s, defaultPattern)
		}
	})
	e.RegisterFilter("to_integer", func(n int) int { return n })
	e.RegisterFilter("number_of_words", func(s string) int {
		wordPattern := regexp.MustCompile(`\w+`)
		m := wordPattern.FindAllStringIndex(s, -1)
		if m == nil {
			return 0
		}
		return len(m)
	})

	// string escapes
	e.RegisterFilter("cgi_escape", url.QueryEscape)
	e.RegisterFilter("smartify", smartifyFilter)
	e.RegisterFilter("uri_escape", func(s string) string {
		return regexp.MustCompile(`\?(.+?)=([^&]*)(?:\&(.+?)=([^&]*))*`).ReplaceAllStringFunc(s, func(m string) string {
			pair := strings.SplitN(m, "=", 2)
			return pair[0] + "=" + url.QueryEscape(pair[1])
		})
	})
}

// helpers

func requireNonEmptyArray(fn func([]interface{}) interface{}) func([]interface{}) interface{} {
	return func(array []interface{}) interface{} {
		if len(array) == 0 {
			return nil
		}
		return fn(array)
	}
}

// array filters

func arrayToSentenceStringFilter(array []string, conjunction func(string) string) string {
	conj := conjunction("and ")
	switch len(array) {
	case 1:
		return array[0]
	default:
		rt := reflect.ValueOf(array)
		ar := make([]string, rt.Len())
		for i, v := range array {
			ar[i] = v
			if i == rt.Len()-1 {
				ar[i] = conj + v
			}
		}
		return strings.Join(ar, ", ")
	}
}

func groupByExpFilter(array []map[string]interface{}, name string, expr expressions.Closure) ([]map[string]interface{}, error) {
	rt := reflect.ValueOf(array)
	if !(rt.Kind() != reflect.Array || rt.Kind() == reflect.Slice) {
		return nil, nil
	}
	groups := map[interface{}][]interface{}{}
	for i := 0; i < rt.Len(); i++ {
		item := rt.Index(i).Interface()
		key, err := expr.Bind(name, item).Evaluate()
		if err != nil {
			return nil, err
		}
		if group, found := groups[key]; found {
			groups[key] = append(group, item)
		} else {
			groups[key] = []interface{}{item}
		}
	}
	var result []map[string]interface{}
	for k, v := range groups {
		result = append(result, map[string]interface{}{"name": k, "items": v})
	}
	return result, nil
}

func groupByFilter(array []map[string]interface{}, property string) []map[string]interface{} {
	rt := reflect.ValueOf(array)
	if !(rt.Kind() != reflect.Array || rt.Kind() == reflect.Slice) {
		return nil
	}
	groups := map[interface{}][]interface{}{}
	for i := 0; i < rt.Len(); i++ {
		irt := rt.Index(i)
		if irt.Kind() == reflect.Map && irt.Type().Key().Kind() == reflect.String {
			krt := irt.MapIndex(reflect.ValueOf(property))
			if krt.IsValid() && krt.CanInterface() {
				key := krt.Interface()
				if group, found := groups[key]; found {
					groups[key] = append(group, irt.Interface())
				} else {
					groups[key] = []interface{}{irt.Interface()}
				}
			}
		}
	}
	var result []map[string]interface{}
	for k, v := range groups {
		result = append(result, map[string]interface{}{"name": k, "items": v})
	}
	return result
}

func sortFilter(array []interface{}, key interface{}, nilFirst func(bool) bool) []interface{} {
	nf := nilFirst(true)
	result := make([]interface{}, len(array))
	copy(result, array)
	if key == nil {
		evaluator.Sort(result)
	} else {
		// TODO error if key is not a string
		evaluator.SortByProperty(result, key.(string), nf)
	}
	return result
}

func whereExpFilter(array []interface{}, name string, expr expressions.Closure) ([]interface{}, error) {
	rt := reflect.ValueOf(array)
	if rt.Kind() != reflect.Array && rt.Kind() != reflect.Slice {
		return nil, nil
	}
	var result []interface{}
	for i := 0; i < rt.Len(); i++ {
		item := rt.Index(i).Interface()
		value, err := expr.Bind(name, item).Evaluate()
		if err != nil {
			return nil, err
		}
		if value != nil && value != false {
			result = append(result, item)
		}
	}
	return result, nil
}

func whereFilter(array []map[string]interface{}, key string, value interface{}) []interface{} {
	rt := reflect.ValueOf(array)
	if rt.Kind() != reflect.Array && rt.Kind() != reflect.Slice {
		return nil
	}
	var result []interface{}
	for i := 0; i < rt.Len(); i++ {
		item := rt.Index(i)
		if item.Kind() == reflect.Map && item.Type().Key().Kind() == reflect.String {
			attr := item.MapIndex(reflect.ValueOf(key))
			if attr.IsValid() && fmt.Sprint(attr) == value {
				result = append(result, item.Interface())
			}
		}
	}
	return result
}

// string filters

// -----------------------------------------------------------------------------
//
// https://github.com/osteele/gojekyll/blob/5ff68e140697bcb30660f083ad42dc328578b16f/filters/smartify.go
//
// -----------------------------------------------------------------------------

var smartifyTransforms = []struct {
	match *regexp.Regexp
	repl  string
}{
	{regexp.MustCompile("(^|[^[:alnum:]])``(.+?)''"), "$1“$2”"},
	{regexp.MustCompile(`(^|[^[:alnum:]])'`), "$1‘"},
	{regexp.MustCompile(`'`), "’"},
	{regexp.MustCompile(`(^|[^[:alnum:]])"`), "$1“"},
	{regexp.MustCompile(`"($|[^[:alnum:]])`), "”$1"},
	{regexp.MustCompile(`(^|\s)--($|\s)`), "$1–$2"},
	{regexp.MustCompile(`(^|\s)---($|\s)`), "$1—$2"},
}

// replace these wherever they appear
var smartifyReplaceSpans = map[string]string{
	"...":  "…",
	"(c)":  "©",
	"(r)":  "®",
	"(tm)": "™",
}

// replace these only if bounded by space or word boundaries
var smartifyReplaceWords = map[string]string{
	// "---": "–",
	// "--":  "—",
}

var smartifyReplacements map[string]string
var smartifyReplacementPattern *regexp.Regexp

func init() {
	smartifyReplacements = map[string]string{}
	var disjuncts []string
	regexQuoter := regexp.MustCompile(`[\(\)\.]`)
	escape := func(s string) string {
		return regexQuoter.ReplaceAllString(s, `\$0`)
	}
	for k, v := range smartifyReplaceSpans {
		disjuncts = append(disjuncts, escape(k))
		smartifyReplacements[k] = v
	}
	for k, v := range smartifyReplaceWords {
		disjuncts = append(disjuncts, fmt.Sprintf(`(\b|\s|^)%s(\b|\s|$)`, escape(k)))
		smartifyReplacements[k] = fmt.Sprintf("$1%s$2", v)
	}
	p := fmt.Sprintf(`(%s)`, strings.Join(disjuncts, `|`))
	smartifyReplacementPattern = regexp.MustCompile(p)
}

func smartifyFilter(s string) string {
	for _, rule := range smartifyTransforms {
		s = rule.match.ReplaceAllString(s, rule.repl)
	}
	s = smartifyReplacementPattern.ReplaceAllStringFunc(s, func(w string) string {
		return smartifyReplacements[w]
	})
	return s
}
