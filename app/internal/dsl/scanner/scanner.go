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

		case nextByte == '/':
			return s.scanLineComment(start)

		case isIdentifierByte(nextByte):
			return s.scanIdentifier(start, nextByte)

		default:
			return token.Token{}, fmt.Errorf("DSL Scanner: Unexpected byte %q.", nextByte)
		}
	}
}

func (s *Scanner) scanIdentifier(start position.Position, firstByteInput byte) (token.Token, error) {
	if firstByteInput != 's' {
		return s.scanUnknownIdentifier([]byte{firstByteInput})
	}

	identifierPrefixOutput := []byte{'s'}

	for _, wantByte := range "cope" {
		nextByte, err := s.readByte()

		switch {
		case err == io.EOF:
			return token.Token{}, fmt.Errorf("DSL Scanner: Unexpected identifier %q.", "scope"[:int(s.offset-start)])

		case err != nil:
			return token.Token{}, err

		case nextByte != byte(wantByte):
			identifierPrefixOutput = append(identifierPrefixOutput, nextByte)

			return s.scanUnknownIdentifier(identifierPrefixOutput)

		default:
			identifierPrefixOutput = append(identifierPrefixOutput, nextByte)
		}
	}

	nextByte, err := s.readByte()

	switch {
	case err == io.EOF:
		return token.New(token.Scope, position.NewRange(start, s.offset)), nil

	case err != nil:
		return token.Token{}, err

	case isIdentifierByte(nextByte):
		identifierPrefixOutput = append(identifierPrefixOutput, nextByte)

		return s.scanUnknownIdentifier(identifierPrefixOutput)

	default:
		s.unreadByte()

		return token.New(token.Scope, position.NewRange(start, s.offset)), nil
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

		case !isIdentifierByte(nextByte):
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

func isIdentifierByte(valueInput byte) bool { return valueInput >= 'a' && valueInput <= 'z' }
