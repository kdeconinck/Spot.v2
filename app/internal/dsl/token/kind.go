// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

// Package token provides token definitions for the DSL scanner.
package token

import "strconv"

// Kind identifies the type of a scanned token.
type Kind uint8

// Supported token kinds.
const (
	Invalid Kind = iota
	EndOfFile

	// Keywords.
	Scope

	// Symbols.
	LeftBrace
	RightBrace

	// Trivia.
	LineComment
)

// String returns a human-readable name for kind.
func (kind Kind) String() string {
	switch kind {
	case Invalid:
		return "invalid"

	case EndOfFile:
		return "end-of-file"

	case Scope:
		return "scope"

	case LeftBrace:
		return "left-brace"

	case RightBrace:
		return "right-brace"

	case LineComment:
		return "line-comment"

	default:
		return strconv.Itoa(int(kind))
	}
}
