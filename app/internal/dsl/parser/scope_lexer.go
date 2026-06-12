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

func (p *parser) parseLexer(tok token.Token) (ast.Lexer, error) {
	p.advance()

	if _, err := p.expect(token.LeftBrace); err != nil {
		return ast.Lexer{}, err
	}

	lexerOutput := ast.Lexer{
		Range: position.NewRange(tok.Range.Start, tok.Range.End),
	}

	for p.current().Kind != token.RightBrace {
		ruleOutput, ruleErr := p.parseLexerRule()

		if ruleErr != nil {
			return ast.Lexer{}, ruleErr
		}

		lexerOutput.Rules = append(lexerOutput.Rules, ruleOutput)
	}

	rightBraceToken := p.current()
	p.advance()
	lexerOutput.Range = position.NewRange(tok.Range.Start, rightBraceToken.Range.End)

	return lexerOutput, nil
}

func (p *parser) parseLexerRule() (ast.LexerRule, error) {
	nameToken, err := p.expect(token.Identifier)

	if err != nil {
		return ast.LexerRule{}, err
	}

	if _, err := p.expect(token.Equal); err != nil {
		return ast.LexerRule{}, err
	}

	valueOutput, err := p.parseExpression(true)

	if err != nil {
		return ast.LexerRule{}, err
	}

	return ast.LexerRule{
		Range: position.NewRange(nameToken.Range.Start, expressionRange(valueOutput).End),
		Value: valueOutput,
	}, nil
}
