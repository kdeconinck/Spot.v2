// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

// Verify the public API of the position package.
package position_test

import (
	"testing"

	"github.com/kdeconinck/spot/internal/position"
	"github.com/kdeconinck/spot/internal/qa/claim"
)

// Labels, used in the different assertion methods.
const (
	PositionLabel = "Position"
	RangeLabel    = "Range"
	LocationLabel = "Location"
	IsValidLabel  = "Is valid?"
)

func Test_NewPosition(t *testing.T) {
	t.Parallel()

	for tcName, tc := range map[string]struct {
		offsetInput position.Position
		want        position.Position
	}{
		"When the offset is non-negative, the returned value is correct.": {
			offsetInput: 7,
			want:        position.Position(7),
		},
		"When the offset is negative, the returned value is correct.": {
			offsetInput: -1,
			want:        position.InvalidPosition,
		},
	} {
		t.Run(tcName, func(t *testing.T) {
			t.Parallel()

			// Act.
			got := position.NewPosition(int(tc.offsetInput))

			// Assert.
			claim.Equal(t, tcName, tc.want, got, PositionLabel)
		})
	}
}

func Test_Position_IsValid(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name          string
		positionInput position.Position
		want          bool
	}{
		{
			name:          "When the position is invalid, the returned value is correct.",
			positionInput: position.NewPosition(-1),
			want:          false,
		},
		{
			name:          "When the position is zero, the returned value is correct.",
			positionInput: position.NewPosition(0),
			want:          true,
		},
		{
			name:          "When the position is positive, the returned value is correct.",
			positionInput: position.NewPosition(9),
			want:          true,
		},
		{
			name:          "When the position is the invalid position sentinel, the returned value is correct.",
			positionInput: position.InvalidPosition,
			want:          false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Act.
			got := testCase.positionInput.IsValid()

			// Assert.
			claim.Equal(t, testCase.name, testCase.want, got, IsValidLabel)
		})
	}
}
