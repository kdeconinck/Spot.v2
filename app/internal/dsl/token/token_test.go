// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

// Verify the public API of the token package.
package token_test

import (
	"testing"

	"github.com/kdeconinck/spot/internal/dsl/token"
	"github.com/kdeconinck/spot/internal/position"
	"github.com/kdeconinck/spot/internal/qa/claim"
)

// Labels, used in the different assertion methods.
const (
	StringLabel = "String"
	TokenLabel  = "Token"
)

func Test_New(t *testing.T) {
	t.Parallel()

	for tcName, tc := range map[string]struct {
		kindInput  token.Kind
		rangeInput position.Range
		want       token.Token
	}{
		"When creating a token, the returned token is correct.": {
			kindInput: token.Scope,
			rangeInput: position.NewRange(
				position.NewPosition(0),
				position.NewPosition(5)),
			want: token.Token{
				Kind: token.Scope,
				Range: position.Range{
					Start: position.Position(0),
					End:   position.Position(5),
				},
			},
		},
	} {
		t.Run(tcName, func(t *testing.T) {
			t.Parallel()

			// Act.
			got := token.New(tc.kindInput, tc.rangeInput)

			// Assert.
			claim.Equal(t, tcName, tc.want, got, TokenLabel)
		})
	}
}
