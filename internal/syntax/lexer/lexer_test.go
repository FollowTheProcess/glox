package lexer_test

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/FollowTheProcess/glox/internal/syntax/lexer"
	"github.com/FollowTheProcess/glox/internal/syntax/token"
	"github.com/FollowTheProcess/test"
	"github.com/FollowTheProcess/txtar"
)

var update = flag.Bool("update", false, "Update snapshots and testdata")

const filePermissions = 0o644

func TestLexer(t *testing.T) {
	tests := []struct {
		name string        // Name of the test case
		src  string        // Source text to lex, string for convenience
		want []token.Token // Expected tokens
	}{
		{
			name: "empty",
			src:  "",
			want: []token.Token{
				{Kind: token.EOF},
			},
		},
		{
			name: "unexpected",
			src:  "%",
			want: []token.Token{
				{Kind: token.Error, Start: 0, End: 1},
				{Kind: token.EOF, Start: 1, End: 1},
			},
		},
		{
			name: "open paren",
			src:  "(",
			want: []token.Token{
				{Kind: token.OpenParen, Start: 0, End: 1},
				{Kind: token.EOF, Start: 1, End: 1},
			},
		},
		{
			name: "close paren",
			src:  ")",
			want: []token.Token{
				{Kind: token.CloseParen, Start: 0, End: 1},
				{Kind: token.EOF, Start: 1, End: 1},
			},
		},
		{
			name: "open brace",
			src:  "{",
			want: []token.Token{
				{Kind: token.OpenBrace, Start: 0, End: 1},
				{Kind: token.EOF, Start: 1, End: 1},
			},
		},
		{
			name: "close brace",
			src:  "}",
			want: []token.Token{
				{Kind: token.CloseBrace, Start: 0, End: 1},
				{Kind: token.EOF, Start: 1, End: 1},
			},
		},
		{
			name: "comma",
			src:  ",",
			want: []token.Token{
				{Kind: token.Comma, Start: 0, End: 1},
				{Kind: token.EOF, Start: 1, End: 1},
			},
		},
		{
			name: "dot",
			src:  ".",
			want: []token.Token{
				{Kind: token.Dot, Start: 0, End: 1},
				{Kind: token.EOF, Start: 1, End: 1},
			},
		},
		{
			name: "minus",
			src:  "-",
			want: []token.Token{
				{Kind: token.Minus, Start: 0, End: 1},
				{Kind: token.EOF, Start: 1, End: 1},
			},
		},
		{
			name: "plus",
			src:  "+",
			want: []token.Token{
				{Kind: token.Plus, Start: 0, End: 1},
				{Kind: token.EOF, Start: 1, End: 1},
			},
		},
		{
			name: "semicolon",
			src:  ";",
			want: []token.Token{
				{Kind: token.SemiColon, Start: 0, End: 1},
				{Kind: token.EOF, Start: 1, End: 1},
			},
		},
		{
			name: "forward slash",
			src:  "/",
			want: []token.Token{
				{Kind: token.ForwardSlash, Start: 0, End: 1},
				{Kind: token.EOF, Start: 1, End: 1},
			},
		},
		{
			name: "star",
			src:  "*",
			want: []token.Token{
				{Kind: token.Star, Start: 0, End: 1},
				{Kind: token.EOF, Start: 1, End: 1},
			},
		},
		{
			name: "bang",
			src:  "!",
			want: []token.Token{
				{Kind: token.Bang, Start: 0, End: 1},
				{Kind: token.EOF, Start: 1, End: 1},
			},
		},
		{
			name: "bang equal",
			src:  "!=",
			want: []token.Token{
				{Kind: token.BangEqual, Start: 0, End: 2},
				{Kind: token.EOF, Start: 2, End: 2},
			},
		},
		{
			name: "equal",
			src:  "=",
			want: []token.Token{
				{Kind: token.Equal, Start: 0, End: 1},
				{Kind: token.EOF, Start: 1, End: 1},
			},
		},
		{
			name: "double equal",
			src:  "==",
			want: []token.Token{
				{Kind: token.DoubleEqual, Start: 0, End: 2},
				{Kind: token.EOF, Start: 2, End: 2},
			},
		},
		{
			name: "greater than",
			src:  ">",
			want: []token.Token{
				{Kind: token.GreaterThan, Start: 0, End: 1},
				{Kind: token.EOF, Start: 1, End: 1},
			},
		},
		{
			name: "greater than equal",
			src:  ">=",
			want: []token.Token{
				{Kind: token.GreaterThanEqual, Start: 0, End: 2},
				{Kind: token.EOF, Start: 2, End: 2},
			},
		},
		{
			name: "less than",
			src:  "<",
			want: []token.Token{
				{Kind: token.LessThan, Start: 0, End: 1},
				{Kind: token.EOF, Start: 1, End: 1},
			},
		},
		{
			name: "less than equal",
			src:  "<=",
			want: []token.Token{
				{Kind: token.LessThanEqual, Start: 0, End: 2},
				{Kind: token.EOF, Start: 2, End: 2},
			},
		},
		{
			name: "comment",
			src:  "// I'm a comment to be ignored",
			want: []token.Token{
				{Kind: token.EOF, Start: 30, End: 30},
			},
		},
		{
			name: "ignore whitespace",
			src:  "  \t\t\n\n ()!=",
			want: []token.Token{
				{Kind: token.OpenParen, Start: 7, End: 8},
				{Kind: token.CloseParen, Start: 8, End: 9},
				{Kind: token.BangEqual, Start: 9, End: 11},
				{Kind: token.EOF, Start: 11, End: 11},
			},
		},
		{
			name: "string",
			src:  `"I'm a string literal"`,
			want: []token.Token{
				{Kind: token.String, Start: 0, End: 22},
				{Kind: token.EOF, Start: 22, End: 22},
			},
		},
		{
			name: "unterminated string",
			src:  `"I'm a string literal`,
			want: []token.Token{
				{Kind: token.Error, Start: 21, End: 21},
				{Kind: token.EOF, Start: 21, End: 21},
			},
		},
		{
			name: "integer",
			src:  "42",
			want: []token.Token{
				{Kind: token.Number, Start: 0, End: 2},
				{Kind: token.EOF, Start: 2, End: 2},
			},
		},
		{
			name: "float",
			src:  "42.69",
			want: []token.Token{
				{Kind: token.Number, Start: 0, End: 5},
				{Kind: token.EOF, Start: 5, End: 5},
			},
		},
		{
			name: "ident",
			src:  "some_variable",
			want: []token.Token{
				{Kind: token.Ident, Start: 0, End: 13},
				{Kind: token.EOF, Start: 13, End: 13},
			},
		},
		{
			name: "keyword if",
			src:  "if",
			want: []token.Token{
				{Kind: token.If, Start: 0, End: 2},
				{Kind: token.EOF, Start: 2, End: 2},
			},
		},
		{
			name: "keyword else",
			src:  "else",
			want: []token.Token{
				{Kind: token.Else, Start: 0, End: 4},
				{Kind: token.EOF, Start: 4, End: 4},
			},
		},
		{
			name: "keyword or",
			src:  "or",
			want: []token.Token{
				{Kind: token.Or, Start: 0, End: 2},
				{Kind: token.EOF, Start: 2, End: 2},
			},
		},
		{
			name: "keyword and",
			src:  "and",
			want: []token.Token{
				{Kind: token.And, Start: 0, End: 3},
				{Kind: token.EOF, Start: 3, End: 3},
			},
		},
		{
			name: "keyword for",
			src:  "for",
			want: []token.Token{
				{Kind: token.For, Start: 0, End: 3},
				{Kind: token.EOF, Start: 3, End: 3},
			},
		},
		{
			name: "keyword while",
			src:  "while",
			want: []token.Token{
				{Kind: token.While, Start: 0, End: 5},
				{Kind: token.EOF, Start: 5, End: 5},
			},
		},
		{
			name: "keyword true",
			src:  "true",
			want: []token.Token{
				{Kind: token.True, Start: 0, End: 4},
				{Kind: token.EOF, Start: 4, End: 4},
			},
		},
		{
			name: "keyword false",
			src:  "false",
			want: []token.Token{
				{Kind: token.False, Start: 0, End: 5},
				{Kind: token.EOF, Start: 5, End: 5},
			},
		},
		{
			name: "keyword class",
			src:  "class",
			want: []token.Token{
				{Kind: token.Class, Start: 0, End: 5},
				{Kind: token.EOF, Start: 5, End: 5},
			},
		},
		{
			name: "keyword super",
			src:  "super",
			want: []token.Token{
				{Kind: token.Super, Start: 0, End: 5},
				{Kind: token.EOF, Start: 5, End: 5},
			},
		},
		{
			name: "keyword this",
			src:  "this",
			want: []token.Token{
				{Kind: token.This, Start: 0, End: 4},
				{Kind: token.EOF, Start: 4, End: 4},
			},
		},
		{
			name: "keyword fun",
			src:  "fun",
			want: []token.Token{
				{Kind: token.Fun, Start: 0, End: 3},
				{Kind: token.EOF, Start: 3, End: 3},
			},
		},
		{
			name: "keyword var",
			src:  "var",
			want: []token.Token{
				{Kind: token.Var, Start: 0, End: 3},
				{Kind: token.EOF, Start: 3, End: 3},
			},
		},
		{
			name: "keyword nil",
			src:  "nil",
			want: []token.Token{
				{Kind: token.Nil, Start: 0, End: 3},
				{Kind: token.EOF, Start: 3, End: 3},
			},
		},
		{
			name: "keyword print",
			src:  "print",
			want: []token.Token{
				{Kind: token.Print, Start: 0, End: 5},
				{Kind: token.EOF, Start: 5, End: 5},
			},
		},
		{
			name: "keyword return",
			src:  "return",
			want: []token.Token{
				{Kind: token.Return, Start: 0, End: 6},
				{Kind: token.EOF, Start: 6, End: 6},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tokens := collect(tt.src)
			test.EqualFunc(t, tokens, tt.want, tokenStreamEqual)
		})
	}
}

func TestLexerIntegration(t *testing.T) {
	base := filepath.Join("testdata", "lex")
	pattern := filepath.Join(base, "*.txtar")
	files, err := filepath.Glob(pattern)
	test.Ok(t, err)

	for _, file := range files {
		name := filepath.Base(file)
		t.Run(name, func(t *testing.T) {
			archive, err := txtar.ParseFile(file)
			test.Ok(t, err)

			src, err := archive.Read("src.lox")
			test.Ok(t, err)

			expected, err := archive.Read("expected.txt")
			test.Ok(t, err)

			tokens := collectBytes(src)

			var formattedTokens strings.Builder
			for _, tok := range tokens {
				formattedTokens.WriteString(tok.String())
				formattedTokens.WriteByte('\n')
			}

			got := formattedTokens.String()

			if *update {
				// Update the expected with what's actually been seen
				err := archive.Add("expected.txt", []byte(got))
				test.Ok(t, err)

				err = os.WriteFile(file, []byte(archive.String()), filePermissions)
				test.Ok(t, err)
				return
			}

			test.Diff(t, got, string(expected))
		})
	}
}

func BenchmarkLexer(b *testing.B) {
	file := filepath.Join("testdata", "bench", "binary_trees.lox")

	contents, err := os.ReadFile(file)
	test.Ok(b, err)

	for b.Loop() {
		// Must initialise the lexer inside the loop as it's internal state is
		// modified on each scan
		lex := lexer.New(contents)
		for {
			tok := lex.NextToken()
			if tok.Is(token.EOF) || tok.Is(token.Error) {
				break
			}
		}
	}
}

// collect gathers the emitted tokens into a slice for comparison.
func collect(src string) []token.Token {
	var tokens []token.Token
	l := lexer.New([]byte(src))
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Kind == token.EOF {
			break
		}
	}
	return tokens
}

// collect gathers the emitted tokens into a slice for comparison.
func collectBytes(src []byte) []token.Token {
	var tokens []token.Token
	l := lexer.New(src)
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Kind == token.EOF {
			break
		}
	}
	return tokens
}

// tokenStreamEqual compares to slices of tokens for equality.
func tokenStreamEqual(t1, t2 []token.Token) bool {
	if len(t1) != len(t2) {
		return false
	}
	for i := range t1 {
		if t1[i].Kind != t2[i].Kind {
			return false
		}

		if t1[i].Start != t2[i].Start {
			return false
		}

		if t1[i].End != t2[i].End {
			return false
		}
	}
	return true
}
