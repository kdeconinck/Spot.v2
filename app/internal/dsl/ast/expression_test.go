// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

package ast_test

import (
	"testing"

	"github.com/kdeconinck/spot/internal/dsl/ast"
	"github.com/kdeconinck/spot/internal/qa/claim"
)

func Test_Expression(t *testing.T) {
	t.Parallel()

	boundedMaximum := 8

	testCases := []struct {
		name            string
		expressionInput ast.Expression
	}{
		{
			name: "When using a reference expression, it is assignable to the expression interface.",
			expressionInput: ast.Reference{
				Range: newRange(0, 5),
			},
		},
		{
			name: "When using a string literal expression, it is assignable to the expression interface.",
			expressionInput: ast.StringLiteral{
				Range: newRange(5, 13),
			},
		},
		{
			name: "When using a character literal expression, it is assignable to the expression interface.",
			expressionInput: ast.CharacterLiteral{
				Range: newRange(13, 17),
			},
		},
		{
			name: "When using a character range expression, it is assignable to the expression interface.",
			expressionInput: ast.CharacterRange{
				Range: newRange(17, 25),
				Start: ast.CharacterLiteral{
					Range: newRange(17, 20)},
				End: ast.CharacterLiteral{
					Range: newRange(22, 25)},
			},
		},
		{
			name: "When using an alternation expression, it is assignable to the expression interface.",
			expressionInput: ast.Alternation{
				Range: newRange(25, 36),
				Expressions: []ast.Expression{
					ast.Reference{
						Range: newRange(25, 31)},
					ast.Reference{
						Range: newRange(34, 36)},
				},
			},
		},
		{
			name: "When using a concatenation expression, it is assignable to the expression interface.",
			expressionInput: ast.Concatenation{
				Range: newRange(36, 52),
				Expressions: []ast.Expression{
					ast.Reference{
						Range: newRange(36, 47)},
					ast.Reference{
						Range: newRange(48, 52)},
				},
			},
		},
		{
			name: "When using an unbounded repetition expression, it is assignable to the expression interface.",
			expressionInput: ast.Repetition{
				Range: newRange(52, 63),
				Expression: ast.Reference{
					Range: newRange(52, 57)},
				Minimum: 1,
				Maximum: nil,
			},
		},
		{
			name: "When using a bounded repetition expression, it is assignable to the expression interface.",
			expressionInput: ast.Repetition{
				Range: newRange(63, 79),
				Expression: ast.Reference{
					Range: newRange(63, 74)},
				Minimum: 1,
				Maximum: &boundedMaximum,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			claim.Equal(t, testCase.name, false, testCase.expressionInput == nil, "Is nil?")
		})
	}
}
