// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

package position_test

import (
	"testing"

	"github.com/kdeconinck/spot/internal/position"
	"github.com/kdeconinck/spot/internal/qa/claim"
)

func Test_NewRange(t *testing.T) {
	t.Parallel()

	for tcName, tc := range map[string]struct {
		startInput position.Position
		endInput   position.Position
		want       position.Range
	}{
		"When creating a range, the returned range preserves both endpoints.": {
			startInput: position.NewPosition(3),
			endInput:   position.NewPosition(8),
			want: position.Range{
				Start: position.Position(3),
				End:   position.Position(8),
			},
		},
	} {
		t.Run(tcName, func(t *testing.T) {
			t.Parallel()

			// Act.
			got := position.NewRange(tc.startInput, tc.endInput)

			// Assert.
			claim.Equal(t, tcName, tc.want, got, RangeLabel)
		})
	}
}

func Test_Range_IsValid(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		rangeInput position.Range
		want       bool
	}{
		{
			name: "When the range has increasing endpoints, the returned value is true.",
			rangeInput: position.NewRange(
				position.NewPosition(2),
				position.NewPosition(5)),
			want: true,
		},
		{
			name: "When the range is empty, the returned value is true.",
			rangeInput: position.NewRange(
				position.NewPosition(4),
				position.NewPosition(4)),
			want: true,
		},
		{
			name:       "When the range start is invalid, the returned value is false.",
			rangeInput: position.NewRange(position.NewPosition(-1), position.NewPosition(4)),
			want:       false,
		},
		{
			name:       "When the range end is invalid, the returned value is false.",
			rangeInput: position.NewRange(position.NewPosition(1), position.NewPosition(-1)),
			want:       false,
		},
		{
			name: "When the range end precedes the start, the returned value is false.",
			rangeInput: position.NewRange(
				position.NewPosition(5),
				position.NewPosition(2)),
			want: false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Act.
			got := testCase.rangeInput.IsValid()

			// Assert.
			claim.Equal(t, testCase.name, testCase.want, got, IsValidLabel)
		})
	}
}
