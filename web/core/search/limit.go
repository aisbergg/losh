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

package search

import (
	"regexp"
	"strings"

	"losh/web/core/search/parser"
)

const (
	maxQueryStringLength = 500
	maxWordsCount        = 30
	maxWildcardsCount    = 10
	maxNodesCount        = 15
)

type limiter struct {
	nodes     int
	words     int
	wildcards int

	operators map[string]struct{}
}

// newLimiter returns a new limiter.
func newLimiter() *limiter {
	return &limiter{
		operators: make(map[string]struct{}),
	}
}

// getOperators returns the list of operators used in the query.
func (l *limiter) getOperators() []string {
	oprs := make([]string, 0, len(l.operators))
	for op := range l.operators {
		oprs = append(oprs, op)
	}
	return oprs
}

// check checks whether the given query exceeds a limit.
func (l *limiter) check(query *parser.Query) error {
	l.checkQuery(query)
	if l.nodes > maxNodesCount {
		return &Error{"too many terms in query", ErrorLimitExceeded}
	}
	if l.words > maxWordsCount {
		return &Error{"too many words in query", ErrorLimitExceeded}
	}
	if l.wildcards > maxWildcardsCount {
		return &Error{"too many wildcards (*) in query", ErrorLimitExceeded}
	}
	return nil
}

func (l *limiter) checkQuery(query *parser.Query) {
	if query == nil {
		return
	}
	for _, orCnd := range query.Or {
		l.checkOrCondition(orCnd)
	}
}

func (l *limiter) checkOrCondition(orCondition *parser.OrCondition) {
	if orCondition == nil {
		return
	}
	for _, andCnd := range orCondition.And {
		l.checkAndCondition(andCnd)
	}
}

func (l *limiter) checkAndCondition(andCondition *parser.AndCondition) {
	if andCondition == nil {
		return
	}
	if andCondition.Not != nil {
		l.nodes++
		l.checkAndCondition(andCondition.Not)
		return
	}
	l.checkExpression(andCondition.Operand)
}

var (
	wordPattern     = regexp.MustCompile(`[\S]+`)
	wildcardPattern = regexp.MustCompile(`\*+`)
)

func (l *limiter) checkExpression(expression *parser.Expression) {
	if expression == nil {
		return
	}
	l.nodes++
	if expression.Sub != nil {
		l.checkQuery(expression.Sub)
	}
	if expression.Text != nil {
		l.checkText(expression.Text)
	}
	if expression.Operator != nil {
		l.checkOperator(expression.Operator)
	}
}

func (l *limiter) checkText(text *parser.Text) {
	s := ""
	if text.Words != nil {
		s = *text.Words
	} else if text.Exact != nil {
		s = *text.Exact
	}
	if s != "" {
		l.words += len(wordPattern.FindAllString(s, -1))
		l.wildcards += len(wildcardPattern.FindAllString(s, -1))
	}
}

// checkOperator
func (l *limiter) checkOperator(opr *parser.Operator) {
	oprName := strings.ToLower(opr.Name)
	if oprName == "is" || oprName == "has" {
		if opr.Value != nil {
			if opr.Value.Exact != nil {
				oprName += strings.ToLower(*opr.Value.Exact)
			} else if opr.Value.Words != nil {
				oprName += strings.ToLower(*opr.Value.Words)
			}
			l.operators[oprName] = struct{}{}
		}
	} else {
		l.operators[oprName] = struct{}{}
	}

	if opr == nil {
		return
	}
	if opr.Comparison != nil {
		l.checkText(opr.Comparison.Value)
	}
	if opr.Value != nil {
		l.checkText(opr.Value)
	}
}

// checkLimits checks if the given query is valid and if it exceeds the limit.
func checkLimits(query *parser.Query) error {
	return (&limiter{}).check(query)
}
