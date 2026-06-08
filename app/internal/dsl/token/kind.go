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
	Charset
	Include
	Exclude
	Identifier

	// Symbols.
	LeftBrace
	RightBrace
	LeftParen
	RightParen
	Equal
	Pipe
	DotDot
	Star
	Plus
	Question
	Comma

	// Literals.
	String
	Character
	Number

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

	case Charset:
		return "charset"

	case Include:
		return "include"

	case Exclude:
		return "exclude"

	case Identifier:
		return "identifier"

	case LeftBrace:
		return "left-brace"

	case RightBrace:
		return "right-brace"

	case LeftParen:
		return "left-paren"

	case RightParen:
		return "right-paren"

	case Equal:
		return "equal"

	case Pipe:
		return "pipe"

	case DotDot:
		return "dot-dot"

	case Star:
		return "star"

	case Plus:
		return "plus"

	case Question:
		return "question"

	case Comma:
		return "comma"

	case String:
		return "string"

	case Character:
		return "character"

	case Number:
		return "number"

	case LineComment:
		return "line-comment"

	default:
		return strconv.Itoa(int(kind))
	}
}
