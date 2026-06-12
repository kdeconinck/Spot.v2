// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

package validation

import (
	"fmt"
	"strings"

	"github.com/kdeconinck/spot/internal/dsl/ast"
	"github.com/kdeconinck/spot/internal/position"
)

// Validate validates the parsed DSL tree against semantic rules.
func Validate(textInput string, scopeInput ast.Scope) []Error {
	validatorInput := validator{
		text: textInput,
	}

	charsetNamesOutput := make(map[string]position.Range)

	if scopeInput.Charset != nil {
		for _, memberInput := range scopeInput.Charset.Members {
			nameOutput := validatorInput.memberName(memberInput)

			if firstRangeOutput, ok := charsetNamesOutput[nameOutput]; ok {
				validatorInput.addError(
					memberInput.Range,
					fmt.Sprintf("Duplicate charset member %q, first declared at %d", nameOutput, firstRangeOutput.Start))
			} else {
				charsetNamesOutput[nameOutput] = memberInput.Range
			}
		}

		for _, memberInput := range scopeInput.Charset.Members {
			validatorInput.validateExpression(memberInput.Value, false, charsetNamesOutput)
		}
	}

	if scopeInput.Lexer != nil {
		lexerNamesOutput := make(map[string]position.Range)

		for _, ruleInput := range scopeInput.Lexer.Rules {
			nameOutput := validatorInput.ruleName(ruleInput)

			if firstRangeOutput, ok := lexerNamesOutput[nameOutput]; ok {
				validatorInput.addError(
					ruleInput.Range,
					fmt.Sprintf("Duplicate lexer rule %q, first declared at %d", nameOutput, firstRangeOutput.Start))
			} else {
				lexerNamesOutput[nameOutput] = ruleInput.Range
			}
		}

		for _, ruleInput := range scopeInput.Lexer.Rules {
			validatorInput.validateExpression(ruleInput.Value, true, charsetNamesOutput)
		}
	}

	return validatorInput.errors
}

type validator struct {
	text   string
	errors []Error
}

func (v *validator) addError(rangeInput position.Range, messageInput string) {
	v.errors = append(v.errors, Error{
		Range:   rangeInput,
		Message: messageInput,
	})
}

func (v *validator) validateExpression(expressionInput ast.Expression, allowReferencesInput bool, charsetNamesInput map[string]position.Range) {
	switch value := expressionInput.(type) {
	case ast.Reference:
		referenceNameOutput := v.referenceName(value)

		if !allowReferencesInput {
			v.addError(value.Range, fmt.Sprintf("Charset expressions cannot reference %q", referenceNameOutput))

			return
		}

		if _, ok := charsetNamesInput[referenceNameOutput]; !ok {
			v.addError(value.Range, fmt.Sprintf("Unknown charset reference %q", referenceNameOutput))
		}

	case ast.StringLiteral:
		return

	case ast.CharacterLiteral:
		return

	case ast.CharacterRange:
		startOutput := v.characterValue(value.Start)
		endOutput := v.characterValue(value.End)

		if startOutput > endOutput {
			v.addError(value.Range, fmt.Sprintf("Invalid character range %q..%q", startOutput, endOutput))
		}

	case ast.Alternation:
		for _, expressionValue := range value.Expressions {
			v.validateExpression(expressionValue, allowReferencesInput, charsetNamesInput)
		}

	case ast.Concatenation:
		for _, expressionValue := range value.Expressions {
			v.validateExpression(expressionValue, allowReferencesInput, charsetNamesInput)
		}

	case ast.Repetition:
		if value.Maximum != nil && value.Minimum > *value.Maximum {
			v.addError(value.Range, fmt.Sprintf("Invalid repetition bounds %d,%d", value.Minimum, *value.Maximum))
		}

		v.validateExpression(value.Expression, allowReferencesInput, charsetNamesInput)
	}
}

func (v *validator) memberName(memberInput ast.CharsetMember) string {
	return v.nameBeforeExpression(memberInput.Range.Start, expressionRange(memberInput.Value).Start)
}

func (v *validator) ruleName(ruleInput ast.LexerRule) string {
	return v.nameBeforeExpression(ruleInput.Range.Start, expressionRange(ruleInput.Value).Start)
}

func (v *validator) referenceName(referenceInput ast.Reference) string {
	return v.text[referenceInput.Range.Start:referenceInput.Range.End]
}

func (v *validator) nameBeforeExpression(startInput, expressionStartInput position.Position) string {
	textOutput := v.text[startInput:expressionStartInput]
	equalIndexOutput := strings.IndexByte(textOutput, '=')

	if equalIndexOutput >= 0 {
		textOutput = textOutput[:equalIndexOutput]
	}

	return strings.TrimSpace(textOutput)
}

func (v *validator) characterValue(characterInput ast.CharacterLiteral) rune {
	textOutput := v.text[characterInput.Range.Start:characterInput.Range.End]

	if len(textOutput) == 3 {
		return rune(textOutput[1])
	}

	switch textOutput[2] {
	case '"':
		return '"'

	case '\\':
		return '\\'

	case 'n':
		return '\n'

	case 'r':
		return '\r'

	case 't':
		return '\t'

	default:
		return rune(textOutput[2])
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
		panic("unreachable")
	}
}
