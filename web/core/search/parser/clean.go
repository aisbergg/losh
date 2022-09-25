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

package parser

import "strings"

func cleanQuery(query *Query) *Query {
	if query == nil {
		return nil
	}
	orCnds := make([]*OrCondition, 0, len(query.Or))
	for _, orCnd := range query.Or {
		val := cleanOrCondition(orCnd)
		if val != nil {
			// flatten nested OR conditions
			if len(val.And) == 1 &&
				val.And[0].Operand != nil &&
				val.And[0].Operand.Sub != nil {
				orCnds = append(orCnds, val.And[0].Operand.Sub.Or...)
			} else {
				orCnds = append(orCnds, val)
			}
		}
	}
	query.Or = orCnds
	if len(query.Or) == 0 {
		return nil
	}
	return query
}

func cleanOrCondition(orCondition *OrCondition) *OrCondition {
	if orCondition == nil {
		return nil
	}
	tmpAndCnds := make([]*AndCondition, 0, len(orCondition.And))
	for _, andCnd := range orCondition.And {
		val := cleanAndCondition(andCnd)
		if val != nil {
			// flatten nested AND conditions
			if val.Operand != nil &&
				val.Operand.Sub != nil &&
				len(val.Operand.Sub.Or) == 1 {
				tmpAndCnds = append(tmpAndCnds, val.Operand.Sub.Or[0].And...)
			} else {
				tmpAndCnds = append(tmpAndCnds, val)
			}
		}
	}
	andCnds := make([]*AndCondition, 0, len(tmpAndCnds))
	words := []string{}
	for _, andCnd := range tmpAndCnds {
		if andCnd.Not == nil && andCnd.Operand.Text != nil && andCnd.Operand.Text.Words != nil {
			words = append(words, *andCnd.Operand.Text.Words)
			continue
		}
		andCnds = append(andCnds, andCnd)
	}
	if len(words) > 0 {
		s := strings.Join(words, " ")
		andCnds = append(andCnds, &AndCondition{
			Operand: &Expression{
				Text: &Text{
					Words: &s,
				},
			},
		})
	}
	orCondition.And = andCnds

	if len(orCondition.And) == 0 {
		return nil
	}
	return orCondition
}

func cleanAndCondition(andCondition *AndCondition) *AndCondition {
	if andCondition == nil {
		return nil
	}
	if andCondition.Not != nil {
		andCondition.Not = cleanAndCondition(andCondition.Not)
		if andCondition.Not != nil && andCondition.Not.Not != nil {
			return andCondition.Not.Not
		}
	}
	if andCondition.Operand != nil {
		andCondition.Operand = cleanExpression(andCondition.Operand)
	}
	if andCondition.Operand == nil && andCondition.Not == nil {
		return nil
	}
	return andCondition
}

func cleanExpression(expression *Expression) *Expression {
	if expression == nil {
		return nil
	}
	if expression.Sub != nil {
		expression.Sub = cleanQuery(expression.Sub)
		if expression.Sub == nil {
			return nil
		}
	}
	if expression.Operator == nil && expression.Text == nil && expression.Sub == nil {
		return nil
	}
	return expression
}
