// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

package scanner_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/kdeconinck/spot/internal/dsl/scanner"
	"github.com/kdeconinck/spot/internal/dsl/token"
	"github.com/kdeconinck/spot/internal/position"
	"github.com/kdeconinck/spot/internal/qa/claim"
)

// Labels, used in the different assertion methods.
const (
	HasErrorLabel = "Has error?"
	TokenLabel    = "Token"
)

func Test_New(t *testing.T) {
	t.Parallel()

	for tcName, tc := range map[string]struct {
		textInput string
		want      token.Token
	}{
		"When creating a scanner for empty input, the first token is end of file.": {
			textInput: "",
			want: token.Token{
				Kind: token.EndOfFile,
				Range: position.Range{
					Start: 0,
					End:   0,
				},
			},
		},
	} {
		t.Run(tcName, func(t *testing.T) {
			t.Parallel()

			// Arrange.
			scn := scanner.New(strings.NewReader(tc.textInput))

			// Act.
			got, err := scn.Next()

			// Assert.
			claim.Equal(t, tcName, false, err != nil, HasErrorLabel)
			claim.Equal(t, tcName, tc.want, got, TokenLabel)
		})
	}
}

func Test_Scanner_Next(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		textInput          string
		amountOfCallsInput int
		wantTokensOutput   []token.Token
		wantErrorOutput    string
	}{
		{
			name:               "When scanning valid input, the returned tokens are correct.",
			textInput:          "scope {\n  include \"**/*.txt\"\n  exclude \"**/vendor/**\"\n  // comment\n}",
			amountOfCallsInput: 9,
			wantTokensOutput: []token.Token{
				newToken(token.Scope, 0, 5),
				newToken(token.LeftBrace, 6, 7),
				newToken(token.Include, 10, 17),
				newToken(token.String, 18, 28),
				newToken(token.Exclude, 31, 38),
				newToken(token.String, 39, 53),
				newToken(token.LineComment, 56, 66),
				newToken(token.RightBrace, 67, 68),
				newToken(token.EndOfFile, 68, 68),
			},
		},
		{
			name:               "When scanning an unknown identifier, the returned error is correct.",
			textInput:          "unknown",
			amountOfCallsInput: 1,
			wantErrorOutput:    "DSL Scanner: Unexpected identifier \"unknown\".",
		},
		{
			name:               "When scanning a partial scope keyword at end of file, the returned error is correct.",
			textInput:          "sc",
			amountOfCallsInput: 1,
			wantErrorOutput:    "DSL Scanner: Unexpected identifier \"sc\".",
		},
		{
			name:               "When scanning an identifier that diverges from an unknown identifier, the returned error is correct.",
			textInput:          "sknown",
			amountOfCallsInput: 1,
			wantErrorOutput:    "DSL Scanner: Unexpected identifier \"sknown\".",
		},
		{
			name:               "When scanning a scope prefix followed by more identifier bytes, the returned error is correct.",
			textInput:          "scopex",
			amountOfCallsInput: 1,
			wantErrorOutput:    "DSL Scanner: Unexpected identifier \"scopex\".",
		},
		{
			name:               "When scanning an unknown identifier followed by a delimiter, the returned error is correct.",
			textInput:          "unknown}",
			amountOfCallsInput: 1,
			wantErrorOutput:    "DSL Scanner: Unexpected identifier \"unknown\".",
		},
		{
			name:               "When scanning the scope keyword at end of file, the returned token is correct.",
			textInput:          "scope",
			amountOfCallsInput: 2,
			wantTokensOutput: []token.Token{
				newToken(token.Scope, 0, 5),
				newToken(token.EndOfFile, 5, 5),
			},
		},
		{
			name:               "When scanning an unexpected byte, the returned error is correct.",
			textInput:          "@",
			amountOfCallsInput: 1,
			wantErrorOutput:    "DSL Scanner: Unexpected byte '@'.",
		},
		{
			name:               "When scanning a trailing slash at end of file, the returned error is correct.",
			textInput:          "/",
			amountOfCallsInput: 1,
			wantErrorOutput:    "DSL Scanner: Unexpected byte '/'.",
		},
		{
			name:               "When scanning a slash that is not a comment, the returned error is correct.",
			textInput:          "/}",
			amountOfCallsInput: 1,
			wantErrorOutput:    "DSL Scanner: Unexpected byte '/'.",
		},
		{
			name:               "When scanning a comment at end of file, the returned token is correct.",
			textInput:          "// comment",
			amountOfCallsInput: 2,
			wantTokensOutput: []token.Token{
				newToken(token.LineComment, 0, 10),
				newToken(token.EndOfFile, 10, 10),
			},
		},
		{
			name:               "When scanning a string with supported escapes, the returned token is correct.",
			textInput:          "\"a\\\"b\\\\c\\nd\\re\\tf\"",
			amountOfCallsInput: 2,
			wantTokensOutput: []token.Token{
				newToken(token.String, 0, 18),
				newToken(token.EndOfFile, 18, 18),
			},
		},
		{
			name:               "When scanning an unterminated string, the returned error is correct.",
			textInput:          "\"abc",
			amountOfCallsInput: 1,
			wantErrorOutput:    "DSL Scanner: Unterminated string literal.",
		},
		{
			name:               "When scanning a string that ends after an escape prefix, the returned error is correct.",
			textInput:          "\"\\",
			amountOfCallsInput: 1,
			wantErrorOutput:    "DSL Scanner: Unterminated string literal.",
		},
		{
			name:               "When scanning a string with an unsupported escape, the returned error is correct.",
			textInput:          "\"\\x\"",
			amountOfCallsInput: 1,
			wantErrorOutput:    "DSL Scanner: Unexpected escape sequence \"\\\\x\".",
		},
		{
			name:               "When scanning empty input repeatedly, the returned end of file token is stable.",
			textInput:          "",
			amountOfCallsInput: 2,
			wantTokensOutput: []token.Token{
				newToken(token.EndOfFile, 0, 0),
				newToken(token.EndOfFile, 0, 0),
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Arrange.
			scn := scanner.New(strings.NewReader(testCase.textInput))

			// Act.
			gotTokensOutput, gotErrorOutput := collectTokens(scn, testCase.amountOfCallsInput)

			// Assert.
			claim.Equal(t, testCase.name, len(testCase.wantTokensOutput), len(gotTokensOutput), "Amount of tokens")

			for tokenIndex := range gotTokensOutput {
				claim.Equal(t, testCase.name, testCase.wantTokensOutput[tokenIndex], gotTokensOutput[tokenIndex], "Token")
			}

			if testCase.wantErrorOutput == "" {
				claim.Equal(t, testCase.name, true, gotErrorOutput == nil, "Has no error?")
				return
			}

			claim.Equal(t, testCase.name, false, gotErrorOutput == nil, "Has error?")
			claim.Equal(t, testCase.name, testCase.wantErrorOutput, gotErrorOutput.Error(), "Error")
		})
	}
}

func Test_Scanner_Next_ReaderError(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name                   string
		readerInput            *errorAfterNReader
		amountOfTokensInput    int
		wantTokensOutput       []token.Token
		wantHasErrorOutput     bool
		wantErrorMessageOutput string
	}{
		{
			name:                   "When the input reader fails before any token is complete, the returned error is propagated.",
			readerInput:            newErrorAfterNReader("sc", 2, errors.New("boom")),
			amountOfTokensInput:    1,
			wantHasErrorOutput:     true,
			wantErrorMessageOutput: "boom",
		},
		{
			name:                   "When the input reader fails while scanning a line comment, the returned error is propagated.",
			readerInput:            newErrorAfterNReader("// com", 6, errors.New("boom")),
			amountOfTokensInput:    1,
			wantHasErrorOutput:     true,
			wantErrorMessageOutput: "boom",
		},
		{
			name:                   "When the input reader fails after scanning scope before a delimiter is read, the returned error is propagated.",
			readerInput:            newErrorAfterNReader("scope", 5, errors.New("boom")),
			amountOfTokensInput:    1,
			wantHasErrorOutput:     true,
			wantErrorMessageOutput: "boom",
		},
		{
			name:                   "When the input reader fails while scanning an unknown identifier, the returned error is propagated.",
			readerInput:            newErrorAfterNReader("unknown", 3, errors.New("boom")),
			amountOfTokensInput:    1,
			wantHasErrorOutput:     true,
			wantErrorMessageOutput: "boom",
		},
		{
			name:                   "When the input reader fails while scanning string contents, the returned error is propagated.",
			readerInput:            newErrorAfterNReader("\"a", 2, errors.New("boom")),
			amountOfTokensInput:    1,
			wantHasErrorOutput:     true,
			wantErrorMessageOutput: "boom",
		},
		{
			name:                   "When the input reader fails while scanning an escaped byte in a string, the returned error is propagated.",
			readerInput:            newErrorAfterNReader("\"\\", 2, errors.New("boom")),
			amountOfTokensInput:    1,
			wantHasErrorOutput:     true,
			wantErrorMessageOutput: "boom",
		},
		{
			name:                "When the input reader fails after a complete token, the returned error is propagated on the next call.",
			readerInput:         newErrorAfterNReader("{", 1, errors.New("boom")),
			amountOfTokensInput: 2,
			wantTokensOutput: []token.Token{
				newToken(token.LeftBrace, 0, 1),
			},
			wantHasErrorOutput:     true,
			wantErrorMessageOutput: "boom",
		},
		{
			name:                   "When the input reader fails after a slash before a comment is confirmed, the returned error is propagated.",
			readerInput:            newErrorAfterNReader("/", 1, errors.New("boom")),
			amountOfTokensInput:    1,
			wantHasErrorOutput:     true,
			wantErrorMessageOutput: "boom",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Arrange.
			scn := scanner.New(testCase.readerInput)

			// Act.
			gotTokensOutput, gotErrorOutput := collectTokens(scn, testCase.amountOfTokensInput)

			// Assert.
			claim.Equal(t, testCase.name, len(testCase.wantTokensOutput), len(gotTokensOutput), "Amount of tokens")

			for tokenIndex := range gotTokensOutput {
				claim.Equal(t, testCase.name, testCase.wantTokensOutput[tokenIndex], gotTokensOutput[tokenIndex], "Token")
			}

			claim.Equal(t, testCase.name, testCase.wantHasErrorOutput, gotErrorOutput != nil, "Has error?")

			if gotErrorOutput != nil {
				claim.Equal(t, testCase.name, testCase.wantErrorMessageOutput, gotErrorOutput.Error(), "Error")
			}
		})
	}
}

func benchmark_New(b *testing.B, blockCountInput int) {
	b.Helper()

	text := benchmarkText(blockCountInput)

	for b.Loop() {
		_ = scanner.New(strings.NewReader(text))
	}
}

func Benchmark_New_1Block(b *testing.B)     { benchmark_New(b, 1) }
func Benchmark_New_10Blocks(b *testing.B)   { benchmark_New(b, 10) }
func Benchmark_New_100Blocks(b *testing.B)  { benchmark_New(b, 100) }
func Benchmark_New_1000Blocks(b *testing.B) { benchmark_New(b, 1000) }

func benchmark_Scanner_Next(b *testing.B, blockCountInput int) {
	b.Helper()

	text := benchmarkText(blockCountInput)

	for b.Loop() {
		scn := scanner.New(strings.NewReader(text))

		for {
			nextToken, err := scn.Next()
			if err != nil {
				b.Fatalf("Next() returned an unexpected error: %v.", err)
			}

			if nextToken.Kind == token.EndOfFile {
				break
			}
		}
	}
}

func Benchmark_Scanner_Next_1Block(b *testing.B)     { benchmark_Scanner_Next(b, 1) }
func Benchmark_Scanner_Next_10Blocks(b *testing.B)   { benchmark_Scanner_Next(b, 10) }
func Benchmark_Scanner_Next_100Blocks(b *testing.B)  { benchmark_Scanner_Next(b, 100) }
func Benchmark_Scanner_Next_1000Blocks(b *testing.B) { benchmark_Scanner_Next(b, 1000) }

func collectTokens(scannerInput *scanner.Scanner, amountOfTokensInput int) ([]token.Token, error) {
	tokensOutput := make([]token.Token, 0, amountOfTokensInput)

	for range amountOfTokensInput {
		nextToken, err := scannerInput.Next()

		if err != nil {
			return tokensOutput, err
		}

		tokensOutput = append(tokensOutput, nextToken)
	}

	return tokensOutput, nil
}

func benchmarkText(blockCountInput int) string {
	var sb strings.Builder

	for range blockCountInput {
		sb.WriteString("scope {\n  include \"**/*.txt\"\n  exclude \"**/vendor/**\"\n  // comment\n}\n")
	}

	return sb.String()
}

func newToken(kindInput token.Kind, startInput, endInput int) token.Token {
	return token.New(
		kindInput,
		position.NewRange(position.NewPosition(startInput), position.NewPosition(endInput)))
}

type errorAfterNReader struct {
	text      string
	readCount int
	failAfter int
	failWith  error
}

func newErrorAfterNReader(textInput string, failAfterInput int, failWithInput error) *errorAfterNReader {
	return &errorAfterNReader{
		text:      textInput,
		failAfter: failAfterInput,
		failWith:  failWithInput,
	}
}

func (reader *errorAfterNReader) Read(buffer []byte) (int, error) {
	if reader.readCount >= reader.failAfter {
		return 0, reader.failWith
	}

	if reader.readCount >= len(reader.text) {
		return 0, io.EOF
	}

	buffer[0] = reader.text[reader.readCount]
	reader.readCount++

	return 1, nil
}
