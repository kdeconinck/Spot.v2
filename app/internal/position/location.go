// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

// Package position provides minimal source position tracking primitives.
package position

// Location identifies a 1-based line and byte column.
type Location struct {
	// Line is the 1-based line number.
	Line int

	// Column is the 1-based byte column number.
	Column int
}
