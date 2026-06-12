// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

// Package validation validates semantic constraints of the DSL AST.
package validation

import (
	"fmt"

	"github.com/kdeconinck/spot/internal/position"
)

// Error describes a semantic validation error.
type Error struct {
	Range   position.Range
	Message string
}

// Error returns the user-facing error message.
func (errorInput Error) Error() string {
	return fmt.Sprintf("DSL Validation: %s.", errorInput.Message)
}
