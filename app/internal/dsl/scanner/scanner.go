// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

// Package scanner provides a streaming DSL scanner.
package scanner

import (
	"bufio"
	"fmt"
	"io"

	"github.com/kdeconinck/spot/internal/dsl/token"
	"github.com/kdeconinck/spot/internal/position"
)

// Scanner tokenizes DSL input from a stream.
type Scanner struct {
	reader        *bufio.Reader
	offset        position.Position
	endOfFile     token.Token
	hasEmittedEOF bool
}

// New returns a scanner for reader.
func New(reader io.Reader) *Scanner {
	return &Scanner{
		reader: bufio.NewReader(reader),
		offset: position.NewPosition(0),
		endOfFile: token.New(
			token.EndOfFile,
			position.NewRange(
				position.NewPosition(0), position.NewPosition(0))),
	}
}

// Next returns the next token from the input stream.
// When err is non-nil, the returned token must be ignored.
func (s *Scanner) Next() (token.Token, error) {
	if s.hasEmittedEOF {
		return s.endOfFile, nil
	}

	for {
		start := s.offset
		nextByte, err := s.readByte()

		if err == io.EOF {
			s.hasEmittedEOF = true
			s.endOfFile = token.New(token.EndOfFile, position.NewRange(start, start))

			return s.endOfFile, nil
		}

		if err != nil {
			return token.Token{}, err
		}

		if isWhitespace(nextByte) {
			continue
		}

		switch {
		case nextByte == '{':
			return token.New(token.LeftBrace, position.NewRange(start, s.offset)), nil

		case nextByte == '}':
			return token.New(token.RightBrace, position.NewRange(start, s.offset)), nil

		case nextByte == '(':
			return token.New(token.LeftParen, position.NewRange(start, s.offset)), nil

		case nextByte == ')':
			return token.New(token.RightParen, position.NewRange(start, s.offset)), nil

		case nextByte == '=':
			return token.New(token.Equal, position.NewRange(start, s.offset)), nil

		case nextByte == '|':
			return token.New(token.Pipe, position.NewRange(start, s.offset)), nil

		case nextByte == '*':
			return token.New(token.Star, position.NewRange(start, s.offset)), nil

		case nextByte == '+':
			return token.New(token.Plus, position.NewRange(start, s.offset)), nil

		case nextByte == '?':
			return token.New(token.Question, position.NewRange(start, s.offset)), nil

		case nextByte == ',':
			return token.New(token.Comma, position.NewRange(start, s.offset)), nil

		case nextByte == '.':
			return s.scanDotDot(start)

		case nextByte == '"':
			return s.scanString(start)

		case nextByte == '\'':
			return s.scanCharacter(start)

		case nextByte == '/':
			return s.scanLineComment(start)

		case isNumberByte(nextByte):
			return s.scanNumber(start)

		case isIdentByte(nextByte):
			return s.scanIdentifier(start, nextByte)

		default:
			return token.Token{}, fmt.Errorf("DSL Scanner: Unexpected byte %q.", nextByte)
		}
	}
}

func (s *Scanner) scanIdentifier(start position.Position, firstByteInput byte) (token.Token, error) {
	identifierPrefixOutput := []byte{firstByteInput}

	for {
		nextByte, err := s.readByte()

		switch {
		case err == io.EOF:
			return token.New(identifierKind(string(identifierPrefixOutput)), position.NewRange(start, s.offset)), nil

		case err != nil:
			return token.Token{}, err

		case isIdentBodyByte(nextByte):
			identifierPrefixOutput = append(identifierPrefixOutput, nextByte)

		case isNumberByte(nextByte):
			identifierPrefixOutput = append(identifierPrefixOutput, nextByte)

			return s.scanUnknownIdentifier(identifierPrefixOutput)

		default:
			s.unreadByte()

			return token.New(identifierKind(string(identifierPrefixOutput)), position.NewRange(start, s.offset)), nil
		}
	}
}

func (s *Scanner) scanUnknownIdentifier(identifierPrefixInput []byte) (token.Token, error) {
	for {
		nextByte, err := s.readByte()

		switch {
		case err == io.EOF:
			return token.Token{}, fmt.Errorf("DSL Scanner: Unexpected identifier %q.", string(identifierPrefixInput))

		case err != nil:
			return token.Token{}, err

		case !isIdentBodyByte(nextByte) && !isNumberByte(nextByte):
			return token.Token{}, fmt.Errorf("DSL Scanner: Unexpected identifier %q.", string(identifierPrefixInput))

		default:
			identifierPrefixInput = append(identifierPrefixInput, nextByte)
		}
	}
}

func (s *Scanner) scanLineComment(start position.Position) (token.Token, error) {
	nextByte, err := s.readByte()

	switch {
	case err == io.EOF:
		return token.Token{}, fmt.Errorf("DSL Scanner: Unexpected byte %q.", '/')

	case err != nil:
		return token.Token{}, err

	case nextByte != '/':
		return token.Token{}, fmt.Errorf("DSL Scanner: Unexpected byte %q.", '/')
	}

	for {
		nextByte, err := s.readByte()

		switch {
		case err == io.EOF:
			return token.New(token.LineComment, position.NewRange(start, s.offset)), nil

		case err != nil:
			return token.Token{}, err

		case nextByte == '\n':
			s.unreadByte()

			return token.New(token.LineComment, position.NewRange(start, s.offset)), nil
		}
	}
}

func (s *Scanner) scanString(start position.Position) (token.Token, error) {
	for {
		nextByte, err := s.readByte()

		switch {
		case err == io.EOF:
			return token.Token{}, fmt.Errorf("DSL Scanner: Unterminated string literal.")

		case err != nil:
			return token.Token{}, err

		case nextByte == '"':
			return token.New(token.String, position.NewRange(start, s.offset)), nil

		case nextByte == '\\':
			escapedByte, escapedErr := s.readByte()

			switch {
			case escapedErr == io.EOF:
				return token.Token{}, fmt.Errorf("DSL Scanner: Unterminated string literal.")

			case escapedErr != nil:
				return token.Token{}, escapedErr

			case !isEscapedByte(escapedByte):
				return token.Token{}, fmt.Errorf("DSL Scanner: Unexpected escape sequence %q.", "\\"+string(escapedByte))
			}
		}
	}
}

func (s *Scanner) scanCharacter(start position.Position) (token.Token, error) {
	nextByte, err := s.readByte()

	switch {
	case err == io.EOF:
		return token.Token{}, fmt.Errorf("DSL Scanner: Unterminated character literal.")

	case err != nil:
		return token.Token{}, err

	case nextByte == '\'':
		return token.Token{}, fmt.Errorf("DSL Scanner: Unexpected character literal.")

	case nextByte == '\\':
		escapedByte, escapedErr := s.readByte()

		switch {
		case escapedErr == io.EOF:
			return token.Token{}, fmt.Errorf("DSL Scanner: Unterminated character literal.")

		case escapedErr != nil:
			return token.Token{}, escapedErr

		case !isEscapedByte(escapedByte):
			return token.Token{}, fmt.Errorf("DSL Scanner: Unexpected escape sequence %q.", "\\"+string(escapedByte))
		}
	}

	closingByte, closingErr := s.readByte()

	switch {
	case closingErr == io.EOF:
		return token.Token{}, fmt.Errorf("DSL Scanner: Unterminated character literal.")

	case closingErr != nil:
		return token.Token{}, closingErr

	case closingByte != '\'':
		return token.Token{}, fmt.Errorf("DSL Scanner: Unexpected character literal.")

	default:
		return token.New(token.Character, position.NewRange(start, s.offset)), nil
	}
}

func (s *Scanner) scanNumber(start position.Position) (token.Token, error) {
	for {
		nextByte, err := s.readByte()

		switch {
		case err == io.EOF:
			return token.New(token.Number, position.NewRange(start, s.offset)), nil

		case err != nil:
			return token.Token{}, err

		case !isNumberByte(nextByte):
			s.unreadByte()

			return token.New(token.Number, position.NewRange(start, s.offset)), nil
		}
	}
}

func (s *Scanner) scanDotDot(start position.Position) (token.Token, error) {
	nextByte, err := s.readByte()

	switch {
	case err == io.EOF:
		return token.Token{}, fmt.Errorf("DSL Scanner: Unexpected byte %q.", '.')

	case err != nil:
		return token.Token{}, err

	case nextByte != '.':
		return token.Token{}, fmt.Errorf("DSL Scanner: Unexpected byte %q.", '.')

	default:
		return token.New(token.DotDot, position.NewRange(start, s.offset)), nil
	}
}

func (s *Scanner) readByte() (byte, error) {
	nextByte, err := s.reader.ReadByte()

	if err != nil {
		return 0, err
	}

	s.offset++

	return nextByte, nil
}

func (s *Scanner) unreadByte() {
	_ = s.reader.UnreadByte()
	s.offset--
}

func isWhitespace(valueInput byte) bool {
	switch valueInput {
	case ' ', '\t', '\n', '\r':
		return true

	default:
		return false
	}
}

func isEscapedByte(valueInput byte) bool {
	switch valueInput {
	case '"', '\\', 'n', 'r', 't':
		return true

	default:
		return false
	}
}

func isIdentByte(valueInput byte) bool     { return valueInput >= 'a' && valueInput <= 'z' }
func isIdentBodyByte(valueInput byte) bool { return isIdentByte(valueInput) || valueInput == '_' }
func isNumberByte(valueInput byte) bool    { return valueInput >= '0' && valueInput <= '9' }

func identifierKind(identifierInput string) token.Kind {
	switch identifierInput {
	case "scope":
		return token.Scope

	case "charset":
		return token.Charset

	case "lexer":
		return token.Lexer

	case "include":
		return token.Include

	case "exclude":
		return token.Exclude

	default:
		return token.Identifier
	}
}
