// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

// Package position provides minimal source position tracking primitives.
package position

// Position identifies a byte offset in source text.
type Position int

// InvalidPosition is a sentinel for an invalid source position.
const InvalidPosition Position = -1

// NewPosition returns a position for offset.
func NewPosition(offset int) Position {
	if offset < 0 {
		return InvalidPosition
	}

	return Position(offset)
}

// IsValid reports whether pos refers to a non-negative source offset.
func (pos Position) IsValid() bool { return pos >= 0 }
