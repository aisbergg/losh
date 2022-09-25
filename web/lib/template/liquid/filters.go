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

package liquid

import (
	"fmt"
	"math"
	gourl "net/url"
	"reflect"
	"regexp"
	"strings"
	"unicode/utf8"

	"losh/internal/lib/util/reflectutil"
	"losh/internal/lib/util/stringutil"

	"github.com/aisbergg/go-jsonpointer/pkg/jsonpointer"
	"github.com/golang-module/carbon/v2"
	"github.com/osteele/liquid"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

func addFilters(e *liquid.Engine) {
	addJekyllFilters(e)

	// map/slice filters
	e.RegisterFilter("dict2items", dict2itemsFilter)
	e.RegisterFilter("in", inFilter)
	e.RegisterFilter("get", getFilter)

	// format filters
	e.RegisterFilter("format_datetime", formatDatetimeFilter)
	e.RegisterFilter("format_number", formatNumberFilter)

	// string filters
	e.RegisterFilter("ellipses", ellipsesFilter)
	e.RegisterFilter("first_letter", firstLetterFilter)
	e.RegisterFilter("first_letters", firstLettersFilter)
	e.RegisterFilter("replace_regex", replaceRegexFilter)
	e.RegisterFilter("idhex", idhexFilter)
	e.RegisterFilter("short_version", shortenVersionFilter)

	// number filters
	e.RegisterFilter("min", minFilter)
	e.RegisterFilter("max", maxFilter)
	e.RegisterFilter("random_number", randomNumberFilter)
	e.RegisterFilter("random_item", randomItemFilter)

	// type conversion filters
	e.RegisterFilter("int", toIntFilter)
	e.RegisterFilter("int64", toInt64Filter)
	e.RegisterFilter("float", toFloatFilter)
	e.RegisterFilter("string", toStringFilter)

	// URL filters
	e.RegisterFilter("url_with_params", urlWithParamsFilter)

	// other
	e.RegisterFilter("ternary", ternaryFilter)
	e.RegisterFilter("is_nil", isNilFilter)
	e.RegisterFilter("is", isFilter)
	e.RegisterFilter("deref", dereferenceFilter)

}

func dict2itemsFilter(value map[string]interface{}) []interface{} {
	items := make([]interface{}, 0, len(value))
	for k, v := range value {
		item := map[string]interface{}{
			"key":   k,
			"value": v,
		}
		items = append(items, item)
	}
	return items
}

func inFilter(variable string, in interface{}) bool {
	inVal := reflect.ValueOf(in)
	switch inVal.Kind() {
	case reflect.Slice:
		for i := 0; i < inVal.Len(); i++ {
			if inVal.Index(i).String() == variable {
				return true
			}
		}
	case reflect.Map:
		for _, key := range inVal.MapKeys() {
			if key.String() == variable {
				return true
			}
		}
	}
	return false
}

func getFilter(obj interface{}, path interface{}) interface{} {
	ptr, err := jsonpointer.New(path)
	if err != nil {
		return nil
	}
	res, err := ptr.Get(obj)
	if err != nil {
		return nil
	}
	return res
}

func ellipsesFilter(input string, length int) string {
	if len(input) == 0 {
		return ""
	}
	return stringutil.Ellipses(input, length)
}

func firstLetterFilter(input string) string {
	if len(input) == 0 {
		return ""
	}
	firstLetter, _ := utf8.DecodeRuneInString(input)
	return string(firstLetter)
}

func firstLettersFilter(input string) string {
	words := strings.Fields(input)
	fl := make([]string, 0, len(words))
	for _, word := range words {
		fl = append(fl, firstLetterFilter(word))
	}
	return strings.Join(fl, "")
}

func formatDatetimeFilter(input, format, timezone string) (string, error) {
	input = strings.TrimSpace(strings.ToLower(input))
	var datetime carbon.Carbon
	if input == "now" {
		datetime = carbon.Now(timezone)
	} else {
		datetime = carbon.Parse(input, timezone)
	}
	if datetime.Error != nil {
		return "", datetime.Error
	}
	formatted := datetime.Format(format)
	if datetime.Error != nil {
		return "", datetime.Error
	}
	return formatted, nil
}

func formatNumberFilter(input interface{}) string {
	// TODO: make language specific
	p := message.NewPrinter(language.English)
	n := number.Decimal(input)
	return p.Sprintf("%v", n)
}

func idhexFilter(input string) string {
	return strings.TrimPrefix(strings.TrimSpace(input), "0x")
}

var hashVersionPattern = regexp.MustCompile(`^([a-f0-9]{40}|[a-f0-9]{64})$`)

func shortenVersionFilter(version string) string {
	if hashVersionPattern.MatchString(version) {
		// return first 7 characters of hash
		return version[:7]
	}

	// return as is
	return version
}

func isFilter(value interface{}, t ...string) bool {
	rv := reflectutil.Indirect(reflect.ValueOf(value))
	for _, t := range t {
		if rv.Type().String() == t {
			return true
		}
	}
	return false
}

func dereferenceFilter(value interface{}) interface{} {
	return reflectutil.Indirect(reflect.ValueOf(value)).Interface()
}

func isNilFilter(value interface{}) bool {
	return reflectutil.IsNil(value)
}

func minFilter(arg1 interface{}, argN ...interface{}) (min interface{}) {
	min = arg1
	minF := toFloatFilter(arg1)
	for _, arg := range argN {
		num := toFloatFilter(arg)
		if num < minF {
			min = arg
			minF = num
		}
	}
	return
}

func maxFilter(arg1 interface{}, argN ...interface{}) (max interface{}) {
	max = arg1
	maxF := toFloatFilter(arg1)
	for _, arg := range argN {
		num := toFloatFilter(arg)
		if num > maxF {
			max = arg
			maxF = num
		}
	}
	return
}

func randomNumberFilter(x, min, max float64, round uint) float64 {
	// based on https://github.com/codecalm/jekyll-random
	value := math.Mod((x*x*math.Pi*math.E*(max+1)*(math.Sin(x)/math.Cos(x*x))), (max+1-min)) + min

	if value > max {
		value = max
	}
	if value < min {
		value = min
	}

	value = roundFloat(value, round)
	return value
}

func randomItemFilter(x float64, items []interface{}) interface{} {
	index := int(randomNumberFilter(x, 0, float64(len(items)), 0))
	return items[index]
}

func roundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

func replaceRegexFilter(input, search, replace string) (string, error) {
	searchPattern, err := regexp.Compile(search)
	if err != nil {
		return "", err
	}
	return searchPattern.ReplaceAllString(input, replace), nil
}

func toIntFilter(input interface{}) int {
	numF := toFloatFilter(input)
	if math.IsNaN(numF) {
		return 0
	}
	return int(numF)
}

func toInt64Filter(input interface{}) int64 {
	numF := toFloatFilter(input)
	if math.IsNaN(numF) {
		return 0
	}
	return int64(numF)
}

func toFloatFilter(num interface{}) float64 {
	switch i := num.(type) {
	case float64:
		return i
	case float32:
		return float64(i)
	case int64:
		return float64(i)
	case int32:
		return float64(i)
	case int:
		return float64(i)
	case uint64:
		return float64(i)
	case uint32:
		return float64(i)
	case uint:
		return float64(i)
	default:
		return math.NaN()
	}
}

func toStringFilter(input interface{}) string {
	return fmt.Sprint(input)
}

func urlWithParamsFilter(url interface{}, args ...interface{}) (string, error) {
	var parsedURL *gourl.URL
	switch u := url.(type) {
	case *gourl.URL:
		parsedURL = u
	case string:
		var err error
		parsedURL, err = gourl.Parse(u)
		if err != nil {
			return "", err
		}
	default:
		return "", fmt.Errorf("url_with_params: invalid url type %T", url)
	}
	params := parsedURL.Query()
	for i := 0; i < len(args); i += 2 {
		key := fmt.Sprint(args[i])
		value := fmt.Sprint(args[i+1])
		params.Set(key, value)
	}
	parsedURL.RawQuery = params.Encode()
	return parsedURL.String(), nil
}

func ternaryFilter(cond bool, t, f interface{}) interface{} {
	if cond {
		return t
	}
	return f
}
