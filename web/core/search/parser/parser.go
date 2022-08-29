package parser

import (
	"fmt"
	"strings"

	"github.com/aisbergg/go-errors/pkg/errors"
	"github.com/alecthomas/participle/v2"
	"github.com/alecthomas/participle/v2/lexer"
)

type Boolean bool

func (b *Boolean) Capture(values []string) error {
	*b = values[0] == "true"
	return nil
}

type CompOperator int

const (
	CompOpEq CompOperator = iota
	CompOpNe
	CompOpLt
	CompOpLe
	CompOpGt
	CompOpGe
)

var compOpMap = map[string]CompOperator{
	"==": CompOpEq,
	"!=": CompOpNe,
	"<":  CompOpLt,
	"<=": CompOpLe,
	">":  CompOpGt,
	">=": CompOpGe,
}

// -----------------------------------------------------------------------------

type Query struct {
	Or []*OrCondition `@@ (Whitespace ("OR" | "|") Whitespace @@)*`
}

type OrCondition struct {
	And []*AndCondition `@@ ((Whitespace ("AND" | "&") Whitespace @@) | Whitespace @@)*`
}

type AndCondition struct {
	Not     *AndCondition `("NOT" Whitespace @@) | ("-" @@)`
	Operand *Expression   `| @@`
}

type Expression struct {
	Operator *Operator `  @@`
	Text     *Text     `| @@`
	Sub      *Query    `| "(" @@ ")"`
	Discard  *string   `| Whitespace`
}

type Text struct {
	Exact *string `  @QuotedString`
	Words *string `| (@BacktickQuotedString | (@Identifier | @Number | @String | @Specials) (@Identifier | @Keyword | @Number | @String | @Specials)*)`
}

func (co *CompOperator) Capture(s []string) error {
	str := strings.Join(s, "")
	*co = compOpMap[str]
	return nil
}

type Comparison struct {
	Operator CompOperator `@("=" "=" | "!" "=" | "<" | ">" | "<" "=" | ">" "=")`
	Value    *Text        `@@`
}

type Range struct {
	OpenStart bool    `( @"*"`
	Start     *string `| @QuotedString | (@String | @Identifier | @Keyword | @Number | @Specials)+) DoubleDot`
	OpenEnd   bool    `( @"*"`
	End       *string `| @QuotedString | (@String | @Identifier | @Keyword | @Number | @Specials)+)`
}

type Operator struct {
	Name       string      `@Identifier ":"`
	Comparison *Comparison `( @@`
	Range      *Range      `| @@`
	Value      *Text       `| @@ )?`
}

// TODO: remove later
// type Value struct {
// 	Number *float64 `  @Number`
// 	Bool   *Boolean `| @("true" | "false")`
// 	Text   *Text    `| @@`
// }

var queryLexer = lexer.MustSimple([]lexer.SimpleRule{
	{"Number", `[-+]?\d+(_\d+)*(\.\d+(_\d+)*)?`},
	{"Keyword", `AND|&|OR|\||NOT|-`},
	{"Identifier", `[a-zA-Z]+`},
	{"QuotedString", `("(?:[^"\\]|\\.)+")|('(?:[^'\\]|\\.)+')`},
	{"BacktickQuotedString", "`(?:[^`\\\\]|\\\\.)+`"},
	{"Group", `\(|\)`},
	{"DoubleDot", `\.\.`},
	{"Specials", `[-[!@#$%^&*+_={}\|:;"'<,>.?/\]]`},
	{"String", `\S+`},
	{"Whitespace", `\s+`},
})

var parser = participle.MustBuild[Query](
	participle.Lexer(queryLexer),
	participle.UseLookahead(2),
	participle.Unquote("QuotedString", "BacktickQuotedString"),
)

// Parse parses the given query string and returns the AST.
func Parse(query string) (*Query, error) {
	parsed, err := parser.ParseString("", query, participle.AllowTrailing(true))
	parsed = cleanQuery(parsed)
	if err != nil {
		fmt.Println("error parsing query:", err)
		return parsed, errors.CEWrap(err, "failed to parse query").Add("query", query)
	}
	return parsed, nil
}
