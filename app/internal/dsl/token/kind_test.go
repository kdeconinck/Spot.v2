// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

// Verify the public API of the token package.
package token_test

import (
	"testing"

	"github.com/kdeconinck/spot/internal/dsl/token"
	"github.com/kdeconinck/spot/internal/qa/claim"
)

func Test_Kind_String(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name      string
		kindInput token.Kind
		want      string
	}{
		{
			name:      "When using the invalid token kind, the returned string is correct.",
			kindInput: token.Invalid,
			want:      "invalid",
		},
		{
			name:      "When using the end of file token kind, the returned string is correct.",
			kindInput: token.EndOfFile,
			want:      "end-of-file",
		},
		{
			name:      "When using the scope token kind, the returned string is correct.",
			kindInput: token.Scope,
			want:      "scope",
		},
		{
			name:      "When using the include token kind, the returned string is correct.",
			kindInput: token.Include,
			want:      "include",
		},
		{
			name:      "When using the exclude token kind, the returned string is correct.",
			kindInput: token.Exclude,
			want:      "exclude",
		},
		{
			name:      "When using the left brace token kind, the returned string is correct.",
			kindInput: token.LeftBrace,
			want:      "left-brace",
		},
		{
			name:      "When using the right brace token kind, the returned string is correct.",
			kindInput: token.RightBrace,
			want:      "right-brace",
		},
		{
			name:      "When using the string token kind, the returned string is correct.",
			kindInput: token.String,
			want:      "string",
		},
		{
			name:      "When using the line comment token kind, the returned string is correct.",
			kindInput: token.LineComment,
			want:      "line-comment",
		},
		{
			name:      "When using an unknown token kind, the returned string is correct.",
			kindInput: token.Kind(255),
			want:      "255",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Act.
			got := testCase.kindInput.String()

			// Assert.
			claim.Equal(t, testCase.name, testCase.want, got, StringLabel)
		})
	}
}
