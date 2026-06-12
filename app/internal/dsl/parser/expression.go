// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

// Package parser converts DSL source text into an abstract syntax tree.
package parser

import (
	"fmt"
	"strconv"

	"github.com/kdeconinck/spot/internal/dsl/ast"
	"github.com/kdeconinck/spot/internal/dsl/token"
	"github.com/kdeconinck/spot/internal/position"
)

func (p *parser) parseExpression(allowStringInput bool) (ast.Expression, error) {
	return p.parseAlternation(allowStringInput)
}

func (p *parser) parseAlternation(allowStringInput bool) (ast.Expression, error) {
	expressionsOutput := make([]ast.Expression, 0, 1)
	firstOutput, err := p.parseConcatenation(allowStringInput)

	if err != nil {
		return nil, err
	}

	expressionsOutput = append(expressionsOutput, firstOutput)

	for p.current().Kind == token.Pipe {
		p.advance()

		nextOutput, nextErr := p.parseConcatenation(allowStringInput)

		if nextErr != nil {
			return nil, nextErr
		}

		expressionsOutput = append(expressionsOutput, nextOutput)
	}

	if len(expressionsOutput) == 1 {
		return expressionsOutput[0], nil
	}

	return ast.Alternation{
		Range:       position.NewRange(expressionRange(expressionsOutput[0]).Start, expressionRange(expressionsOutput[len(expressionsOutput)-1]).End),
		Expressions: expressionsOutput,
	}, nil
}

func (p *parser) parseConcatenation(allowStringInput bool) (ast.Expression, error) {
	expressionsOutput := make([]ast.Expression, 0, 1)

	for isExpressionStart(p.current().Kind, allowStringInput) {
		if p.current().Kind == token.Identifier && p.tokens[p.index+1].Kind == token.Equal {
			break
		}

		nextOutput, err := p.parseRepeatedExpression(allowStringInput)

		if err != nil {
			return nil, err
		}

		expressionsOutput = append(expressionsOutput, nextOutput)
	}

	if len(expressionsOutput) == 0 {
		if !allowStringInput && p.current().Kind == token.String {
			return nil, fmt.Errorf("DSL Parser: Expected charset expression, got %q.", p.current().Kind.String())
		}

		return nil, p.unexpectedTokenError("expression")
	}

	if len(expressionsOutput) == 1 {
		return expressionsOutput[0], nil
	}

	return ast.Concatenation{
		Range:       position.NewRange(expressionRange(expressionsOutput[0]).Start, expressionRange(expressionsOutput[len(expressionsOutput)-1]).End),
		Expressions: expressionsOutput,
	}, nil
}

func (p *parser) parseRepeatedExpression(allowStringInput bool) (ast.Expression, error) {
	expressionOutput, err := p.parseRangedExpression(allowStringInput)

	if err != nil {
		return nil, err
	}

	for {
		switch p.current().Kind {
		case token.Star:
			starToken := p.current()
			p.advance()

			expressionOutput = ast.Repetition{
				Range:      position.NewRange(expressionRange(expressionOutput).Start, starToken.Range.End),
				Expression: expressionOutput,
				Minimum:    0,
				Maximum:    nil,
			}

		case token.Plus:
			plusToken := p.current()
			p.advance()

			expressionOutput = ast.Repetition{
				Range:      position.NewRange(expressionRange(expressionOutput).Start, plusToken.Range.End),
				Expression: expressionOutput,
				Minimum:    1,
				Maximum:    nil,
			}

		case token.Question:
			questionToken := p.current()
			p.advance()

			maximumOutput := 1
			expressionOutput = ast.Repetition{
				Range:      position.NewRange(expressionRange(expressionOutput).Start, questionToken.Range.End),
				Expression: expressionOutput,
				Minimum:    0,
				Maximum:    &maximumOutput,
			}

		case token.LeftBrace:
			repetitionOutput, repetitionErr := p.parseBoundedRepetition(expressionOutput)

			if repetitionErr != nil {
				return nil, repetitionErr
			}

			expressionOutput = repetitionOutput

		default:
			return expressionOutput, nil
		}
	}
}

func (p *parser) parseRangedExpression(allowStringInput bool) (ast.Expression, error) {
	leftOutput, err := p.parsePrimary(allowStringInput)

	if err != nil {
		return nil, err
	}

	if p.current().Kind != token.DotDot {
		return leftOutput, nil
	}

	leftCharacterOutput, ok := leftOutput.(ast.CharacterLiteral)

	if !ok {
		return nil, fmt.Errorf("DSL Parser: Expected character literal before %q.", token.DotDot.String())
	}

	p.advance()

	rightOutput, err := p.parsePrimary(allowStringInput)

	if err != nil {
		return nil, err
	}

	rightCharacterOutput, ok := rightOutput.(ast.CharacterLiteral)

	if !ok {
		return nil, fmt.Errorf("DSL Parser: Expected character literal after %q.", token.DotDot.String())
	}

	return ast.CharacterRange{
		Range: position.NewRange(leftCharacterOutput.Range.Start, rightCharacterOutput.Range.End),
		Start: leftCharacterOutput,
		End:   rightCharacterOutput,
	}, nil
}

func (p *parser) parsePrimary(allowStringInput bool) (ast.Expression, error) {
	switch p.current().Kind {
	case token.Identifier:
		identifierToken := p.current()
		p.advance()

		return ast.Reference{
			Range: identifierToken.Range,
		}, nil

	case token.String:
		if !allowStringInput {
			return nil, fmt.Errorf("DSL Parser: Expected charset expression, got %q.", p.current().Kind.String())
		}

		stringToken := p.current()
		p.advance()

		return p.stringLiteralFromToken(stringToken), nil

	case token.Character:
		characterToken := p.current()
		p.advance()

		return p.characterLiteralFromToken(characterToken), nil

	case token.LeftParen:
		leftParenToken := p.current()
		p.advance()

		expressionOutput, err := p.parseExpression(allowStringInput)

		if err != nil {
			return nil, err
		}

		rightParenToken, err := p.expect(token.RightParen)

		if err != nil {
			return nil, err
		}

		return withExpressionRange(expressionOutput, leftParenToken.Range.Start, rightParenToken.Range.End), nil

	default:
		return nil, p.unexpectedTokenError("expression")
	}
}

func (p *parser) parseBoundedRepetition(expressionInput ast.Expression) (ast.Repetition, error) {
	p.advance()

	minimumToken, err := p.expect(token.Number)
	if err != nil {
		return ast.Repetition{}, err
	}

	minimumOutput, err := strconv.Atoi(p.lexeme(minimumToken))
	if err != nil {
		return ast.Repetition{}, err
	}

	if p.current().Kind == token.RightBrace {
		rightBraceToken := p.current()
		p.advance()

		maximumOutput := minimumOutput

		return ast.Repetition{
			Range:      position.NewRange(expressionRange(expressionInput).Start, rightBraceToken.Range.End),
			Expression: expressionInput,
			Minimum:    minimumOutput,
			Maximum:    &maximumOutput,
		}, nil
	}

	if _, err := p.expect(token.Comma); err != nil {
		return ast.Repetition{}, err
	}

	if p.current().Kind == token.RightBrace {
		rightBraceToken := p.current()
		p.advance()

		return ast.Repetition{
			Range:      position.NewRange(expressionRange(expressionInput).Start, rightBraceToken.Range.End),
			Expression: expressionInput,
			Minimum:    minimumOutput,
			Maximum:    nil,
		}, nil
	}

	maximumToken, err := p.expect(token.Number)
	if err != nil {
		return ast.Repetition{}, err
	}

	maximumOutput, err := strconv.Atoi(p.lexeme(maximumToken))
	if err != nil {
		return ast.Repetition{}, err
	}

	rightBraceToken, err := p.expect(token.RightBrace)
	if err != nil {
		return ast.Repetition{}, err
	}

	return ast.Repetition{
		Range:      position.NewRange(expressionRange(expressionInput).Start, rightBraceToken.Range.End),
		Expression: expressionInput,
		Minimum:    minimumOutput,
		Maximum:    &maximumOutput,
	}, nil
}
