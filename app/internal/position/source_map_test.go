// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

package position_test

import (
	"strings"
	"testing"

	"github.com/kdeconinck/spot/internal/position"
	"github.com/kdeconinck/spot/internal/qa/claim"
)

func Test_SourceMap_Resolve(t *testing.T) {
	t.Parallel()

	text := "alpha\nbeta\ngamma\n"
	sourceMap := position.NewSourceMap(text)

	testCases := []struct {
		name              string
		positionInput     position.Position
		want              position.Location
		wantIsValidOutput bool
	}{
		{
			name:          "When resolving the start of the file, the returned location is [1:1].",
			positionInput: position.NewPosition(0),
			want: position.Location{
				Line:   1,
				Column: 1,
			},
			wantIsValidOutput: true,
		},
		{
			name:          "When resolving a newline byte, the returned location is [1:6].",
			positionInput: position.NewPosition(5),
			want: position.Location{
				Line:   1,
				Column: 6,
			},
			wantIsValidOutput: true,
		},
		{
			name:          "When resolving the start of the second line, the returned location is [2:1].",
			positionInput: position.NewPosition(6),
			want: position.Location{
				Line:   2,
				Column: 1,
			},
			wantIsValidOutput: true,
		},
		{
			name:          "When resolving the middle of the second line, the returned location is [2:3].",
			positionInput: position.NewPosition(8),
			want: position.Location{
				Line:   2,
				Column: 3,
			},
			wantIsValidOutput: true,
		},
		{
			name:          "When resolving the start of the third line, the returned location is [3:1].",
			positionInput: position.NewPosition(11),
			want: position.Location{
				Line:   3,
				Column: 1,
			},
			wantIsValidOutput: true,
		},
		{
			name:          "When resolving the trailing newline byte, the returned location is [3:6].",
			positionInput: position.NewPosition(16),
			want: position.Location{
				Line:   3,
				Column: 6,
			},
			wantIsValidOutput: true,
		},
		{
			name:          "When resolving the end of text after a trailing newline, the returned location is [4:1].",
			positionInput: position.NewPosition(len(text)),
			want: position.Location{
				Line:   4,
				Column: 1,
			},
			wantIsValidOutput: true,
		},
		{
			name:              "When resolving an invalid position sentinel, the returned location is not valid.",
			positionInput:     position.InvalidPosition,
			wantIsValidOutput: false,
		},
		{
			name:              "When resolving a position past the end of the text, the returned location is not valid.",
			positionInput:     position.NewPosition(len(text) + 1),
			wantIsValidOutput: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Act.
			got, gotIsValid := sourceMap.Resolve(testCase.positionInput)

			// Assert.
			claim.Equal(t, testCase.name, testCase.wantIsValidOutput, gotIsValid, IsValidLabel)

			if !gotIsValid {
				return
			}

			claim.Equal(t, testCase.name, testCase.want, got, LocationLabel)
		})
	}
}

func Test_SourceMap_Resolve_EmptyText(t *testing.T) {
	t.Parallel()

	sourceMap := position.NewSourceMap("")

	// Act.
	got, valid := sourceMap.Resolve(position.NewPosition(0))

	// Assert.
	claim.Equal(t, "When resolving position zero in empty text, the returned location location is valid.", true, valid, IsValidLabel)

	claim.Equal(t, "When resolving position zero in empty text, the returned location is [1:1].", position.Location{
		Line: 1, Column: 1,
	}, got, LocationLabel)
}

func benchmark_NewSourceMap(b *testing.B, lineCountInput int) {
	b.Helper()

	text := benchmarkText(lineCountInput)

	for b.Loop() {
		_ = position.NewSourceMap(text)
	}
}

func Benchmark_NewSourceMap_1Line(b *testing.B)     { benchmark_NewSourceMap(b, 1) }
func Benchmark_NewSourceMap_10Lines(b *testing.B)   { benchmark_NewSourceMap(b, 10) }
func Benchmark_NewSourceMap_100Lines(b *testing.B)  { benchmark_NewSourceMap(b, 100) }
func Benchmark_NewSourceMap_1000Lines(b *testing.B) { benchmark_NewSourceMap(b, 1000) }

func benchmark_SourceMap_Resolve(b *testing.B, lineCountInput int) {
	b.Helper()

	text := benchmarkText(lineCountInput)
	sourceMap := position.NewSourceMap(text)

	positions := []position.Position{
		position.NewPosition(0),
		position.NewPosition(len(text) / 4),
		position.NewPosition(len(text) / 2),
		position.NewPosition((len(text) * 3) / 4),
		position.NewPosition(len(text)),
	}

	for i := 0; b.Loop(); i++ {
		sourceMap.Resolve(positions[i%len(positions)])
	}
}

func Benchmark_SourceMap_Resolve_1Line(b *testing.B)     { benchmark_SourceMap_Resolve(b, 1) }
func Benchmark_SourceMap_Resolve_10Lines(b *testing.B)   { benchmark_SourceMap_Resolve(b, 10) }
func Benchmark_SourceMap_Resolve_100Lines(b *testing.B)  { benchmark_SourceMap_Resolve(b, 100) }
func Benchmark_SourceMap_Resolve_1000Lines(b *testing.B) { benchmark_SourceMap_Resolve(b, 1000) }

func benchmarkText(lineCountInput int) string {
	var sb strings.Builder

	for range lineCountInput {
		sb.WriteString("let value = 12345;\n")
	}

	return sb.String()
}
