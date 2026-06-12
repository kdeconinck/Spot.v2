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

func (p *parser) parseScope() (ast.Scope, error) {
	scopeKeyword, err := p.expect(token.Scope)

	if err != nil {
		return ast.Scope{}, err
	}

	if _, err := p.expect(token.LeftBrace); err != nil {
		return ast.Scope{}, err
	}

	scopeOutput := ast.Scope{
		Range: position.NewRange(scopeKeyword.Range.Start, scopeKeyword.Range.End),
	}

	for p.current().Kind != token.RightBrace {
		switch p.current().Kind {
		case token.Include:
			includeOutput, includeErr := p.parseInclude(p.current())

			if includeErr != nil {
				return ast.Scope{}, includeErr
			}

			scopeOutput.Includes = append(scopeOutput.Includes, includeOutput)

		case token.Exclude:
			excludeOutput, excludeErr := p.parseExclude(p.current())

			if excludeErr != nil {
				return ast.Scope{}, excludeErr
			}

			scopeOutput.Excludes = append(scopeOutput.Excludes, excludeOutput)

		case token.Charset:
			charsetOutput, charsetErr := p.parseCharset(p.current())

			if charsetErr != nil {
				return ast.Scope{}, charsetErr
			}

			scopeOutput.Charset = &charsetOutput

		case token.Lexer:
			lexerOutput, lexerErr := p.parseLexer(p.current())

			if lexerErr != nil {
				return ast.Scope{}, lexerErr
			}

			scopeOutput.Lexer = &lexerOutput

		default:
			return ast.Scope{}, p.unexpectedTokenError("scope member")
		}
	}

	rightBraceToken := p.current()
	p.advance()

	if _, err := p.expect(token.EndOfFile); err != nil {
		return ast.Scope{}, err
	}

	scopeOutput.Range = position.NewRange(scopeKeyword.Range.Start, rightBraceToken.Range.End)

	return scopeOutput, nil
}
