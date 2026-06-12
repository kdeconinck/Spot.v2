// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

package ast

import "github.com/kdeconinck/spot/internal/position"

// Lexer represents a lexer block.
type Lexer struct {
	Range position.Range
	Rules []LexerRule
}

// LexerRule represents a named lexer rule.
type LexerRule struct {
	Range position.Range
	Value Expression
}
