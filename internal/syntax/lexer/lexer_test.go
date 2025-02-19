package lexer_test

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/FollowTheProcess/glox/internal/syntax/lexer"
	"github.com/FollowTheProcess/glox/internal/syntax/token"
	"github.com/FollowTheProcess/test"
	"golang.org/x/tools/txtar"
)

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
				{Kind: token.Error, Text: []byte("unexpected char %"), Offset: 0},
				{Kind: token.EOF, Offset: 0},
			},
		},
		{
			name: "open paren",
			src:  "(",
			want: []token.Token{
				{Kind: token.OpenParen, Text: []byte("("), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "close paren",
			src:  ")",
			want: []token.Token{
				{Kind: token.CloseParen, Text: []byte(")"), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "open brace",
			src:  "{",
			want: []token.Token{
				{Kind: token.OpenBrace, Text: []byte("{"), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "close brace",
			src:  "}",
			want: []token.Token{
				{Kind: token.CloseBrace, Text: []byte("}"), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "comma",
			src:  ",",
			want: []token.Token{
				{Kind: token.Comma, Text: []byte(","), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "dot",
			src:  ".",
			want: []token.Token{
				{Kind: token.Dot, Text: []byte("."), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "minus",
			src:  "-",
			want: []token.Token{
				{Kind: token.Minus, Text: []byte("-"), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "plus",
			src:  "+",
			want: []token.Token{
				{Kind: token.Plus, Text: []byte("+"), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "semicolon",
			src:  ";",
			want: []token.Token{
				{Kind: token.SemiColon, Text: []byte(";"), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "forward slash",
			src:  "/",
			want: []token.Token{
				{Kind: token.ForwardSlash, Text: []byte("/"), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "star",
			src:  "*",
			want: []token.Token{
				{Kind: token.Star, Text: []byte("*"), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "bang",
			src:  "!",
			want: []token.Token{
				{Kind: token.Bang, Text: []byte("!"), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "bang equal",
			src:  "!=",
			want: []token.Token{
				{Kind: token.BangEqual, Text: []byte("!="), Offset: 0, Width: 2},
				{Kind: token.EOF, Offset: 2},
			},
		},
		{
			name: "equal",
			src:  "=",
			want: []token.Token{
				{Kind: token.Equal, Text: []byte("="), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "double equal",
			src:  "==",
			want: []token.Token{
				{Kind: token.DoubleEqual, Text: []byte("=="), Offset: 0, Width: 2},
				{Kind: token.EOF, Offset: 2},
			},
		},
		{
			name: "greater than",
			src:  ">",
			want: []token.Token{
				{Kind: token.GreaterThan, Text: []byte(">"), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "greater than equal",
			src:  ">=",
			want: []token.Token{
				{Kind: token.GreaterThanEqual, Text: []byte(">="), Offset: 0, Width: 2},
				{Kind: token.EOF, Offset: 2},
			},
		},
		{
			name: "less than",
			src:  "<",
			want: []token.Token{
				{Kind: token.LessThan, Text: []byte("<"), Offset: 0, Width: 1},
				{Kind: token.EOF, Offset: 1},
			},
		},
		{
			name: "less than equal",
			src:  "<=",
			want: []token.Token{
				{Kind: token.LessThanEqual, Text: []byte("<="), Offset: 0, Width: 2},
				{Kind: token.EOF, Offset: 2},
			},
		},
		{
			name: "comment",
			src:  "// I'm a comment to be ignored",
			want: []token.Token{
				{Kind: token.EOF, Offset: 30},
			},
		},
		{
			name: "ignore whitespace",
			src:  "  \t\t\n\n ()!=",
			want: []token.Token{
				{Kind: token.OpenParen, Text: []byte("("), Offset: 7, Width: 1},
				{Kind: token.CloseParen, Text: []byte(")"), Offset: 8, Width: 1},
				{Kind: token.BangEqual, Text: []byte("!="), Offset: 9, Width: 2},
				{Kind: token.EOF, Offset: 11},
			},
		},
		{
			name: "string",
			src:  `"I'm a string literal"`,
			want: []token.Token{
				{Kind: token.String, Text: []byte(`"I'm a string literal"`), Offset: 0, Width: 22},
				{Kind: token.EOF, Offset: 22},
			},
		},
		{
			name: "unterminated string",
			src:  `"I'm a string literal`,
			want: []token.Token{
				{Kind: token.Error, Text: []byte("unterminated string literal"), Offset: 21},
				{Kind: token.EOF, Offset: 21},
			},
		},
		{
			name: "integer",
			src:  "42",
			want: []token.Token{
				{Kind: token.Number, Text: []byte("42"), Offset: 0, Width: 2},
				{Kind: token.EOF, Offset: 2},
			},
		},
		{
			name: "float",
			src:  "42.69",
			want: []token.Token{
				{Kind: token.Number, Text: []byte("42.69"), Offset: 0, Width: 5},
				{Kind: token.EOF, Offset: 5},
			},
		},
		{
			name: "ident",
			src:  "some_variable",
			want: []token.Token{
				{Kind: token.Ident, Text: []byte("some_variable"), Offset: 0, Width: 13},
				{Kind: token.EOF, Offset: 13},
			},
		},
		{
			name: "keyword if",
			src:  "if",
			want: []token.Token{
				{Kind: token.If, Text: []byte("if"), Offset: 0, Width: 2},
				{Kind: token.EOF, Offset: 2},
			},
		},
		{
			name: "keyword else",
			src:  "else",
			want: []token.Token{
				{Kind: token.Else, Text: []byte("else"), Offset: 0, Width: 4},
				{Kind: token.EOF, Offset: 4},
			},
		},
		{
			name: "keyword or",
			src:  "or",
			want: []token.Token{
				{Kind: token.Or, Text: []byte("or"), Offset: 0, Width: 2},
				{Kind: token.EOF, Offset: 2},
			},
		},
		{
			name: "keyword and",
			src:  "and",
			want: []token.Token{
				{Kind: token.And, Text: []byte("and"), Offset: 0, Width: 3},
				{Kind: token.EOF, Offset: 3},
			},
		},
		{
			name: "keyword for",
			src:  "for",
			want: []token.Token{
				{Kind: token.For, Text: []byte("for"), Offset: 0, Width: 3},
				{Kind: token.EOF, Offset: 3},
			},
		},
		{
			name: "keyword while",
			src:  "while",
			want: []token.Token{
				{Kind: token.While, Text: []byte("while"), Offset: 0, Width: 5},
				{Kind: token.EOF, Offset: 5},
			},
		},
		{
			name: "keyword true",
			src:  "true",
			want: []token.Token{
				{Kind: token.True, Text: []byte("true"), Offset: 0, Width: 4},
				{Kind: token.EOF, Offset: 4},
			},
		},
		{
			name: "keyword false",
			src:  "false",
			want: []token.Token{
				{Kind: token.False, Text: []byte("false"), Offset: 0, Width: 5},
				{Kind: token.EOF, Offset: 5},
			},
		},
		{
			name: "keyword class",
			src:  "class",
			want: []token.Token{
				{Kind: token.Class, Text: []byte("class"), Offset: 0, Width: 5},
				{Kind: token.EOF, Offset: 5},
			},
		},
		{
			name: "keyword super",
			src:  "super",
			want: []token.Token{
				{Kind: token.Super, Text: []byte("super"), Offset: 0, Width: 5},
				{Kind: token.EOF, Offset: 5},
			},
		},
		{
			name: "keyword this",
			src:  "this",
			want: []token.Token{
				{Kind: token.This, Text: []byte("this"), Offset: 0, Width: 4},
				{Kind: token.EOF, Offset: 4},
			},
		},
		{
			name: "keyword fun",
			src:  "fun",
			want: []token.Token{
				{Kind: token.Fun, Text: []byte("fun"), Offset: 0, Width: 3},
				{Kind: token.EOF, Offset: 3},
			},
		},
		{
			name: "keyword var",
			src:  "var",
			want: []token.Token{
				{Kind: token.Var, Text: []byte("var"), Offset: 0, Width: 3},
				{Kind: token.EOF, Offset: 3},
			},
		},
		{
			name: "keyword nil",
			src:  "nil",
			want: []token.Token{
				{Kind: token.Nil, Text: []byte("nil"), Offset: 0, Width: 3},
				{Kind: token.EOF, Offset: 3},
			},
		},
		{
			name: "keyword print",
			src:  "print",
			want: []token.Token{
				{Kind: token.Print, Text: []byte("print"), Offset: 0, Width: 5},
				{Kind: token.EOF, Offset: 5},
			},
		},
		{
			name: "keyword return",
			src:  "return",
			want: []token.Token{
				{Kind: token.Return, Text: []byte("return"), Offset: 0, Width: 6},
				{Kind: token.EOF, Offset: 6},
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

			if len(archive.Files) != 2 {
				t.Fatalf("%s expected to contain 2 files, actual: %d", file, len(archive.Files))
			}

			if archive.Files[0].Name != "src.lox" {
				t.Fatalf("%s expected first file to be 'src.lox' got %s", file, archive.Files[0].Name)
			}

			if archive.Files[1].Name != "expected.txt" {
				t.Fatalf("%s expected second file to be 'expected.txt' got %s", file, archive.Files[1].Name)
			}

			src := archive.Files[0].Data
			expected := archive.Files[1].Data

			tokens := collectBytes(src)

			var formattedTokens strings.Builder
			for _, tok := range tokens {
				formattedTokens.WriteString(tok.String())
				formattedTokens.WriteByte('\n')
			}

			test.Diff(t, formattedTokens.String(), string(expected))
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
		if !bytes.Equal(t1[i].Text, t2[i].Text) {
			return false
		}

		if t1[i].Offset != t2[i].Offset {
			return false
		}

		if t1[i].Width != t2[i].Width {
			return false
		}
	}
	return true
}
