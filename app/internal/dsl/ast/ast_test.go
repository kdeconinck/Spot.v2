// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

package ast_test

import (
	"github.com/kdeconinck/spot/internal/position"
)

func newRange(startInput, endInput int) position.Range {
	return position.NewRange(
		position.NewPosition(startInput), position.NewPosition(endInput))
}
