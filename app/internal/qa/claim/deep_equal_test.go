// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

package claim_test

import (
	"testing"

	"github.com/kdeconinck/spot/internal/qa/claim"
)

func Test_DeepEqual(t *testing.T) {
	t.Parallel()

	type sample struct {
		Name   string
		Values []int
	}

	for testName, testCase := range map[string]struct {
		wantInput, gotInput              any
		wantMsg                          string
		wantHelperCalls, wantFatalfCalls int
	}{
		"When the compared values are deeply equal, no failure is reported.": {
			wantInput: sample{
				Name: "ok", Values: []int{1, 2},
			},
			gotInput: sample{
				Name: "ok", Values: []int{1, 2},
			},
			wantMsg:         "",
			wantHelperCalls: 1,
			wantFatalfCalls: 0,
		},
		"When the compared values are not deeply equal, a failure is reported.": {
			wantInput: sample{
				Name: "ok", Values: []int{1, 2},
			},
			gotInput: sample{
				Name: "ok", Values: []int{1, 3},
			},
			wantHelperCalls: 1,
			wantFatalfCalls: 1,
			wantMsg:         "\n\nTest name:            Deep equality.\n\033[32mExpected (Structure): {ok [1 2]}\033[0m\n\033[31mActual (Structure):   {ok [1 3]}\033[0m\n\n",
		},
	} {
		t.Run(testName, func(t *testing.T) {
			t.Parallel()

			// Arrange.
			spy := new(tbSpy)

			// Act.
			claim.DeepEqual(spy, "Deep equality.", testCase.wantInput, testCase.gotInput, "Structure")

			// Assert.
			spy.verifyFailure(t, testCase.wantMsg, testCase.wantHelperCalls, testCase.wantFatalfCalls)
		})
	}
}

func benchmark_DeepEqualFailure(b *testing.B) {
	b.Helper()

	tb := new(tbNoop)

	type sample struct {
		Name   string
		Values []int
	}

	want := sample{Name: "ok", Values: []int{1, 2}}
	got := sample{Name: "ok", Values: []int{1, 3}}

	for b.Loop() {
		claim.DeepEqual(tb, "Benchmark deep equality.", want, got, "Label")
	}
}

func benchmark_DeepEqualSuccess(b *testing.B) {
	b.Helper()

	tb := new(tbNoop)

	type sample struct {
		Name   string
		Values []int
	}

	want := sample{Name: "ok", Values: []int{1, 2}}
	got := sample{Name: "ok", Values: []int{1, 2}}

	for b.Loop() {
		claim.DeepEqual(tb, "Benchmark deep equality.", want, got, "Label")
	}
}

func Benchmark_DeepEqual_Failure(b *testing.B) { benchmark_DeepEqualFailure(b) }
func Benchmark_DeepEqual_Success(b *testing.B) { benchmark_DeepEqualSuccess(b) }
