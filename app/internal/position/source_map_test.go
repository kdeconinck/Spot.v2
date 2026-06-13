// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

// Verify the public API of the position package.
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
			name:          "When resolving the start of the file, the returned value is correct.",
			positionInput: position.NewPosition(0),
			want: position.Location{
				Line:   1,
				Column: 1,
			},
			wantIsValidOutput: true,
		},
		{
			name:          "When resolving a newline byte, the returned value is correct.",
			positionInput: position.NewPosition(5),
			want: position.Location{
				Line:   1,
				Column: 6,
			},
			wantIsValidOutput: true,
		},
		{
			name:          "When resolving the start of the second line, the returned value is correct.",
			positionInput: position.NewPosition(6),
			want: position.Location{
				Line:   2,
				Column: 1,
			},
			wantIsValidOutput: true,
		},
		{
			name:          "When resolving the middle of the second line, the returned value is correct.",
			positionInput: position.NewPosition(8),
			want: position.Location{
				Line:   2,
				Column: 3,
			},
			wantIsValidOutput: true,
		},
		{
			name:          "When resolving the start of the third line, the returned value is correct.",
			positionInput: position.NewPosition(11),
			want: position.Location{
				Line:   3,
				Column: 1,
			},
			wantIsValidOutput: true,
		},
		{
			name:          "When resolving the trailing newline byte, the returned value is correct.",
			positionInput: position.NewPosition(16),
			want: position.Location{
				Line:   3,
				Column: 6,
			},
			wantIsValidOutput: true,
		},
		{
			name:          "When resolving the end of text after a trailing newline, the returned value is correct.",
			positionInput: position.NewPosition(len(text)),
			want: position.Location{
				Line:   4,
				Column: 1,
			},
			wantIsValidOutput: true,
		},
		{
			name:              "When resolving an invalid position sentinel, the returned value is correct.",
			positionInput:     position.InvalidPosition,
			wantIsValidOutput: false,
		},
		{
			name:              "When resolving a position past the end of the text, the returned value is correct.",
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
			claim.Equal(t, testCase.name, testCase.want, got, LocationLabel)
		})
	}
}

func Test_SourceMap_Resolve_EndOfTextWithoutTrailingNewline(t *testing.T) {
	t.Parallel()

	// Arrange.
	text := "alpha\nbeta\ngamma"
	sourceMap := position.NewSourceMap(text)

	// Act.
	got, valid := sourceMap.Resolve(position.NewPosition(len(text)))
	want := position.Location{
		Line:   3,
		Column: 6,
	}

	// Assert.
	claim.Equal(t, "When resolving at the end, using an input without a trailing newline, the returned value is correct.", true, valid, IsValidLabel)
	claim.Equal(t, "When resolving at the end using an input without a trailing newline, the returned value is correct.", want, got, LocationLabel)
}

func Test_SourceMap_Resolve_WithEmptyLines(t *testing.T) {
	t.Parallel()

	text := "alpha\n\nbeta"
	sourceMap := position.NewSourceMap(text)

	testCases := []struct {
		name              string
		positionInput     position.Position
		want              position.Location
		wantIsValidOutput bool
	}{
		{
			name:          "When resolving the newline byte ending the first line, the returned value is correct.",
			positionInput: position.NewPosition(5),
			want: position.Location{
				Line:   1,
				Column: 6,
			},
			wantIsValidOutput: true,
		},
		{
			name:          "When resolving the start of an empty line, the returned value is correct.",
			positionInput: position.NewPosition(6),
			want: position.Location{
				Line:   2,
				Column: 1,
			},
			wantIsValidOutput: true,
		},
		{
			name:          "When resolving the start of the line after an empty line, the returned value is correct.",
			positionInput: position.NewPosition(7),
			want: position.Location{
				Line:   3,
				Column: 1,
			},
			wantIsValidOutput: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Act.
			got, gotIsValid := sourceMap.Resolve(testCase.positionInput)

			// Assert.
			claim.Equal(t, testCase.name, testCase.wantIsValidOutput, gotIsValid, IsValidLabel)
			claim.Equal(t, testCase.name, testCase.want, got, LocationLabel)
		})
	}
}

func Test_SourceMap_Resolve_UsesByteColumns(t *testing.T) {
	t.Parallel()

	// "é" occupies two bytes in UTF-8.
	text := "aé\n"
	sourceMap := position.NewSourceMap(text)

	testCases := []struct {
		name              string
		positionInput     position.Position
		want              position.Location
		wantIsValidOutput bool
	}{
		{
			name:          "When resolving the first byte of a multibyte rune, the returned value is correct.",
			positionInput: position.NewPosition(1),
			want: position.Location{
				Line:   1,
				Column: 2,
			},
			wantIsValidOutput: true,
		},
		{
			name:          "When resolving the second byte of a multibyte rune, the returned value is correct.",
			positionInput: position.NewPosition(2),
			want: position.Location{
				Line:   1,
				Column: 3,
			},
			wantIsValidOutput: true,
		},
		{
			name:          "When resolving the newline after a multibyte rune, the returned value is correct.",
			positionInput: position.NewPosition(3),
			want: position.Location{
				Line:   1,
				Column: 4,
			},
			wantIsValidOutput: true,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Act.
			got, gotIsValid := sourceMap.Resolve(testCase.positionInput)

			// Assert.
			claim.Equal(t, testCase.name, testCase.wantIsValidOutput, gotIsValid, IsValidLabel)
			claim.Equal(t, testCase.name, testCase.want, got, LocationLabel)
		})
	}
}

func Test_SourceMap_Resolve_EmptyText(t *testing.T) {
	t.Parallel()

	// Arrange.
	sourceMap := position.NewSourceMap("")

	// Act.
	got, valid := sourceMap.Resolve(position.NewPosition(0))
	want := position.Location{
		Line: 1, Column: 1,
	}

	// Assert.
	claim.Equal(t, "When resolving position zero in empty text, the returned value is correct.", true, valid, IsValidLabel)
	claim.Equal(t, "When resolving position zero in empty text, the returned value is correct.", want, got, LocationLabel)
}

func benchmark_NewSourceMap(b *testing.B, lineCountInput int) {
	b.Helper()

	text := generateText(lineCountInput)

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

	text := generateText(lineCountInput)
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

func generateText(lineCountInput int) string {
	var sb strings.Builder

	for range lineCountInput {
		sb.WriteString("let value = 12345;\n")
	}

	return sb.String()
}
