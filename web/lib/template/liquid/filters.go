package liquid

import (
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/golang-module/carbon/v2"
	"github.com/osteele/liquid"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

func addFilters(e *liquid.Engine) {
	addJekyllFilters(e)
	e.RegisterFilter("dict2items", dict2itemsFilter)
	e.RegisterFilter("first_letter", firstLetterFilter)
	e.RegisterFilter("first_letters", firstLettersFilter)
	e.RegisterFilter("format_datetime", formatDatetimeFilter)
	e.RegisterFilter("format_number", formatNumberFilter)
	e.RegisterFilter("replace_regex", replaceRegexFilter)
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

func replaceRegexFilter(input, search, replace string) (string, error) {
	searchPattern, err := regexp.Compile(search)
	if err != nil {
		return "", err
	}
	return searchPattern.ReplaceAllString(input, replace), nil
}
