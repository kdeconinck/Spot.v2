// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

// Package parser converts DSL source text into an abstract syntax tree.
package parser

import (
	"github.com/kdeconinck/spot/internal/dsl/ast"
	"github.com/kdeconinck/spot/internal/dsl/token"
	"github.com/kdeconinck/spot/internal/position"
)

func (p *parser) parseInclude(tok token.Token) (ast.Include, error) {
	p.advance()
	patternToken, err := p.expect(token.String)

	if err != nil {
		return ast.Include{}, err
	}

	return ast.Include{
		Range:   position.NewRange(tok.Range.Start, patternToken.Range.End),
		Pattern: p.stringLiteralFromToken(patternToken),
	}, nil
}
