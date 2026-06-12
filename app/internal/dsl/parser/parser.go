// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

// Package parser converts DSL source text into an abstract syntax tree.
package parser

import (
	"fmt"
	"strings"

	"github.com/kdeconinck/spot/internal/dsl/ast"
	"github.com/kdeconinck/spot/internal/dsl/scanner"
	"github.com/kdeconinck/spot/internal/dsl/token"
	"github.com/kdeconinck/spot/internal/position"
)

// Parse parses the DSL source text into an abstract syntax tree.
func Parse(textInput string) (ast.Scope, error) {
	tokensOutput, err := scanTokens(textInput)

	if err != nil {
		return ast.Scope{}, err
	}

	parserOutput := parser{
		text:   textInput,
		tokens: tokensOutput,
	}

	return parserOutput.parseScope()
}

type parser struct {
	text   string
	tokens []token.Token
	index  int
}

func (p *parser) advance() { p.index++ }

func (p *parser) expect(kindInput token.Kind) (token.Token, error) {
	if p.current().Kind != kindInput {
		return token.Token{}, p.unexpectedTokenError(kindInput.String())
	}

	currentOutput := p.current()

	if currentOutput.Kind != token.EndOfFile {
		p.advance()
	}

	return currentOutput, nil
}

func (p *parser) unexpectedTokenError(expectedInput string) error {
	return fmt.Errorf("DSL Parser: Expected %s, got %q.", expectedInput, p.current().Kind.String())
}

func (p *parser) lexeme(tokenInput token.Token) string {
	return p.text[tokenInput.Range.Start:tokenInput.Range.End]
}

func (p *parser) stringLiteralFromToken(tokenInput token.Token) ast.StringLiteral {
	return ast.StringLiteral{
		Range: tokenInput.Range,
	}
}

func (p *parser) characterLiteralFromToken(tokenInput token.Token) ast.CharacterLiteral {
	return ast.CharacterLiteral{
		Range: tokenInput.Range,
	}
}

func scanTokens(textInput string) ([]token.Token, error) {
	scannerInput := scanner.New(strings.NewReader(textInput))
	tokensOutput := make([]token.Token, 0, 32)

	for {
		nextTokenOutput, err := scannerInput.Next()

		if err != nil {
			return nil, err
		}

		if nextTokenOutput.Kind == token.LineComment {
			continue
		}

		tokensOutput = append(tokensOutput, nextTokenOutput)

		if nextTokenOutput.Kind == token.EndOfFile {
			return tokensOutput, nil
		}
	}
}

func (p *parser) current() token.Token { return p.tokens[p.index] }

func isExpressionStart(kindInput token.Kind, allowStringInput bool) bool {
	switch kindInput {
	case token.Identifier, token.Character, token.LeftParen:
		return true

	case token.String:
		return allowStringInput

	default:
		return false
	}
}

func expressionRange(expressionInput ast.Expression) position.Range {
	switch value := expressionInput.(type) {
	case ast.Reference:
		return value.Range

	case ast.StringLiteral:
		return value.Range

	case ast.CharacterLiteral:
		return value.Range

	case ast.CharacterRange:
		return value.Range

	case ast.Alternation:
		return value.Range

	case ast.Concatenation:
		return value.Range

	case ast.Repetition:
		return value.Range

	default:
		panic("DSL Parser: Unreachable expression type.")
	}
}

func withExpressionRange(expressionInput ast.Expression, startInput, endInput position.Position) ast.Expression {
	expressionRangeOutput := position.NewRange(startInput, endInput)

	switch value := expressionInput.(type) {
	case ast.Reference:
		value.Range = expressionRangeOutput

		return value

	case ast.StringLiteral:
		value.Range = expressionRangeOutput

		return value

	case ast.CharacterLiteral:
		value.Range = expressionRangeOutput

		return value

	case ast.CharacterRange:
		value.Range = expressionRangeOutput

		return value

	case ast.Alternation:
		value.Range = expressionRangeOutput

		return value

	case ast.Concatenation:
		value.Range = expressionRangeOutput

		return value

	case ast.Repetition:
		value.Range = expressionRangeOutput

		return value

	default:
		panic("unreachable")
	}
}
