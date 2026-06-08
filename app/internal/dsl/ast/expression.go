// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

package ast

import "github.com/kdeconinck/spot/internal/position"

// Reference represents a named reference in an expression.
type Reference struct {
	Range position.Range
	Name  string
}

func (Reference) isExpression() {
	// NOTE: Intentionally left blank.
}

// StringLiteral represents a string literal expression.
type StringLiteral struct {
	Range position.Range
	Value string
}

func (StringLiteral) isExpression() {
	// NOTE: Intentionally left blank.
}

// CharacterLiteral represents a character literal expression.
type CharacterLiteral struct {
	Range position.Range
	Value rune
}

func (CharacterLiteral) isExpression() {
	// NOTE: Intentionally left blank.
}

// CharacterRange represents an inclusive character range expression.
type CharacterRange struct {
	Range position.Range
	Start CharacterLiteral
	End   CharacterLiteral
}

func (CharacterRange) isExpression() {
	// NOTE: Intentionally left blank.
}

// Alternation represents an alternation expression.
type Alternation struct {
	Range       position.Range
	Expressions []Expression
}

func (Alternation) isExpression() {
	// NOTE: Intentionally left blank.
}

// Concatenation represents a concatenation expression.
type Concatenation struct {
	Range       position.Range
	Expressions []Expression
}

func (Concatenation) isExpression() {
	// NOTE: Intentionally left blank.
}

// Repetition represents a repetition expression.
// A nil maximum means the repetition is unbounded.
type Repetition struct {
	Range      position.Range
	Expression Expression
	Minimum    int
	Maximum    *int
}

func (Repetition) isExpression() {
	// NOTE: Intentionally left blank.
}
