// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

// Package position provides minimal source position tracking primitives.
package position

import "sort"

// SourceMap resolves byte offsets into 1-based line and byte column pairs.
type SourceMap struct {
	lineStartOffsets []int
	size             int
}

// NewSourceMap returns a source map for text.
func NewSourceMap(text string) SourceMap {
	lineStartOffsets := []int{0}

	for idx := range len(text) {
		if text[idx] == '\n' {
			lineStartOffsets = append(lineStartOffsets, idx+1)
		}
	}

	return SourceMap{
		lineStartOffsets: lineStartOffsets,
		size:             len(text),
	}
}

// Resolve converts pos into a 1-based line and byte column.
// It reports false if pos is invalid or outside the indexed text.
func (m SourceMap) Resolve(pos Position) (Location, bool) {
	if !pos.IsValid() || int(pos) > m.size {
		return Location{}, false
	}

	// Find the first line start after pos, then step back to get the line containing pos.
	lineIndex := sort.Search(len(m.lineStartOffsets), func(i int) bool {
		return m.lineStartOffsets[i] > int(pos)
	}) - 1

	lineStartOffset := m.lineStartOffsets[lineIndex]

	return Location{
		Line:   lineIndex + 1,
		Column: (int(pos) - lineStartOffset) + 1,
	}, true
}
