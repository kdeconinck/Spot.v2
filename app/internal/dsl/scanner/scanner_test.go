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
	ErrorLabel    = "Error"
	TokenLabel    = "Token"
	TokensLabel   = "Tokens"
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
			name:               "When scanning valid input with a charset block, the returned tokens are correct.",
			textInput:          "scope {\n  include \"**/*.txt\"\n  exclude \"**/vendor/**\"\n  charset {\n    letters = ('a'..'z' | 'A'..'Z')+\n    decimal_numbers = ('0'..'9'){1,8}\n    space = ' ' | '\\t'?\n    optional = ('x' | 'y')*\n  }\n  // comment\n}",
			amountOfCallsInput: 50,
			wantTokensOutput: []token.Token{
				newToken(token.Scope, 0, 5),
				newToken(token.LeftBrace, 6, 7),
				newToken(token.Include, 10, 17),
				newToken(token.String, 18, 28),
				newToken(token.Exclude, 31, 38),
				newToken(token.String, 39, 53),
				newToken(token.Charset, 56, 63),
				newToken(token.LeftBrace, 64, 65),
				newToken(token.Identifier, 70, 77),
				newToken(token.Equal, 78, 79),
				newToken(token.LeftParen, 80, 81),
				newToken(token.Character, 81, 84),
				newToken(token.DotDot, 84, 86),
				newToken(token.Character, 86, 89),
				newToken(token.Pipe, 90, 91),
				newToken(token.Character, 92, 95),
				newToken(token.DotDot, 95, 97),
				newToken(token.Character, 97, 100),
				newToken(token.RightParen, 100, 101),
				newToken(token.Plus, 101, 102),
				newToken(token.Identifier, 107, 122),
				newToken(token.Equal, 123, 124),
				newToken(token.LeftParen, 125, 126),
				newToken(token.Character, 126, 129),
				newToken(token.DotDot, 129, 131),
				newToken(token.Character, 131, 134),
				newToken(token.RightParen, 134, 135),
				newToken(token.LeftBrace, 135, 136),
				newToken(token.Number, 136, 137),
				newToken(token.Comma, 137, 138),
				newToken(token.Number, 138, 139),
				newToken(token.RightBrace, 139, 140),
				newToken(token.Identifier, 145, 150),
				newToken(token.Equal, 151, 152),
				newToken(token.Character, 153, 156),
				newToken(token.Pipe, 157, 158),
				newToken(token.Character, 159, 163),
				newToken(token.Question, 163, 164),
				newToken(token.Identifier, 169, 177),
				newToken(token.Equal, 178, 179),
				newToken(token.LeftParen, 180, 181),
				newToken(token.Character, 181, 184),
				newToken(token.Pipe, 185, 186),
				newToken(token.Character, 187, 190),
				newToken(token.RightParen, 190, 191),
				newToken(token.Star, 191, 192),
				newToken(token.RightBrace, 195, 196),
				newToken(token.LineComment, 199, 209),
				newToken(token.RightBrace, 210, 211),
				newToken(token.EndOfFile, 211, 211),
			},
		},
		{
			name:               "When scanning valid input with a lexer block, the returned tokens are correct.",
			textInput:          "scope {\n  charset {\n    ident_start = 'a'..'z' | 'A'..'Z' | '_'\n    ident_part = ident_start | '0'..'9'\n    space = ' ' | '\\t'\n  }\n  lexer {\n    identifier = ident_start ident_part*\n    whitespace = space+\n    public_keyword = \"public\"\n  }\n}",
			amountOfCallsInput: 45,
			wantTokensOutput: []token.Token{
				newToken(token.Scope, 0, 5),
				newToken(token.LeftBrace, 6, 7),
				newToken(token.Charset, 10, 17),
				newToken(token.LeftBrace, 18, 19),
				newToken(token.Identifier, 24, 35),
				newToken(token.Equal, 36, 37),
				newToken(token.Character, 38, 41),
				newToken(token.DotDot, 41, 43),
				newToken(token.Character, 43, 46),
				newToken(token.Pipe, 47, 48),
				newToken(token.Character, 49, 52),
				newToken(token.DotDot, 52, 54),
				newToken(token.Character, 54, 57),
				newToken(token.Pipe, 58, 59),
				newToken(token.Character, 60, 63),
				newToken(token.Identifier, 68, 78),
				newToken(token.Equal, 79, 80),
				newToken(token.Identifier, 81, 92),
				newToken(token.Pipe, 93, 94),
				newToken(token.Character, 95, 98),
				newToken(token.DotDot, 98, 100),
				newToken(token.Character, 100, 103),
				newToken(token.Identifier, 108, 113),
				newToken(token.Equal, 114, 115),
				newToken(token.Character, 116, 119),
				newToken(token.Pipe, 120, 121),
				newToken(token.Character, 122, 126),
				newToken(token.RightBrace, 129, 130),
				newToken(token.Lexer, 133, 138),
				newToken(token.LeftBrace, 139, 140),
				newToken(token.Identifier, 145, 155),
				newToken(token.Equal, 156, 157),
				newToken(token.Identifier, 158, 169),
				newToken(token.Identifier, 170, 180),
				newToken(token.Star, 180, 181),
				newToken(token.Identifier, 186, 196),
				newToken(token.Equal, 197, 198),
				newToken(token.Identifier, 199, 204),
				newToken(token.Plus, 204, 205),
				newToken(token.Identifier, 210, 224),
				newToken(token.Equal, 225, 226),
				newToken(token.String, 227, 235),
				newToken(token.RightBrace, 238, 239),
				newToken(token.RightBrace, 240, 241),
				newToken(token.EndOfFile, 241, 241),
			},
		},
		{
			name:               "When scanning an unknown identifier, the returned token is correct.",
			textInput:          "unknown",
			amountOfCallsInput: 2,
			wantTokensOutput: []token.Token{
				newToken(token.Identifier, 0, 7),
				newToken(token.EndOfFile, 7, 7),
			},
		},
		{
			name:               "When scanning a partial scope keyword at end of file, the returned token is correct.",
			textInput:          "sc",
			amountOfCallsInput: 2,
			wantTokensOutput: []token.Token{
				newToken(token.Identifier, 0, 2),
				newToken(token.EndOfFile, 2, 2),
			},
		},
		{
			name:               "When scanning an identifier that diverges from a keyword, the returned token is correct.",
			textInput:          "sknown",
			amountOfCallsInput: 2,
			wantTokensOutput: []token.Token{
				newToken(token.Identifier, 0, 6),
				newToken(token.EndOfFile, 6, 6),
			},
		},
		{
			name:               "When scanning a scope prefix followed by more identifier bytes, the returned token is correct.",
			textInput:          "scopex",
			amountOfCallsInput: 2,
			wantTokensOutput: []token.Token{
				newToken(token.Identifier, 0, 6),
				newToken(token.EndOfFile, 6, 6),
			},
		},
		{
			name:               "When scanning an unknown identifier followed by a delimiter, the returned tokens are correct.",
			textInput:          "unknown}",
			amountOfCallsInput: 3,
			wantTokensOutput: []token.Token{
				newToken(token.Identifier, 0, 7),
				newToken(token.RightBrace, 7, 8),
				newToken(token.EndOfFile, 8, 8),
			},
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
			name:               "When scanning an identifier, the returned token is correct.",
			textInput:          "letters",
			amountOfCallsInput: 2,
			wantTokensOutput: []token.Token{
				newToken(token.Identifier, 0, 7),
				newToken(token.EndOfFile, 7, 7),
			},
		},
		{
			name:               "When scanning a character literal with a supported escape, the returned token is correct.",
			textInput:          "'\\t'",
			amountOfCallsInput: 2,
			wantTokensOutput: []token.Token{
				newToken(token.Character, 0, 4),
				newToken(token.EndOfFile, 4, 4),
			},
		},
		{
			name:               "When scanning an unexpected byte, the returned error is correct.",
			textInput:          "@",
			amountOfCallsInput: 1,
			wantErrorOutput:    "DSL Scanner: Unexpected byte '@'.",
		},
		{
			name:               "When scanning an identifier with a digit, the returned error is correct.",
			textInput:          "decimal1",
			amountOfCallsInput: 1,
			wantErrorOutput:    "DSL Scanner: Unexpected identifier \"decimal1\".",
		},
		{
			name:               "When scanning an identifier with a digit followed by a delimiter, the returned error is correct.",
			textInput:          "decimal1}",
			amountOfCallsInput: 1,
			wantErrorOutput:    "DSL Scanner: Unexpected identifier \"decimal1\".",
		},
		{
			name:               "When scanning an identifier with a digit followed by more identifier bytes, the returned error is correct.",
			textInput:          "decimal1x",
			amountOfCallsInput: 1,
			wantErrorOutput:    "DSL Scanner: Unexpected identifier \"decimal1x\".",
		},
		{
			name:               "When scanning an identifier that starts with an underscore, the returned error is correct.",
			textInput:          "_temp",
			amountOfCallsInput: 1,
			wantErrorOutput:    "DSL Scanner: Unexpected byte '_'.",
		},
		{
			name:               "When scanning a trailing slash at end of file, the returned error is correct.",
			textInput:          "/",
			amountOfCallsInput: 1,
			wantErrorOutput:    "DSL Scanner: Unexpected byte '/'.",
		},
		{
			name:               "When scanning a single dot, the returned error is correct.",
			textInput:          ".",
			amountOfCallsInput: 1,
			wantErrorOutput:    "DSL Scanner: Unexpected byte '.'.",
		},
		{
			name:               "When scanning a dot that is not followed by another dot, the returned error is correct.",
			textInput:          ".a",
			amountOfCallsInput: 1,
			wantErrorOutput:    "DSL Scanner: Unexpected byte '.'.",
		},
		{
			name:               "When scanning a character literal with no character, the returned error is correct.",
			textInput:          "''",
			amountOfCallsInput: 1,
			wantErrorOutput:    "DSL Scanner: Unexpected character literal.",
		},
		{
			name:               "When scanning a character literal with no character before end of file, the returned error is correct.",
			textInput:          "'",
			amountOfCallsInput: 1,
			wantErrorOutput:    "DSL Scanner: Unterminated character literal.",
		},
		{
			name:               "When scanning an unterminated character literal, the returned error is correct.",
			textInput:          "'a",
			amountOfCallsInput: 1,
			wantErrorOutput:    "DSL Scanner: Unterminated character literal.",
		},
		{
			name:               "When scanning a character literal with too many characters, the returned error is correct.",
			textInput:          "'ab'",
			amountOfCallsInput: 1,
			wantErrorOutput:    "DSL Scanner: Unexpected character literal.",
		},
		{
			name:               "When scanning a character literal with an unsupported escape, the returned error is correct.",
			textInput:          "'\\x'",
			amountOfCallsInput: 1,
			wantErrorOutput:    "DSL Scanner: Unexpected escape sequence \"\\\\x\".",
		},
		{
			name:               "When scanning a character literal that ends after an escape prefix, the returned error is correct.",
			textInput:          "'\\",
			amountOfCallsInput: 1,
			wantErrorOutput:    "DSL Scanner: Unterminated character literal.",
		},
		{
			name:               "When scanning a number, the returned token is correct.",
			textInput:          "123",
			amountOfCallsInput: 2,
			wantTokensOutput: []token.Token{
				newToken(token.Number, 0, 3),
				newToken(token.EndOfFile, 3, 3),
			},
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
			nWantTokensOutput := normalizeTokens(testCase.wantTokensOutput)
			nGotTokensOutput := normalizeTokens(gotTokensOutput)

			claim.DeepEqual(t, testCase.name, nWantTokensOutput, nGotTokensOutput, TokensLabel)

			if testCase.wantErrorOutput == "" {
				claim.Equal(t, testCase.name, false, gotErrorOutput != nil, HasErrorLabel)

				return
			}

			claim.Equal(t, testCase.name, false, gotErrorOutput == nil, HasErrorLabel)
			claim.Equal(t, testCase.name, testCase.wantErrorOutput, gotErrorOutput.Error(), ErrorLabel)
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
			name:                   "When the input reader fails while scanning an invalid identifier, the returned error is propagated.",
			readerInput:            newErrorAfterNReader("decimal1x", 8, errors.New("boom")),
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
			name:                   "When the input reader fails while scanning a character literal, the returned error is propagated.",
			readerInput:            newErrorAfterNReader("'a", 2, errors.New("boom")),
			amountOfTokensInput:    1,
			wantHasErrorOutput:     true,
			wantErrorMessageOutput: "boom",
		},
		{
			name:                   "When the input reader fails while scanning an escaped byte in a character literal, the returned error is propagated.",
			readerInput:            newErrorAfterNReader("'\\", 2, errors.New("boom")),
			amountOfTokensInput:    1,
			wantHasErrorOutput:     true,
			wantErrorMessageOutput: "boom",
		},
		{
			name:                   "When the input reader fails before any character literal value is read, the returned error is propagated.",
			readerInput:            newErrorAfterNReader("'", 1, errors.New("boom")),
			amountOfTokensInput:    1,
			wantHasErrorOutput:     true,
			wantErrorMessageOutput: "boom",
		},
		{
			name:                   "When the input reader fails while scanning a number, the returned error is propagated.",
			readerInput:            newErrorAfterNReader("12", 2, errors.New("boom")),
			amountOfTokensInput:    1,
			wantHasErrorOutput:     true,
			wantErrorMessageOutput: "boom",
		},
		{
			name:                   "When the input reader fails while scanning dot-dot, the returned error is propagated.",
			readerInput:            newErrorAfterNReader(".", 1, errors.New("boom")),
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
			nWantTokensOutput := normalizeTokens(testCase.wantTokensOutput)
			nGotTokensOutput := normalizeTokens(gotTokensOutput)

			claim.DeepEqual(t, testCase.name, nWantTokensOutput, nGotTokensOutput, TokensLabel)
			claim.Equal(t, testCase.name, testCase.wantHasErrorOutput, gotErrorOutput != nil, HasErrorLabel)

			if gotErrorOutput != nil {
				claim.Equal(t, testCase.name, testCase.wantErrorMessageOutput, gotErrorOutput.Error(), ErrorLabel)
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

func normalizeTokens(tokensInput []token.Token) []token.Token {
	if tokensInput == nil {
		return []token.Token{}
	}

	return tokensInput
}

func benchmarkText(blockCountInput int) string {
	var sb strings.Builder

	for range blockCountInput {
		sb.WriteString("scope {\n  include \"**/*.txt\"\n  exclude \"**/vendor/**\"\n  charset {\n    letters = ('a'..'z' | 'A'..'Z')+\n    decimal_numbers = ('0'..'9'){1,8}\n    space = ' ' | '\\t'?\n    optional = ('x' | 'y')*\n  }\n  // comment\n}\n")
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
