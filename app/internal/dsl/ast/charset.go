// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

package ast

import "github.com/kdeconinck/spot/internal/position"

// Charset represents a charset block.
type Charset struct {
	Range   position.Range
	Members []CharsetMember
}

// CharsetMember represents a named charset declaration.
type CharsetMember struct {
	Range position.Range
	Name  string
	Value Expression
}
