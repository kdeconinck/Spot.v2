// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

// Package position provides minimal source position tracking primitives.
package position

import "sort"

// SourceMap resolves byte offsets into line and column pairs.
type SourceMap struct {
	lineStarts []int
	size       int
}

// NewSourceMap builds a source map for text.
func NewSourceMap(text string) SourceMap {
	lineStarts := []int{0}

	for idx := range len(text) {
		if text[idx] == '\n' {
			lineStarts = append(lineStarts, idx+1)
		}
	}

	return SourceMap{
		lineStarts: lineStarts,
		size:       len(text),
	}
}

// Resolve converts pos into a 1-based line and byte column.
// It reports false when pos is outside the indexed text.
func (m SourceMap) Resolve(pos Position) (Location, bool) {
	if !pos.IsValid() || int(pos) > m.size {
		return Location{}, false
	}

	// Find the first line start after pos, then step back to get the line containing pos.
	lineIndex := sort.Search(len(m.lineStarts), func(i int) bool {
		return m.lineStarts[i] > int(pos)
	}) - 1

	lineStart := m.lineStarts[lineIndex]

	return Location{
		Line:   lineIndex + 1,
		Column: (int(pos) - lineStart) + 1,
	}, true
}
