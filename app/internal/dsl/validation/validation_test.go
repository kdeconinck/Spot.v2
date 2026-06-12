// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

package validation_test

import (
	"strings"
	"testing"

	"github.com/kdeconinck/spot/internal/dsl/parser"
	"github.com/kdeconinck/spot/internal/dsl/validation"
	"github.com/kdeconinck/spot/internal/position"
	"github.com/kdeconinck/spot/internal/qa/claim"
)

const ErrorsLabel = "Errors"

func Test_Validate(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name             string
		textInput        string
		wantErrorsOutput []validation.Error
	}{
		{
			name:      "When validating a semantically valid scope, no errors are returned.",
			textInput: "scope { charset { letter = 'a'..'z' } lexer { identifier = letter+ } }",
		},
		{
			name:      "When validating duplicate charset member names, the returned error is correct.",
			textInput: "scope { charset { letter = 'a' letter = 'b' } }",
			wantErrorsOutput: []validation.Error{
				{
					Range:   rangeOf("scope { charset { letter = 'a' letter = 'b' } }", "letter = 'b'"),
					Message: `Duplicate charset member "letter", first declared at 18`,
				},
			},
		},
		{
			name:      "When validating duplicate lexer rule names, the returned error is correct.",
			textInput: "scope { lexer { identifier = 'a' identifier = 'b' } }",
			wantErrorsOutput: []validation.Error{
				{
					Range:   rangeOf("scope { lexer { identifier = 'a' identifier = 'b' } }", "identifier = 'b'"),
					Message: `Duplicate lexer rule "identifier", first declared at 16`,
				},
			},
		},
		{
			name:      "When validating a charset reference, the returned error is correct.",
			textInput: "scope { charset { letter = other } }",
			wantErrorsOutput: []validation.Error{
				{
					Range:   rangeOf("scope { charset { letter = other } }", "other"),
					Message: `Charset expressions cannot reference "other"`,
				},
			},
		},
		{
			name:      "When validating an unknown lexer reference, the returned error is correct.",
			textInput: "scope { charset { letter = 'a' } lexer { identifier = missing+ } }",
			wantErrorsOutput: []validation.Error{
				{
					Range:   rangeOf("scope { charset { letter = 'a' } lexer { identifier = missing+ } }", "missing"),
					Message: `Unknown charset reference "missing"`,
				},
			},
		},
		{
			name:      "When validating a reversed charset range, the returned error is correct.",
			textInput: "scope { charset { letter = 'z'..'a' } }",
			wantErrorsOutput: []validation.Error{
				{
					Range:   rangeOf("scope { charset { letter = 'z'..'a' } }", "'z'..'a'"),
					Message: `Invalid character range 'z'..'a'`,
				},
			},
		},
		{
			name:      "When validating a reversed lexer range, the returned error is correct.",
			textInput: "scope { lexer { identifier = 'z'..'a' } }",
			wantErrorsOutput: []validation.Error{
				{
					Range:   rangeOf("scope { lexer { identifier = 'z'..'a' } }", "'z'..'a'"),
					Message: `Invalid character range 'z'..'a'`,
				},
			},
		},
		{
			name:      "When validating invalid charset repetition bounds, the returned error is correct.",
			textInput: "scope { charset { letter = 'a'{8,1} } }",
			wantErrorsOutput: []validation.Error{
				{
					Range:   rangeOf("scope { charset { letter = 'a'{8,1} } }", "'a'{8,1}"),
					Message: "Invalid repetition bounds 8,1",
				},
			},
		},
		{
			name:      "When validating invalid lexer repetition bounds, the returned error is correct.",
			textInput: "scope { lexer { identifier = 'a'{8,1} } }",
			wantErrorsOutput: []validation.Error{
				{
					Range:   rangeOf("scope { lexer { identifier = 'a'{8,1} } }", "'a'{8,1}"),
					Message: "Invalid repetition bounds 8,1",
				},
			},
		},
		{
			name:      "When validating multiple semantic errors, all errors are returned.",
			textInput: "scope { charset { letter = other letter = 'z'..'a' } lexer { identifier = missing{8,1} } }",
			wantErrorsOutput: []validation.Error{
				{
					Range:   rangeOf("scope { charset { letter = other letter = 'z'..'a' } lexer { identifier = missing{8,1} } }", "letter = 'z'..'a'"),
					Message: `Duplicate charset member "letter", first declared at 18`,
				},
				{
					Range:   rangeOf("scope { charset { letter = other letter = 'z'..'a' } lexer { identifier = missing{8,1} } }", "other"),
					Message: `Charset expressions cannot reference "other"`,
				},
				{
					Range:   rangeOf("scope { charset { letter = other letter = 'z'..'a' } lexer { identifier = missing{8,1} } }", "'z'..'a'"),
					Message: `Invalid character range 'z'..'a'`,
				},
				{
					Range:   rangeOf("scope { charset { letter = other letter = 'z'..'a' } lexer { identifier = missing{8,1} } }", "missing{8,1}"),
					Message: "Invalid repetition bounds 8,1",
				},
				{
					Range:   rangeOf("scope { charset { letter = other letter = 'z'..'a' } lexer { identifier = missing{8,1} } }", "missing"),
					Message: `Unknown charset reference "missing"`,
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			scopeOutput, err := parser.Parse(testCase.textInput)

			claim.Equal(t, testCase.name, false, err != nil, "Has parser error?")

			gotErrorsOutput := validation.Validate(testCase.textInput, scopeOutput)

			claim.DeepEqual(t, testCase.name, testCase.wantErrorsOutput, gotErrorsOutput, ErrorsLabel)
		})
	}
}

func rangeOf(textInput, snippetInput string) position.Range {
	startOutput := strings.Index(textInput, snippetInput)

	return position.NewRange(
		position.NewPosition(startOutput),
		position.NewPosition(startOutput+len(snippetInput)))
}
