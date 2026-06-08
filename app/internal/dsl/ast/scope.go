// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

package ast

import "github.com/kdeconinck/spot/internal/position"

// Scope is the top-level DSL scope node.
type Scope struct {
	Range    position.Range
	Includes []Include
	Excludes []Exclude
	Charset  *Charset
	Lexer    *Lexer
}

// Include represents an include statement.
type Include struct {
	Range   position.Range
	Pattern StringLiteral
}

// Exclude represents an exclude statement.
type Exclude struct {
	Range   position.Range
	Pattern StringLiteral
}
