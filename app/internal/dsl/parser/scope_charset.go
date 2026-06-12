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

func (p *parser) parseCharset(tok token.Token) (ast.Charset, error) {
	p.advance()

	if _, err := p.expect(token.LeftBrace); err != nil {
		return ast.Charset{}, err
	}

	charsetOutput := ast.Charset{
		Range: position.NewRange(tok.Range.Start, tok.Range.End),
	}

	for p.current().Kind != token.RightBrace {
		memberOutput, memberErr := p.parseCharsetMember()

		if memberErr != nil {
			return ast.Charset{}, memberErr
		}

		charsetOutput.Members = append(charsetOutput.Members, memberOutput)
	}

	rightBraceToken := p.current()
	p.advance()
	charsetOutput.Range = position.NewRange(tok.Range.Start, rightBraceToken.Range.End)

	return charsetOutput, nil
}

func (p *parser) parseCharsetMember() (ast.CharsetMember, error) {
	nameToken, err := p.expect(token.Identifier)

	if err != nil {
		return ast.CharsetMember{}, err
	}

	if _, err := p.expect(token.Equal); err != nil {
		return ast.CharsetMember{}, err
	}

	valueOutput, err := p.parseExpression(false)

	if err != nil {
		return ast.CharsetMember{}, err
	}

	return ast.CharsetMember{
		Range: position.NewRange(nameToken.Range.Start, expressionRange(valueOutput).End),
		Value: valueOutput,
	}, nil
}
