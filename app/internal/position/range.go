// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

// Package position provides minimal source position tracking primitives.
package position

// Range identifies a half-open source interval [Start, End).
type Range struct {
	Start Position
	End   Position
}

// NewRange returns the range from start to end.
func NewRange(start, end Position) Range {
	return Range{
		Start: start,
		End:   end,
	}
}

// IsValid reports whether r is well-formed.
func (r Range) IsValid() bool { return r.Start.IsValid() && r.End.IsValid() && r.Start <= r.End }
