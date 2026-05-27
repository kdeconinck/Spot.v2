// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

// Package token provides token definitions for the DSL scanner.
package token

import "github.com/kdeconinck/spot/internal/position"

// Token identifies a scanned token and the source range it came from.
type Token struct {
	Kind  Kind
	Range position.Range
}

// New returns a token with kind and source range.
func New(kind Kind, rangeInput position.Range) Token {
	return Token{
		Kind:  kind,
		Range: rangeInput,
	}
}
