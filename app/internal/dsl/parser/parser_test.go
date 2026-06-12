// =====================================================================================================================
// == LICENSE:                 Copyright (c) 2026 Kevin De Coninck.
// == SPDX-License-Identifier: LicenseRef-PolyForm-Noncommercial-1.0.0
// =====================================================================================================================

package parser_test

import (
	"strings"
	"testing"

	"github.com/kdeconinck/spot/internal/dsl/ast"
	"github.com/kdeconinck/spot/internal/dsl/parser"
	"github.com/kdeconinck/spot/internal/position"
	"github.com/kdeconinck/spot/internal/qa/claim"
)

const (
	HasErrorLabel = "Has error?"
	ErrorLabel    = "Error"
	ASTLabel      = "AST"
)

func Test_Parse(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		textInput  string
		wantOutput ast.Scope
	}{
		{
			name:       "When parsing an empty scope, the returned tree is correct.",
			textInput:  "scope {}",
			wantOutput: ast.Scope{},
		},
		{
			name:      "When parsing a scope with an include statement, the returned tree is correct.",
			textInput: `scope { include "**/*.txt" }`,
			wantOutput: ast.Scope{
				Includes: []ast.Include{
					{
						Pattern: ast.StringLiteral{},
					},
				},
			},
		},
		{
			name:      "When parsing a scope with an exclude statement, the returned tree is correct.",
			textInput: `scope { exclude "**/vendor/**" }`,
			wantOutput: ast.Scope{
				Excludes: []ast.Exclude{
					{
						Pattern: ast.StringLiteral{},
					},
				},
			},
		},
		{
			name: "When parsing a scope with include and exclude statements, the returned tree is correct.",
			textInput: `scope {
                include "**/*.txt"
                exclude "**/vendor/**"
        }`,
			wantOutput: ast.Scope{
				Includes: []ast.Include{
					{
						Pattern: ast.StringLiteral{},
					},
				},
				Excludes: []ast.Exclude{
					{
						Pattern: ast.StringLiteral{},
					},
				},
			},
		},
		{
			name: "When parsing an empty charset block, the returned tree is correct.",
			textInput: `scope {
                charset {
                }
        }`,
			wantOutput: ast.Scope{
				Charset: &ast.Charset{},
			},
		},
		{
			name: "When parsing a charset block with a single member containing a character literal, the returned tree is correct.",
			textInput: `scope {
                charset {
                        letter = 'a'
                }
        }`,
			wantOutput: ast.Scope{
				Charset: &ast.Charset{
					Members: []ast.CharsetMember{
						{
							Value: ast.CharacterLiteral{},
						},
					},
				},
			},
		},
		{
			name: "When parsing an empty lexer block, the returned tree is correct.",
			textInput: `scope {
                lexer {
                }
        }`,
			wantOutput: ast.Scope{
				Lexer: &ast.Lexer{},
			},
		},
		{
			name: "When parsing a lexer block with a single rule containing a reference, the returned tree is correct.",
			textInput: `scope {
                lexer {
                        identifier = letter
                }
        }`,
			wantOutput: ast.Scope{
				Lexer: &ast.Lexer{
					Rules: []ast.LexerRule{
						{
							Value: ast.Reference{},
						},
					},
				},
			},
		},

		{
			name: "When parsing a charset member with a parenthesized character literal, the returned tree is correct.",
			textInput: `scope {
                charset {
                        letter = ('a')
                }
        }`,
			wantOutput: ast.Scope{
				Charset: &ast.Charset{
					Members: []ast.CharsetMember{
						{
							Value: ast.CharacterLiteral{},
						},
					},
				},
			},
		},
		{
			name: "When parsing a lexer rule with a parenthesized reference, the returned tree is correct.",
			textInput: `scope {
                lexer {
                        identifier = (letter)
                }
        }`,
			wantOutput: ast.Scope{
				Lexer: &ast.Lexer{
					Rules: []ast.LexerRule{
						{
							Value: ast.Reference{},
						},
					},
				},
			},
		},
		{
			name: "When parsing a lexer rule with a string literal, the returned tree is correct.",
			textInput: `scope {
                lexer {
                        public_keyword = "public"
                }
        }`,
			wantOutput: ast.Scope{
				Lexer: &ast.Lexer{
					Rules: []ast.LexerRule{
						{
							Value: ast.StringLiteral{},
						},
					},
				},
			},
		},
		{
			name: "When parsing a lexer rule with a parenthesized string literal, the returned tree is correct.",
			textInput: `scope {
                lexer {
                        public_keyword = ("public")
                }
        }`,
			wantOutput: ast.Scope{
				Lexer: &ast.Lexer{
					Rules: []ast.LexerRule{
						{
							Value: ast.StringLiteral{},
						},
					},
				},
			},
		},
		{
			name: "When parsing a lexer rule with a character literal, the returned tree is correct.",
			textInput: `scope {
                lexer {
                        letter = 'a'
                }
        }`,
			wantOutput: ast.Scope{
				Lexer: &ast.Lexer{
					Rules: []ast.LexerRule{
						{
							Value: ast.CharacterLiteral{},
						},
					},
				},
			},
		},
		{
			name: "When parsing a lexer rule with a parenthesized character literal, the returned tree is correct.",
			textInput: `scope {
                lexer {
                        letter = ('a')
                }
        }`,
			wantOutput: ast.Scope{
				Lexer: &ast.Lexer{
					Rules: []ast.LexerRule{
						{
							Value: ast.CharacterLiteral{},
						},
					},
				},
			},
		},

		{
			name: "When parsing a charset member with a character range, the returned tree is correct.",
			textInput: `scope {
                charset {
                        letter = 'a'..'z'
                }
        }`,
			wantOutput: ast.Scope{
				Charset: &ast.Charset{
					Members: []ast.CharsetMember{
						{
							Value: ast.CharacterRange{
								Start: ast.CharacterLiteral{},
								End:   ast.CharacterLiteral{},
							},
						},
					},
				},
			},
		},
		{
			name: "When parsing a lexer rule with a character range, the returned tree is correct.",
			textInput: `scope {
                lexer {
                        identifier = 'a'..'z'
                }
        }`,
			wantOutput: ast.Scope{
				Lexer: &ast.Lexer{
					Rules: []ast.LexerRule{
						{
							Value: ast.CharacterRange{
								Start: ast.CharacterLiteral{},
								End:   ast.CharacterLiteral{},
							},
						},
					},
				},
			},
		},
		{
			name: "When parsing a parenthesized charset character range, the returned tree is correct.",
			textInput: `scope {
                charset {
                        letter = ('a'..'z')
                }
        }`,
			wantOutput: ast.Scope{
				Charset: &ast.Charset{
					Members: []ast.CharsetMember{
						{
							Value: ast.CharacterRange{
								Start: ast.CharacterLiteral{},
								End:   ast.CharacterLiteral{},
							},
						},
					},
				},
			},
		},
		{
			name: "When parsing a parenthesized lexer character range, the returned tree is correct.",
			textInput: `scope {
                lexer {
                        identifier = ('a'..'z')
                }
        }`,
			wantOutput: ast.Scope{
				Lexer: &ast.Lexer{
					Rules: []ast.LexerRule{
						{
							Value: ast.CharacterRange{
								Start: ast.CharacterLiteral{},
								End:   ast.CharacterLiteral{},
							},
						},
					},
				},
			},
		},

		{
			name: "When parsing a charset star repetition, the returned tree is correct.",
			textInput: `scope {
                charset {
                        letters = 'a'*
                }
        }`,
			wantOutput: ast.Scope{
				Charset: &ast.Charset{
					Members: []ast.CharsetMember{
						{
							Value: ast.Repetition{
								Expression: ast.CharacterLiteral{},
								Minimum:    0,
							},
						},
					},
				},
			},
		},
		{
			name: "When parsing a charset plus repetition, the returned tree is correct.",
			textInput: `scope {
                charset {
                        letters = 'a'+
                }
        }`,
			wantOutput: ast.Scope{
				Charset: &ast.Charset{
					Members: []ast.CharsetMember{
						{
							Value: ast.Repetition{
								Expression: ast.CharacterLiteral{},
								Minimum:    1,
							},
						},
					},
				},
			},
		},
		{
			name: "When parsing a charset question repetition, the returned tree is correct.",
			textInput: `scope {
                charset {
                        letters = 'a'?
                }
        }`,
			wantOutput: ast.Scope{
				Charset: &ast.Charset{
					Members: []ast.CharsetMember{
						{
							Value: ast.Repetition{
								Expression: ast.CharacterLiteral{},
								Minimum:    0,
								Maximum:    intPointer(1),
							},
						},
					},
				},
			},
		},
		{
			name: "When parsing a charset exact repetition, the returned tree is correct.",
			textInput: `scope {
                charset {
                        letters = 'a'{3}
                }
        }`,
			wantOutput: ast.Scope{
				Charset: &ast.Charset{
					Members: []ast.CharsetMember{
						{
							Value: ast.Repetition{
								Expression: ast.CharacterLiteral{},
								Minimum:    3,
								Maximum:    intPointer(3),
							},
						},
					},
				},
			},
		},
		{
			name: "When parsing a charset lower-bounded repetition, the returned tree is correct.",
			textInput: `scope {
                charset {
                        letters = 'a'{2,}
                }
        }`,
			wantOutput: ast.Scope{
				Charset: &ast.Charset{
					Members: []ast.CharsetMember{
						{
							Value: ast.Repetition{
								Expression: ast.CharacterLiteral{},
								Minimum:    2,
							},
						},
					},
				},
			},
		},
		{
			name: "When parsing a charset bounded repetition, the returned tree is correct.",
			textInput: `scope {
                charset {
                        letters = 'a'{1,8}
                }
        }`,
			wantOutput: ast.Scope{
				Charset: &ast.Charset{
					Members: []ast.CharsetMember{
						{
							Value: ast.Repetition{
								Expression: ast.CharacterLiteral{},
								Minimum:    1,
								Maximum:    intPointer(8),
							},
						},
					},
				},
			},
		},
		{
			name: "When parsing a parenthesized charset repetition, the returned tree is correct.",
			textInput: `scope {
                charset {
                        letters = ('a'+)
                }
        }`,
			wantOutput: ast.Scope{
				Charset: &ast.Charset{
					Members: []ast.CharsetMember{
						{
							Value: ast.Repetition{
								Expression: ast.CharacterLiteral{},
								Minimum:    1,
							},
						},
					},
				},
			},
		},

		{
			name: "When parsing a lexer star repetition, the returned tree is correct.",
			textInput: `scope {
                lexer {
                        identifier = letter*
                }
        }`,
			wantOutput: ast.Scope{
				Lexer: &ast.Lexer{
					Rules: []ast.LexerRule{
						{
							Value: ast.Repetition{
								Expression: ast.Reference{},
								Minimum:    0,
							},
						},
					},
				},
			},
		},
		{
			name: "When parsing a lexer plus repetition, the returned tree is correct.",
			textInput: `scope {
                lexer {
                        identifier = letter+
                }
        }`,
			wantOutput: ast.Scope{
				Lexer: &ast.Lexer{
					Rules: []ast.LexerRule{
						{
							Value: ast.Repetition{
								Expression: ast.Reference{},
								Minimum:    1,
							},
						},
					},
				},
			},
		},
		{
			name: "When parsing a lexer question repetition, the returned tree is correct.",
			textInput: `scope {
                lexer {
                        identifier = letter?
                }
        }`,
			wantOutput: ast.Scope{
				Lexer: &ast.Lexer{
					Rules: []ast.LexerRule{
						{
							Value: ast.Repetition{
								Expression: ast.Reference{},
								Minimum:    0,
								Maximum:    intPointer(1),
							},
						},
					},
				},
			},
		},
		{
			name: "When parsing a lexer exact repetition, the returned tree is correct.",
			textInput: `scope {
                lexer {
                        identifier = letter{3}
                }
        }`,
			wantOutput: ast.Scope{
				Lexer: &ast.Lexer{
					Rules: []ast.LexerRule{
						{
							Value: ast.Repetition{
								Expression: ast.Reference{},
								Minimum:    3,
								Maximum:    intPointer(3),
							},
						},
					},
				},
			},
		},
		{
			name: "When parsing a lexer lower-bounded repetition, the returned tree is correct.",
			textInput: `scope {
                lexer {
                        identifier = letter{2,}
                }
        }`,
			wantOutput: ast.Scope{
				Lexer: &ast.Lexer{
					Rules: []ast.LexerRule{
						{
							Value: ast.Repetition{
								Expression: ast.Reference{},
								Minimum:    2,
							},
						},
					},
				},
			},
		},
		{
			name: "When parsing a lexer bounded repetition, the returned tree is correct.",
			textInput: `scope {
                lexer {
                        identifier = letter{1,8}
                }
        }`,
			wantOutput: ast.Scope{
				Lexer: &ast.Lexer{
					Rules: []ast.LexerRule{
						{
							Value: ast.Repetition{
								Expression: ast.Reference{},
								Minimum:    1,
								Maximum:    intPointer(8),
							},
						},
					},
				},
			},
		},
		{
			name: "When parsing a parenthesized lexer repetition, the returned tree is correct.",
			textInput: `scope {
                lexer {
                        identifier = (letter+)
                }
        }`,
			wantOutput: ast.Scope{
				Lexer: &ast.Lexer{
					Rules: []ast.LexerRule{
						{
							Value: ast.Repetition{
								Expression: ast.Reference{},
								Minimum:    1,
							},
						},
					},
				},
			},
		},

		{
			name: "When parsing a charset concatenation, the returned tree is correct.",
			textInput: `scope {
                charset {
                        letter = 'a' 'b'
                }
        }`,
			wantOutput: ast.Scope{
				Charset: &ast.Charset{
					Members: []ast.CharsetMember{
						{
							Value: ast.Concatenation{
								Expressions: []ast.Expression{
									ast.CharacterLiteral{},
									ast.CharacterLiteral{},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "When parsing a parenthesized charset concatenation, the returned tree is correct.",
			textInput: `scope {
                charset {
                        letter = ('a' 'b')
                }
        }`,
			wantOutput: ast.Scope{
				Charset: &ast.Charset{
					Members: []ast.CharsetMember{
						{
							Value: ast.Concatenation{
								Expressions: []ast.Expression{
									ast.CharacterLiteral{},
									ast.CharacterLiteral{},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "When parsing a lexer concatenation, the returned tree is correct.",
			textInput: `scope {
                lexer {
                        identifier = ident_start ident_part
                }
        }`,
			wantOutput: ast.Scope{
				Lexer: &ast.Lexer{
					Rules: []ast.LexerRule{
						{
							Value: ast.Concatenation{
								Expressions: []ast.Expression{
									ast.Reference{},
									ast.Reference{},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "When parsing a parenthesized lexer concatenation, the returned tree is correct.",
			textInput: `scope {
                lexer {
                        identifier = (ident_start ident_part)
                }
        }`,
			wantOutput: ast.Scope{
				Lexer: &ast.Lexer{
					Rules: []ast.LexerRule{
						{
							Value: ast.Concatenation{
								Expressions: []ast.Expression{
									ast.Reference{},
									ast.Reference{},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "When parsing a lexer concatenation with a string literal, the returned tree is correct.",
			textInput: `scope {
                lexer {
                        public_suffix = "public" letter
                }
        }`,
			wantOutput: ast.Scope{
				Lexer: &ast.Lexer{
					Rules: []ast.LexerRule{
						{
							Value: ast.Concatenation{
								Expressions: []ast.Expression{
									ast.StringLiteral{},
									ast.Reference{},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "When parsing a parenthesized lexer concatenation with a string literal, the returned tree is correct.",
			textInput: `scope {
                lexer {
                        public_suffix = ("public" letter)
                }
        }`,
			wantOutput: ast.Scope{
				Lexer: &ast.Lexer{
					Rules: []ast.LexerRule{
						{
							Value: ast.Concatenation{
								Expressions: []ast.Expression{
									ast.StringLiteral{},
									ast.Reference{},
								},
							},
						},
					},
				},
			},
		},

		{
			name: "When parsing a charset alternation, the returned tree is correct.",
			textInput: `scope {
                charset {
                        letter = 'a' | 'b'
                }
        }`,
			wantOutput: ast.Scope{
				Charset: &ast.Charset{
					Members: []ast.CharsetMember{
						{
							Value: ast.Alternation{
								Expressions: []ast.Expression{
									ast.CharacterLiteral{},
									ast.CharacterLiteral{},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "When parsing a parenthesized charset alternation, the returned tree is correct.",
			textInput: `scope {
                charset {
                        letter = ('a' | 'b')
                }
        }`,
			wantOutput: ast.Scope{
				Charset: &ast.Charset{
					Members: []ast.CharsetMember{
						{
							Value: ast.Alternation{
								Expressions: []ast.Expression{
									ast.CharacterLiteral{},
									ast.CharacterLiteral{},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "When parsing a lexer alternation, the returned tree is correct.",
			textInput: `scope {
                lexer {
                        identifier = letter | digit
                }
        }`,
			wantOutput: ast.Scope{
				Lexer: &ast.Lexer{
					Rules: []ast.LexerRule{
						{
							Value: ast.Alternation{
								Expressions: []ast.Expression{
									ast.Reference{},
									ast.Reference{},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "When parsing a parenthesized lexer alternation, the returned tree is correct.",
			textInput: `scope {
                lexer {
                        identifier = (letter | digit)
                }
        }`,
			wantOutput: ast.Scope{
				Lexer: &ast.Lexer{
					Rules: []ast.LexerRule{
						{
							Value: ast.Alternation{
								Expressions: []ast.Expression{
									ast.Reference{},
									ast.Reference{},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "When parsing a lexer alternation with strings, the returned tree is correct.",
			textInput: `scope {
                lexer {
                        visibility = "public" | "private"
                }
        }`,
			wantOutput: ast.Scope{
				Lexer: &ast.Lexer{
					Rules: []ast.LexerRule{
						{
							Value: ast.Alternation{
								Expressions: []ast.Expression{
									ast.StringLiteral{},
									ast.StringLiteral{},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "When parsing a parenthesized lexer alternation with strings, the returned tree is correct.",
			textInput: `scope {
                lexer {
                        visibility = ("public" | "private")
                }
        }`,
			wantOutput: ast.Scope{
				Lexer: &ast.Lexer{
					Rules: []ast.LexerRule{
						{
							Value: ast.Alternation{
								Expressions: []ast.Expression{
									ast.StringLiteral{},
									ast.StringLiteral{},
								},
							},
						},
					},
				},
			},
		},

		{
			name: "When parsing a scope with include, exclude, charset, and lexer blocks, the returned tree is correct.",
			textInput: `scope {
                include "**/*.txt"
                exclude "**/vendor/**"
                charset {
                        letter = 'a'..'z' | 'A'..'Z'
                        space = ' ' | '\t'
                }
                lexer {
                        identifier = letter+
                        whitespace = space+
                        public_keyword = "public"
                }
        }`,
			wantOutput: ast.Scope{
				Includes: []ast.Include{
					{
						Pattern: ast.StringLiteral{},
					},
				},
				Excludes: []ast.Exclude{
					{
						Pattern: ast.StringLiteral{},
					},
				},
				Charset: &ast.Charset{
					Members: []ast.CharsetMember{
						{
							Value: ast.Alternation{
								Expressions: []ast.Expression{
									ast.CharacterRange{
										Start: ast.CharacterLiteral{},
										End:   ast.CharacterLiteral{},
									},
									ast.CharacterRange{
										Start: ast.CharacterLiteral{},
										End:   ast.CharacterLiteral{},
									},
								},
							},
						},
						{
							Value: ast.Alternation{
								Expressions: []ast.Expression{
									ast.CharacterLiteral{},
									ast.CharacterLiteral{},
								},
							},
						},
					},
				},
				Lexer: &ast.Lexer{
					Rules: []ast.LexerRule{
						{
							Value: ast.Repetition{
								Expression: ast.Reference{},
								Minimum:    1,
							},
						},
						{
							Value: ast.Repetition{
								Expression: ast.Reference{},
								Minimum:    1,
							},
						},
						{
							Value: ast.StringLiteral{},
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Act.
			gotOutput, err := parser.Parse(testCase.textInput)

			// Assert.
			claim.Equal(t, testCase.name, false, err != nil, HasErrorLabel)
			claim.DeepEqual(t, testCase.name, testCase.wantOutput, withoutRanges(gotOutput), ASTLabel)
		})
	}
}

func Test_Parse_Range(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name       string
		textInput  string
		wantOutput ast.Scope
	}{
		{
			name:      "When parsing scope statements and blocks, the returned ranges are correct.",
			textInput: "scope {\n  include \"**/*.txt\"\n  exclude \"**/vendor/**\"\n  charset {\n    letter = 'a'\n  }\n  lexer {\n    identifier = letter\n  }\n}",
			wantOutput: ast.Scope{
				Range: rangeOf("scope {\n  include \"**/*.txt\"\n  exclude \"**/vendor/**\"\n  charset {\n    letter = 'a'\n  }\n  lexer {\n    identifier = letter\n  }\n}", "scope {\n  include \"**/*.txt\"\n  exclude \"**/vendor/**\"\n  charset {\n    letter = 'a'\n  }\n  lexer {\n    identifier = letter\n  }\n}"),
				Includes: []ast.Include{
					{
						Range:   rangeOf("scope {\n  include \"**/*.txt\"\n  exclude \"**/vendor/**\"\n  charset {\n    letter = 'a'\n  }\n  lexer {\n    identifier = letter\n  }\n}", "include \"**/*.txt\""),
						Pattern: ast.StringLiteral{Range: rangeOf("scope {\n  include \"**/*.txt\"\n  exclude \"**/vendor/**\"\n  charset {\n    letter = 'a'\n  }\n  lexer {\n    identifier = letter\n  }\n}", "\"**/*.txt\"")},
					},
				},
				Excludes: []ast.Exclude{
					{
						Range:   rangeOf("scope {\n  include \"**/*.txt\"\n  exclude \"**/vendor/**\"\n  charset {\n    letter = 'a'\n  }\n  lexer {\n    identifier = letter\n  }\n}", "exclude \"**/vendor/**\""),
						Pattern: ast.StringLiteral{Range: rangeOf("scope {\n  include \"**/*.txt\"\n  exclude \"**/vendor/**\"\n  charset {\n    letter = 'a'\n  }\n  lexer {\n    identifier = letter\n  }\n}", "\"**/vendor/**\"")},
					},
				},
				Charset: &ast.Charset{
					Range: rangeOf("scope {\n  include \"**/*.txt\"\n  exclude \"**/vendor/**\"\n  charset {\n    letter = 'a'\n  }\n  lexer {\n    identifier = letter\n  }\n}", "charset {\n    letter = 'a'\n  }"),
					Members: []ast.CharsetMember{
						{
							Range: rangeOf("scope {\n  include \"**/*.txt\"\n  exclude \"**/vendor/**\"\n  charset {\n    letter = 'a'\n  }\n  lexer {\n    identifier = letter\n  }\n}", "letter = 'a'"),
							Value: ast.CharacterLiteral{Range: rangeOf("scope {\n  include \"**/*.txt\"\n  exclude \"**/vendor/**\"\n  charset {\n    letter = 'a'\n  }\n  lexer {\n    identifier = letter\n  }\n}", "'a'")},
						},
					},
				},
				Lexer: &ast.Lexer{
					Range: rangeOf("scope {\n  include \"**/*.txt\"\n  exclude \"**/vendor/**\"\n  charset {\n    letter = 'a'\n  }\n  lexer {\n    identifier = letter\n  }\n}", "lexer {\n    identifier = letter\n  }"),
					Rules: []ast.LexerRule{
						{
							Range: rangeOf("scope {\n  include \"**/*.txt\"\n  exclude \"**/vendor/**\"\n  charset {\n    letter = 'a'\n  }\n  lexer {\n    identifier = letter\n  }\n}", "identifier = letter"),
							Value: ast.Reference{Range: secondRangeOf("scope {\n  include \"**/*.txt\"\n  exclude \"**/vendor/**\"\n  charset {\n    letter = 'a'\n  }\n  lexer {\n    identifier = letter\n  }\n}", "letter")},
						},
					},
				},
			},
		},
		{
			name:      "When parsing composite expressions, the returned ranges are correct.",
			textInput: "scope {\n  charset {\n    letter = ('a'..'z' | 'A'..'Z')+\n  }\n  lexer {\n    identifier = (letter digit)\n  }\n}",
			wantOutput: ast.Scope{
				Range: rangeOf("scope {\n  charset {\n    letter = ('a'..'z' | 'A'..'Z')+\n  }\n  lexer {\n    identifier = (letter digit)\n  }\n}", "scope {\n  charset {\n    letter = ('a'..'z' | 'A'..'Z')+\n  }\n  lexer {\n    identifier = (letter digit)\n  }\n}"),
				Charset: &ast.Charset{
					Range: rangeOf("scope {\n  charset {\n    letter = ('a'..'z' | 'A'..'Z')+\n  }\n  lexer {\n    identifier = (letter digit)\n  }\n}", "charset {\n    letter = ('a'..'z' | 'A'..'Z')+\n  }"),
					Members: []ast.CharsetMember{
						{
							Range: rangeOf("scope {\n  charset {\n    letter = ('a'..'z' | 'A'..'Z')+\n  }\n  lexer {\n    identifier = (letter digit)\n  }\n}", "letter = ('a'..'z' | 'A'..'Z')+"),
							Value: ast.Repetition{
								Range: rangeOf("scope {\n  charset {\n    letter = ('a'..'z' | 'A'..'Z')+\n  }\n  lexer {\n    identifier = (letter digit)\n  }\n}", "('a'..'z' | 'A'..'Z')+"),
								Expression: ast.Alternation{
									Range: rangeOf("scope {\n  charset {\n    letter = ('a'..'z' | 'A'..'Z')+\n  }\n  lexer {\n    identifier = (letter digit)\n  }\n}", "('a'..'z' | 'A'..'Z')"),
									Expressions: []ast.Expression{
										ast.CharacterRange{
											Range: rangeOf("scope {\n  charset {\n    letter = ('a'..'z' | 'A'..'Z')+\n  }\n  lexer {\n    identifier = (letter digit)\n  }\n}", "'a'..'z'"),
											Start: ast.CharacterLiteral{Range: rangeOf("scope {\n  charset {\n    letter = ('a'..'z' | 'A'..'Z')+\n  }\n  lexer {\n    identifier = (letter digit)\n  }\n}", "'a'")},
											End:   ast.CharacterLiteral{Range: rangeOf("scope {\n  charset {\n    letter = ('a'..'z' | 'A'..'Z')+\n  }\n  lexer {\n    identifier = (letter digit)\n  }\n}", "'z'")},
										},
										ast.CharacterRange{
											Range: rangeOf("scope {\n  charset {\n    letter = ('a'..'z' | 'A'..'Z')+\n  }\n  lexer {\n    identifier = (letter digit)\n  }\n}", "'A'..'Z'"),
											Start: ast.CharacterLiteral{Range: rangeOf("scope {\n  charset {\n    letter = ('a'..'z' | 'A'..'Z')+\n  }\n  lexer {\n    identifier = (letter digit)\n  }\n}", "'A'")},
											End:   ast.CharacterLiteral{Range: rangeOf("scope {\n  charset {\n    letter = ('a'..'z' | 'A'..'Z')+\n  }\n  lexer {\n    identifier = (letter digit)\n  }\n}", "'Z'")},
										},
									},
								},
								Minimum: 1,
							},
						},
					},
				},
				Lexer: &ast.Lexer{
					Range: rangeOf("scope {\n  charset {\n    letter = ('a'..'z' | 'A'..'Z')+\n  }\n  lexer {\n    identifier = (letter digit)\n  }\n}", "lexer {\n    identifier = (letter digit)\n  }"),
					Rules: []ast.LexerRule{
						{
							Range: rangeOf("scope {\n  charset {\n    letter = ('a'..'z' | 'A'..'Z')+\n  }\n  lexer {\n    identifier = (letter digit)\n  }\n}", "identifier = (letter digit)"),
							Value: ast.Concatenation{
								Range: rangeOf("scope {\n  charset {\n    letter = ('a'..'z' | 'A'..'Z')+\n  }\n  lexer {\n    identifier = (letter digit)\n  }\n}", "(letter digit)"),
								Expressions: []ast.Expression{
									ast.Reference{Range: secondRangeOf("scope {\n  charset {\n    letter = ('a'..'z' | 'A'..'Z')+\n  }\n  lexer {\n    identifier = (letter digit)\n  }\n}", "letter")},
									ast.Reference{Range: rangeOf("scope {\n  charset {\n    letter = ('a'..'z' | 'A'..'Z')+\n  }\n  lexer {\n    identifier = (letter digit)\n  }\n}", "digit")},
								},
							},
						},
					},
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Act.
			gotOutput, err := parser.Parse(testCase.textInput)

			// Assert.
			claim.Equal(t, testCase.name, false, err != nil, HasErrorLabel)
			claim.DeepEqual(t, testCase.name, testCase.wantOutput, gotOutput, ASTLabel)
		})
	}
}

func Test_Parse_Error(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name            string
		textInput       string
		wantErrorOutput string
	}{
		{
			name:            "When parsing input that does not start with scope, the returned error is correct.",
			textInput:       "include \"a\"",
			wantErrorOutput: "DSL Parser: Expected scope, got \"include\".",
		},
		{
			name:            "When parsing input with a scanner error, the returned error is correct.",
			textInput:       "@",
			wantErrorOutput: "DSL Scanner: Unexpected byte '@'.",
		},
		{
			name:            "When parsing a missing scope block, the returned error is correct.",
			textInput:       "scope",
			wantErrorOutput: "DSL Parser: Expected left-brace, got \"end-of-file\".",
		},
		{
			name:            "When parsing a scope member with an unexpected token, the returned error is correct.",
			textInput:       "scope { identifier }",
			wantErrorOutput: "DSL Parser: Expected scope member, got \"identifier\".",
		},
		{
			name:            "When parsing a scope block without a closing brace, the returned error is correct.",
			textInput:       "scope { include \"a\"",
			wantErrorOutput: "DSL Parser: Expected scope member, got \"end-of-file\".",
		},
		{
			name:            "When parsing extra tokens after the scope block, the returned error is correct.",
			textInput:       "scope {} scope",
			wantErrorOutput: "DSL Parser: Expected end-of-file, got \"scope\".",
		},

		{
			name:            "When parsing an include statement without a string literal, the returned error is correct.",
			textInput:       "scope { include letter }",
			wantErrorOutput: "DSL Parser: Expected string, got \"identifier\".",
		},

		{
			name:            "When parsing an exclude statement without a string literal, the returned error is correct.",
			textInput:       "scope { exclude letter }",
			wantErrorOutput: "DSL Parser: Expected string, got \"identifier\".",
		},

		{
			name:            "When parsing a charset block without a left brace, the returned error is correct.",
			textInput:       "scope { charset }",
			wantErrorOutput: "DSL Parser: Expected left-brace, got \"right-brace\".",
		},
		{
			name:            "When parsing a charset member without a name, the returned error is correct.",
			textInput:       "scope { charset { = 'a' } }",
			wantErrorOutput: "DSL Parser: Expected identifier, got \"equal\".",
		},
		{
			name:            "When parsing a charset member without an equal sign, the returned error is correct.",
			textInput:       "scope { charset { letter 'a' } }",
			wantErrorOutput: "DSL Parser: Expected equal, got \"character\".",
		},
		{
			name:            "When parsing a charset member without a right-hand expression, the returned error is correct.",
			textInput:       "scope { charset { letter = } }",
			wantErrorOutput: "DSL Parser: Expected expression, got \"right-brace\".",
		},
		{
			name:            "When parsing a charset member without a right-hand expression at the end of file, the returned error is correct.",
			textInput:       "scope { charset { letter =",
			wantErrorOutput: "DSL Parser: Expected expression, got \"end-of-file\".",
		},
		{
			name:            "When parsing a charset block without a closing brace, the returned error is correct.",
			textInput:       "scope { charset { letter = 'a' ",
			wantErrorOutput: "DSL Parser: Expected identifier, got \"end-of-file\".",
		},

		{
			name:            "When parsing a lexer block without a left brace, the returned error is correct.",
			textInput:       "scope { lexer }",
			wantErrorOutput: "DSL Parser: Expected left-brace, got \"right-brace\".",
		},
		{
			name:            "When parsing a lexer rule without a name, the returned error is correct.",
			textInput:       "scope { lexer { = letter } }",
			wantErrorOutput: "DSL Parser: Expected identifier, got \"equal\".",
		},
		{
			name:            "When parsing a lexer rule without an equal sign, the returned error is correct.",
			textInput:       "scope { lexer { identifier letter } }",
			wantErrorOutput: "DSL Parser: Expected equal, got \"identifier\".",
		},
		{
			name:            "When parsing a lexer rule without a right-hand expression, the returned error is correct.",
			textInput:       "scope { lexer { identifier = } }",
			wantErrorOutput: "DSL Parser: Expected expression, got \"right-brace\".",
		},
		{
			name:            "When parsing a lexer rule without a right-hand expression at the end of file, the returned error is correct.",
			textInput:       "scope { lexer { identifier =",
			wantErrorOutput: "DSL Parser: Expected expression, got \"end-of-file\".",
		},
		{
			name:            "When parsing a lexer block without a closing brace, the returned error is correct.",
			textInput:       "scope { lexer { identifier = letter ",
			wantErrorOutput: "DSL Parser: Expected identifier, got \"end-of-file\".",
		},

		{
			name:            "When parsing a charset expression with a string literal, the returned error is correct.",
			textInput:       "scope { charset { letter = \"a\" } }",
			wantErrorOutput: "DSL Parser: Expected charset expression, got \"string\".",
		},

		{
			name:            "When parsing a charset alternation without a right-hand expression, the returned error is correct.",
			textInput:       "scope { charset { letter = 'a' | } }",
			wantErrorOutput: "DSL Parser: Expected expression, got \"right-brace\".",
		},
		{
			name:            "When parsing a lexer alternation without a right-hand expression, the returned error is correct.",
			textInput:       "scope { lexer { identifier = letter | } }",
			wantErrorOutput: "DSL Parser: Expected expression, got \"right-brace\".",
		},

		{
			name:            "When parsing an empty charset grouped expression, the returned error is correct.",
			textInput:       "scope { charset { letter = () } }",
			wantErrorOutput: "DSL Parser: Expected expression, got \"right-paren\".",
		},
		{
			name:            "When parsing a charset grouped expression whose contents are missing, the returned error is correct.",
			textInput:       "scope { charset { letter = ( } }",
			wantErrorOutput: "DSL Parser: Expected expression, got \"right-brace\".",
		},
		{
			name:            "When parsing a charset grouped expression without a closing paren, the returned error is correct.",
			textInput:       "scope { charset { letter = ('a' } }",
			wantErrorOutput: "DSL Parser: Expected right-paren, got \"right-brace\".",
		},

		{
			name:            "When parsing an empty lexer grouped expression, the returned error is correct.",
			textInput:       "scope { lexer { identifier = () } }",
			wantErrorOutput: "DSL Parser: Expected expression, got \"right-paren\".",
		},
		{
			name:            "When parsing a lexer grouped expression whose contents are missing, the returned error is correct.",
			textInput:       "scope { lexer { identifier = ( } }",
			wantErrorOutput: "DSL Parser: Expected expression, got \"right-brace\".",
		},
		{
			name:            "When parsing a lexer grouped expression without a closing paren, the returned error is correct.",
			textInput:       "scope { lexer { identifier = (letter } }",
			wantErrorOutput: "DSL Parser: Expected right-paren, got \"right-brace\".",
		},

		{
			name:            "When parsing a charset range whose left-hand side is not a character literal, the returned error is correct.",
			textInput:       "scope { charset { letter = alpha..'z' } }",
			wantErrorOutput: "DSL Parser: Expected character literal before \"dot-dot\".",
		},
		{
			name:            "When parsing a charset range whose right-hand side is not a character literal, the returned error is correct.",
			textInput:       "scope { charset { letter = 'a'..alpha } }",
			wantErrorOutput: "DSL Parser: Expected character literal after \"dot-dot\".",
		},
		{
			name:            "When parsing a charset range whose right-hand side is a string literal, the returned error is correct.",
			textInput:       "scope { charset { letter = 'a'..\"a\" } }",
			wantErrorOutput: "DSL Parser: Expected charset expression, got \"string\".",
		},
		{
			name:            "When parsing a charset range without a right-hand expression, the returned error is correct.",
			textInput:       "scope { charset { letter = 'a'.. } }",
			wantErrorOutput: "DSL Parser: Expected expression, got \"right-brace\".",
		},
		{
			name:            "When parsing a lexer range whose left-hand side is not a character literal, the returned error is correct.",
			textInput:       "scope { lexer { identifier = letter..'z' } }",
			wantErrorOutput: "DSL Parser: Expected character literal before \"dot-dot\".",
		},
		{
			name:            "When parsing a lexer range whose right-hand side is not a character literal, the returned error is correct.",
			textInput:       "scope { lexer { identifier = 'a'..letter } }",
			wantErrorOutput: "DSL Parser: Expected character literal after \"dot-dot\".",
		},
		{
			name:            "When parsing a lexer range without a right-hand expression, the returned error is correct.",
			textInput:       "scope { lexer { identifier = 'a'.. } }",
			wantErrorOutput: "DSL Parser: Expected expression, got \"right-brace\".",
		},

		{
			name:            "When parsing a charset bounded repetition without a minimum number, the returned error is correct.",
			textInput:       "scope { charset { letter = 'a'{} } }",
			wantErrorOutput: "DSL Parser: Expected number, got \"right-brace\".",
		},
		{
			name:            "When parsing a charset bounded repetition without a minimum number at the end of file, the returned error is correct.",
			textInput:       "scope { charset { letter = 'a'{",
			wantErrorOutput: "DSL Parser: Expected number, got \"end-of-file\".",
		},
		{
			name:            "When parsing a charset bounded repetition without a comma or closing brace, the returned error is correct.",
			textInput:       "scope { charset { letter = 'a'{1?} } }",
			wantErrorOutput: "DSL Parser: Expected comma, got \"question\".",
		},
		{
			name:            "When parsing a charset bounded repetition with an identifier maximum, the returned error is correct.",
			textInput:       "scope { charset { letter = 'a'{1,a} } }",
			wantErrorOutput: "DSL Parser: Expected number, got \"identifier\".",
		},
		{
			name:            "When parsing a charset bounded repetition without a maximum number at the end of file, the returned error is correct.",
			textInput:       "scope { charset { letter = 'a'{1,",
			wantErrorOutput: "DSL Parser: Expected number, got \"end-of-file\".",
		},
		{
			name:            "When parsing a charset bounded repetition without a closing brace after the maximum, the returned error is correct.",
			textInput:       "scope { charset { letter = 'a'{1,2) } }",
			wantErrorOutput: "DSL Parser: Expected right-brace, got \"right-paren\".",
		},
		{
			name:            "When parsing a charset repetition with an out-of-range minimum, the returned error is correct.",
			textInput:       "scope { charset { letter = 'a'{999999999999999999999999999999} } }",
			wantErrorOutput: "strconv.Atoi: parsing \"999999999999999999999999999999\": value out of range",
		},
		{
			name:            "When parsing a charset repetition with an out-of-range maximum, the returned error is correct.",
			textInput:       "scope { charset { letter = 'a'{1,999999999999999999999999999999} } }",
			wantErrorOutput: "strconv.Atoi: parsing \"999999999999999999999999999999\": value out of range",
		},

		{
			name:            "When parsing a lexer bounded repetition without a minimum number, the returned error is correct.",
			textInput:       "scope { lexer { identifier = letter{} } }",
			wantErrorOutput: "DSL Parser: Expected number, got \"right-brace\".",
		},
		{
			name:            "When parsing a lexer bounded repetition without a minimum number at the end of file, the returned error is correct.",
			textInput:       "scope { lexer { identifier = letter{",
			wantErrorOutput: "DSL Parser: Expected number, got \"end-of-file\".",
		},
		{
			name:            "When parsing a lexer bounded repetition without a comma or closing brace, the returned error is correct.",
			textInput:       "scope { lexer { identifier = letter{1?} } }",
			wantErrorOutput: "DSL Parser: Expected comma, got \"question\".",
		},
		{
			name:            "When parsing a lexer bounded repetition with an identifier maximum, the returned error is correct.",
			textInput:       "scope { lexer { identifier = letter{1,a} } }",
			wantErrorOutput: "DSL Parser: Expected number, got \"identifier\".",
		},
		{
			name:            "When parsing a lexer bounded repetition without a maximum number at the end of file, the returned error is correct.",
			textInput:       "scope { lexer { identifier = letter{1,",
			wantErrorOutput: "DSL Parser: Expected number, got \"end-of-file\".",
		},
		{
			name:            "When parsing a lexer bounded repetition without a closing brace after the maximum, the returned error is correct.",
			textInput:       "scope { lexer { identifier = letter{1,2) } }",
			wantErrorOutput: "DSL Parser: Expected right-brace, got \"right-paren\".",
		},
		{
			name:            "When parsing a lexer repetition with an out-of-range minimum, the returned error is correct.",
			textInput:       "scope { lexer { identifier = letter{999999999999999999999999999999} } }",
			wantErrorOutput: "strconv.Atoi: parsing \"999999999999999999999999999999\": value out of range",
		},
		{
			name:            "When parsing a lexer repetition with an out-of-range maximum, the returned error is correct.",
			textInput:       "scope { lexer { identifier = letter{1,999999999999999999999999999999} } }",
			wantErrorOutput: "strconv.Atoi: parsing \"999999999999999999999999999999\": value out of range",
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			// Act.
			_, err := parser.Parse(testCase.textInput)

			// Assert.
			claim.Equal(t, testCase.name, false, err == nil, HasErrorLabel)
			claim.Equal(t, testCase.name, testCase.wantErrorOutput, err.Error(), ErrorLabel)
		})
	}
}

func benchmark_Parse(b *testing.B, blockCountInput int) {
	b.Helper()

	text := benchmarkText(blockCountInput)

	for b.Loop() {
		_, err := parser.Parse(text)

		if err != nil {
			b.Fatalf("Parse() returned an unexpected error: %v.", err)
		}
	}
}

func Benchmark_Parse_1Block(b *testing.B)     { benchmark_Parse(b, 1) }
func Benchmark_Parse_10Blocks(b *testing.B)   { benchmark_Parse(b, 10) }
func Benchmark_Parse_100Blocks(b *testing.B)  { benchmark_Parse(b, 100) }
func Benchmark_Parse_1000Blocks(b *testing.B) { benchmark_Parse(b, 1000) }

func withoutRanges(scopeInput ast.Scope) ast.Scope {
	scopeInput.Range = position.Range{}

	for includeIndex := range scopeInput.Includes {
		scopeInput.Includes[includeIndex].Range = position.Range{}
		scopeInput.Includes[includeIndex].Pattern = withoutRangesStringLiteral(scopeInput.Includes[includeIndex].Pattern)
	}

	for excludeIndex := range scopeInput.Excludes {
		scopeInput.Excludes[excludeIndex].Range = position.Range{}
		scopeInput.Excludes[excludeIndex].Pattern = withoutRangesStringLiteral(scopeInput.Excludes[excludeIndex].Pattern)
	}

	if scopeInput.Charset != nil {
		scopeInput.Charset.Range = position.Range{}

		for memberIndex := range scopeInput.Charset.Members {
			scopeInput.Charset.Members[memberIndex].Range = position.Range{}
			scopeInput.Charset.Members[memberIndex].Value = withoutRangesExpression(scopeInput.Charset.Members[memberIndex].Value)
		}
	}

	if scopeInput.Lexer != nil {
		scopeInput.Lexer.Range = position.Range{}

		for ruleIndex := range scopeInput.Lexer.Rules {
			scopeInput.Lexer.Rules[ruleIndex].Range = position.Range{}
			scopeInput.Lexer.Rules[ruleIndex].Value = withoutRangesExpression(scopeInput.Lexer.Rules[ruleIndex].Value)
		}
	}

	return scopeInput
}

func withoutRangesExpression(expressionInput ast.Expression) ast.Expression {
	switch value := expressionInput.(type) {
	case ast.Reference:
		value.Range = position.Range{}

		return value

	case ast.StringLiteral:
		return withoutRangesStringLiteral(value)

	case ast.CharacterLiteral:
		value.Range = position.Range{}

		return value

	case ast.CharacterRange:
		value.Range = position.Range{}
		value.Start = withoutRangesExpression(value.Start).(ast.CharacterLiteral)
		value.End = withoutRangesExpression(value.End).(ast.CharacterLiteral)

		return value

	case ast.Alternation:
		value.Range = position.Range{}

		for expressionIndex := range value.Expressions {
			value.Expressions[expressionIndex] = withoutRangesExpression(value.Expressions[expressionIndex])
		}

		return value

	case ast.Concatenation:
		value.Range = position.Range{}

		for expressionIndex := range value.Expressions {
			value.Expressions[expressionIndex] = withoutRangesExpression(value.Expressions[expressionIndex])
		}

		return value

	case ast.Repetition:
		value.Range = position.Range{}
		value.Expression = withoutRangesExpression(value.Expression)

		return value

	default:
		return expressionInput
	}
}

func withoutRangesStringLiteral(literalInput ast.StringLiteral) ast.StringLiteral {
	literalInput.Range = position.Range{}

	return literalInput
}

func benchmarkText(blockCountInput int) string {
	var sb strings.Builder

	sb.WriteString("scope {\n")
	sb.WriteString("  include \"**/*.txt\"\n")
	sb.WriteString("  exclude \"**/vendor/**\"\n")
	sb.WriteString("  charset {\n")

	for range blockCountInput {
		sb.WriteString("    letters = ('a'..'z' | 'A'..'Z')+\n")
		sb.WriteString("    decimal_numbers = ('0'..'9'){1,8}\n")
		sb.WriteString("    space = ' ' | '\\t'?\n")
	}

	sb.WriteString("  }\n")
	sb.WriteString("  lexer {\n")

	for range blockCountInput {
		sb.WriteString("    identifier = letters+\n")
		sb.WriteString("    whitespace = space+\n")
		sb.WriteString("    public_keyword = \"public\"\n")
	}

	sb.WriteString("  }\n")
	sb.WriteString("}\n")

	return sb.String()
}

func intPointer(valueInput int) *int { return &valueInput }

func rangeOf(textInput, snippetInput string) position.Range {
	startOutput := strings.Index(textInput, snippetInput)

	return position.NewRange(
		position.NewPosition(startOutput),
		position.NewPosition(startOutput+len(snippetInput)))
}

func secondRangeOf(textInput, snippetInput string) position.Range {
	firstStartOutput := strings.Index(textInput, snippetInput)
	secondStartOutput := strings.Index(textInput[firstStartOutput+len(snippetInput):], snippetInput)
	startOutput := firstStartOutput + len(snippetInput) + secondStartOutput

	return position.NewRange(
		position.NewPosition(startOutput),
		position.NewPosition(startOutput+len(snippetInput)))
}
