// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

// Package ast provides the abstract syntax tree for the DSL.
package ast

// Expression represents a DSL expression node.
type Expression interface {
	isExpression()
}
