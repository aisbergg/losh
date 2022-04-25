package liquid

import (
	"regexp"
	"strings"

	"github.com/osteele/liquid"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/number"
)

func addFilters(e *liquid.Engine) {
	addJekyllFilters(e)
	e.RegisterFilter("first_letter", firstLetterFilter)
	e.RegisterFilter("first_letters", firstLettersFilter)
	e.RegisterFilter("replace_regex", replaceRegexFilter)
}

func firstLetterFilter(input string) string {
	if len(input) == 0 {
		return ""
	}
	return input[:1]
}

func firstLettersFilter(input string) string {
	words := strings.Fields(input)
	fl := make([]string, 0, len(words))
	for _, word := range words {
		fl = append(fl, word[:1])
	}
	return strings.Join(fl, "")
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
