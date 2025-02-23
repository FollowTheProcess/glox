package token_test

import (
	"testing"

	"github.com/FollowTheProcess/glox/internal/syntax/token"
	"github.com/FollowTheProcess/test"
)

func TestToken(t *testing.T) {
	tests := []struct {
		name string      // Name of the test case
		want string      // Expected return value from String()
		tok  token.Token // The token under test
	}{
		{
			name: "bad",
			tok:  token.Token{Kind: token.Kind(9999)},
			want: `<Token::<BadToken> start=0, end=0>`,
		},
		{
			name: "eof",
			tok:  token.Token{Kind: token.EOF},
			want: `<Token::EOF start=0, end=0>`,
		},
		{
			name: "error",
			tok:  token.Token{Kind: token.Error, Start: 42, End: 4},
			want: `<Token::Error start=42, end=4>`,
		},
		{
			name: "open paren",
			tok:  token.Token{Kind: token.OpenParen, Start: 1, End: 1},
			want: `<Token::OpenParen start=1, end=1>`,
		},
		{
			name: "close paren",
			tok:  token.Token{Kind: token.CloseParen, Start: 1, End: 1},
			want: `<Token::CloseParen start=1, end=1>`,
		},
		{
			name: "open brace",
			tok:  token.Token{Kind: token.OpenBrace, Start: 1, End: 1},
			want: `<Token::OpenBrace start=1, end=1>`,
		},
		{
			name: "close brace",
			tok:  token.Token{Kind: token.CloseBrace, Start: 1, End: 1},
			want: `<Token::CloseBrace start=1, end=1>`,
		},
		{
			name: "comma",
			tok:  token.Token{Kind: token.Comma, Start: 27, End: 1},
			want: `<Token::Comma start=27, end=1>`,
		},
		{
			name: "dot",
			tok:  token.Token{Kind: token.Dot, Start: 2, End: 1},
			want: `<Token::Dot start=2, end=1>`,
		},
		{
			name: "minus",
			tok:  token.Token{Kind: token.Minus, Start: 32, End: 1},
			want: `<Token::Minus start=32, end=1>`,
		},
		{
			name: "plus",
			tok:  token.Token{Kind: token.Plus, Start: 185, End: 1},
			want: `<Token::Plus start=185, end=1>`,
		},
		{
			name: "semicolon",
			tok:  token.Token{Kind: token.SemiColon, Start: 86, End: 1},
			want: `<Token::SemiColon start=86, end=1>`,
		},
		{
			name: "forward slash",
			tok:  token.Token{Kind: token.ForwardSlash, Start: 17, End: 1},
			want: `<Token::ForwardSlash start=17, end=1>`,
		},
		{
			name: "star",
			tok:  token.Token{Kind: token.Star, Start: 12, End: 1},
			want: `<Token::Star start=12, end=1>`,
		},
		{
			name: "bang",
			tok:  token.Token{Kind: token.Bang, Start: 7, End: 1},
			want: `<Token::Bang start=7, end=1>`,
		},
		{
			name: "equal",
			tok:  token.Token{Kind: token.Equal, Start: 2, End: 1},
			want: `<Token::Equal start=2, end=1>`,
		},
		{
			name: "bang equal",
			tok:  token.Token{Kind: token.BangEqual, Start: 1, End: 2},
			want: `<Token::BangEqual start=1, end=2>`,
		},
		{
			name: "double equal",
			tok:  token.Token{Kind: token.DoubleEqual, Start: 174, End: 2},
			want: `<Token::DoubleEqual start=174, end=2>`,
		},
		{
			name: "greater than",
			tok:  token.Token{Kind: token.GreaterThan, Start: 22, End: 1},
			want: `<Token::GreaterThan start=22, end=1>`,
		},
		{
			name: "less than",
			tok:  token.Token{Kind: token.LessThan, Start: 63, End: 1},
			want: `<Token::LessThan start=63, end=1>`,
		},
		{
			name: "greater than equal",
			tok:  token.Token{Kind: token.GreaterThanEqual, Start: 3, End: 2},
			want: `<Token::GreaterThanEqual start=3, end=2>`,
		},
		{
			name: "less than equal",
			tok:  token.Token{Kind: token.LessThanEqual, Start: 7, End: 2},
			want: `<Token::LessThanEqual start=7, end=2>`,
		},
		{
			name: "string",
			tok:  token.Token{Kind: token.String, Start: 1, End: 22},
			want: `<Token::String start=1, end=22>`,
		},
		{
			name: "number",
			tok:  token.Token{Kind: token.Number, Start: 0, End: 2},
			want: `<Token::Number start=0, end=2>`,
		},
		{
			name: "ident",
			tok:  token.Token{Kind: token.Ident, Start: 0, End: 9},
			want: `<Token::Ident start=0, end=9>`,
		},
		{
			name: "if",
			tok:  token.Token{Kind: token.If, Start: 17, End: 2},
			want: `<Token::If start=17, end=2>`,
		},
		{
			name: "else",
			tok:  token.Token{Kind: token.Else, Start: 37, End: 4},
			want: `<Token::Else start=37, end=4>`,
		},
		{
			name: "or",
			tok:  token.Token{Kind: token.Or, Start: 145, End: 2},
			want: `<Token::Or start=145, end=2>`,
		},
		{
			name: "and",
			tok:  token.Token{Kind: token.And, Start: 1, End: 3},
			want: `<Token::And start=1, end=3>`,
		},
		{
			name: "for",
			tok:  token.Token{Kind: token.For, Start: 5, End: 3},
			want: `<Token::For start=5, end=3>`,
		},
		{
			name: "while",
			tok:  token.Token{Kind: token.While, Start: 2, End: 5},
			want: `<Token::While start=2, end=5>`,
		},
		{
			name: "true",
			tok:  token.Token{Kind: token.True, Start: 0, End: 4},
			want: `<Token::True start=0, end=4>`,
		},
		{
			name: "false",
			tok:  token.Token{Kind: token.False, Start: 19, End: 5},
			want: `<Token::False start=19, end=5>`,
		},
		{
			name: "class",
			tok:  token.Token{Kind: token.Class, Start: 21, End: 5},
			want: `<Token::Class start=21, end=5>`,
		},
		{
			name: "super",
			tok:  token.Token{Kind: token.Super, Start: 67, End: 5},
			want: `<Token::Super start=67, end=5>`,
		},
		{
			name: "this",
			tok:  token.Token{Kind: token.This, Start: 2, End: 4},
			want: `<Token::This start=2, end=4>`,
		},
		{
			name: "fun",
			tok:  token.Token{Kind: token.Fun, Start: 0, End: 3},
			want: `<Token::Fun start=0, end=3>`,
		},
		{
			name: "var",
			tok:  token.Token{Kind: token.Var, Start: 73, End: 3},
			want: `<Token::Var start=73, end=3>`,
		},
		{
			name: "nil",
			tok:  token.Token{Kind: token.Nil, Start: 189, End: 3},
			want: `<Token::Nil start=189, end=3>`,
		},
		{
			name: "print",
			tok:  token.Token{Kind: token.Print, Start: 0, End: 5},
			want: `<Token::Print start=0, end=5>`,
		},
		{
			name: "return",
			tok:  token.Token{Kind: token.Return, Start: 17, End: 6},
			want: `<Token::Return start=17, end=6>`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test.Equal(t, tt.tok.String(), tt.want)
			test.True(t, tt.tok.Is(tt.tok.Kind))
		})
	}
}

func TestKeyword(t *testing.T) {
	tests := []struct {
		name  string     // Name of the test
		ident string     // Ident to lookup in the set of keywords
		want  token.Kind // Expected return kind
		ok    bool       // Expected return bool
	}{
		{name: "empty", ident: "", want: token.Ident, ok: false},
		{name: "garbage", ident: "$%^&*vhbjj", want: token.Ident, ok: false},
		{name: "non keyword", ident: "some_variable", want: token.Ident, ok: false},
		{name: "if", ident: "if", want: token.If, ok: true},
		{name: "else", ident: "else", want: token.Else, ok: true},
		{name: "or", ident: "or", want: token.Or, ok: true},
		{name: "and", ident: "and", want: token.And, ok: true},
		{name: "for", ident: "for", want: token.For, ok: true},
		{name: "while", ident: "while", want: token.While, ok: true},
		{name: "true", ident: "true", want: token.True, ok: true},
		{name: "false", ident: "false", want: token.False, ok: true},
		{name: "class", ident: "class", want: token.Class, ok: true},
		{name: "super", ident: "super", want: token.Super, ok: true},
		{name: "this", ident: "this", want: token.This, ok: true},
		{name: "fun", ident: "fun", want: token.Fun, ok: true},
		{name: "var", ident: "var", want: token.Var, ok: true},
		{name: "nil", ident: "nil", want: token.Nil, ok: true},
		{name: "print", ident: "print", want: token.Print, ok: true},
		{name: "return", ident: "return", want: token.Return, ok: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := token.Keyword(tt.ident)

			test.Equal(t, ok, tt.ok)
			test.Equal(t, got, tt.want)
		})
	}
}

func TestPrecedence(t *testing.T) {
	tests := []struct {
		kind token.Kind // The token kind under test
		want int        // The expected precedence
	}{
		{kind: token.Or, want: token.PrecedenceOr},
		{kind: token.And, want: token.PrecedenceAnd},
		{kind: token.Equal, want: token.PrecedenceComp},
		{kind: token.BangEqual, want: token.PrecedenceComp},
		{kind: token.LessThan, want: token.PrecedenceComp},
		{kind: token.LessThanEqual, want: token.PrecedenceComp},
		{kind: token.GreaterThan, want: token.PrecedenceComp},
		{kind: token.GreaterThanEqual, want: token.PrecedenceComp},
		{kind: token.Plus, want: token.PrecedenceAddSubtract},
		{kind: token.Minus, want: token.PrecedenceAddSubtract},
		{kind: token.Star, want: token.PrecedenceMulDivide},
		{kind: token.ForwardSlash, want: token.PrecedenceMulDivide},
		{kind: token.Bang, want: token.PrecedenceMin},
		{kind: token.Bang, want: token.PrecedenceMin},
		{kind: token.Number, want: token.PrecedenceMin},
		{kind: token.Ident, want: token.PrecedenceMin},
	}

	for _, tt := range tests {
		t.Run(tt.kind.String(), func(t *testing.T) {
			tok := token.Token{Kind: tt.kind}

			test.Equal(t, tok.Precedence(), tt.want)
		})
	}
}
