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
			want: `<Token::<BadToken> text="", offset=0, width=0>`,
		},
		{
			name: "eof",
			tok:  token.Token{Kind: token.EOF},
			want: `<Token::EOF text="", offset=0, width=0>`,
		},
		{
			name: "error",
			tok:  token.Token{Kind: token.Error, Text: []byte("bang"), Offset: 42, Width: 4},
			want: `<Token::Error text="bang", offset=42, width=4>`,
		},
		{
			name: "open paren",
			tok:  token.Token{Kind: token.OpenParen, Text: []byte("("), Offset: 1, Width: 1},
			want: `<Token::OpenParen text="(", offset=1, width=1>`,
		},
		{
			name: "close paren",
			tok:  token.Token{Kind: token.CloseParen, Text: []byte(")"), Offset: 1, Width: 1},
			want: `<Token::CloseParen text=")", offset=1, width=1>`,
		},
		{
			name: "open brace",
			tok:  token.Token{Kind: token.OpenBrace, Text: []byte("{"), Offset: 1, Width: 1},
			want: `<Token::OpenBrace text="{", offset=1, width=1>`,
		},
		{
			name: "close brace",
			tok:  token.Token{Kind: token.CloseBrace, Text: []byte("}"), Offset: 1, Width: 1},
			want: `<Token::CloseBrace text="}", offset=1, width=1>`,
		},
		{
			name: "comma",
			tok:  token.Token{Kind: token.Comma, Text: []byte(","), Offset: 27, Width: 1},
			want: `<Token::Comma text=",", offset=27, width=1>`,
		},
		{
			name: "dot",
			tok:  token.Token{Kind: token.Dot, Text: []byte("."), Offset: 2, Width: 1},
			want: `<Token::Dot text=".", offset=2, width=1>`,
		},
		{
			name: "minus",
			tok:  token.Token{Kind: token.Minus, Text: []byte("-"), Offset: 32, Width: 1},
			want: `<Token::Minus text="-", offset=32, width=1>`,
		},
		{
			name: "plus",
			tok:  token.Token{Kind: token.Plus, Text: []byte("+"), Offset: 185, Width: 1},
			want: `<Token::Plus text="+", offset=185, width=1>`,
		},
		{
			name: "semicolon",
			tok:  token.Token{Kind: token.SemiColon, Text: []byte(";"), Offset: 86, Width: 1},
			want: `<Token::SemiColon text=";", offset=86, width=1>`,
		},
		{
			name: "forward slash",
			tok:  token.Token{Kind: token.ForwardSlash, Text: []byte("/"), Offset: 17, Width: 1},
			want: `<Token::ForwardSlash text="/", offset=17, width=1>`,
		},
		{
			name: "star",
			tok:  token.Token{Kind: token.Star, Text: []byte("*"), Offset: 12, Width: 1},
			want: `<Token::Star text="*", offset=12, width=1>`,
		},
		{
			name: "bang",
			tok:  token.Token{Kind: token.Bang, Text: []byte("!"), Offset: 7, Width: 1},
			want: `<Token::Bang text="!", offset=7, width=1>`,
		},
		{
			name: "equal",
			tok:  token.Token{Kind: token.Equal, Text: []byte("="), Offset: 2, Width: 1},
			want: `<Token::Equal text="=", offset=2, width=1>`,
		},
		{
			name: "bang equal",
			tok:  token.Token{Kind: token.BangEqual, Text: []byte("!="), Offset: 1, Width: 2},
			want: `<Token::BangEqual text="!=", offset=1, width=2>`,
		},
		{
			name: "double equal",
			tok:  token.Token{Kind: token.DoubleEqual, Text: []byte("=="), Offset: 174, Width: 2},
			want: `<Token::DoubleEqual text="==", offset=174, width=2>`,
		},
		{
			name: "greater than",
			tok:  token.Token{Kind: token.GreaterThan, Text: []byte(">"), Offset: 22, Width: 1},
			want: `<Token::GreaterThan text=">", offset=22, width=1>`,
		},
		{
			name: "less than",
			tok:  token.Token{Kind: token.LessThan, Text: []byte("<"), Offset: 63, Width: 1},
			want: `<Token::LessThan text="<", offset=63, width=1>`,
		},
		{
			name: "greater than equal",
			tok:  token.Token{Kind: token.GreaterThanEqual, Text: []byte(">="), Offset: 3, Width: 2},
			want: `<Token::GreaterThanEqual text=">=", offset=3, width=2>`,
		},
		{
			name: "less than equal",
			tok:  token.Token{Kind: token.LessThanEqual, Text: []byte("<="), Offset: 7, Width: 2},
			want: `<Token::LessThanEqual text="<=", offset=7, width=2>`,
		},
		{
			name: "string",
			tok:  token.Token{Kind: token.String, Text: []byte(`"I'm a string literal"`), Offset: 1, Width: 22},
			want: `<Token::String text="I'm a string literal", offset=1, width=22>`,
		},
		{
			name: "number",
			tok:  token.Token{Kind: token.Number, Text: []byte("42"), Offset: 0, Width: 2},
			want: `<Token::Number text="42", offset=0, width=2>`,
		},
		{
			name: "ident",
			tok:  token.Token{Kind: token.Ident, Text: []byte("something"), Offset: 0, Width: 9},
			want: `<Token::Ident text="something", offset=0, width=9>`,
		},
		{
			name: "if",
			tok:  token.Token{Kind: token.If, Text: []byte("if"), Offset: 17, Width: 2},
			want: `<Token::If text="if", offset=17, width=2>`,
		},
		{
			name: "else",
			tok:  token.Token{Kind: token.Else, Text: []byte("else"), Offset: 37, Width: 4},
			want: `<Token::Else text="else", offset=37, width=4>`,
		},
		{
			name: "or",
			tok:  token.Token{Kind: token.Or, Text: []byte("or"), Offset: 145, Width: 2},
			want: `<Token::Or text="or", offset=145, width=2>`,
		},
		{
			name: "and",
			tok:  token.Token{Kind: token.And, Text: []byte("and"), Offset: 1, Width: 3},
			want: `<Token::And text="and", offset=1, width=3>`,
		},
		{
			name: "for",
			tok:  token.Token{Kind: token.For, Text: []byte("for"), Offset: 5, Width: 3},
			want: `<Token::For text="for", offset=5, width=3>`,
		},
		{
			name: "while",
			tok:  token.Token{Kind: token.While, Text: []byte("while"), Offset: 2, Width: 5},
			want: `<Token::While text="while", offset=2, width=5>`,
		},
		{
			name: "true",
			tok:  token.Token{Kind: token.True, Text: []byte("true"), Offset: 0, Width: 4},
			want: `<Token::True text="true", offset=0, width=4>`,
		},
		{
			name: "false",
			tok:  token.Token{Kind: token.False, Text: []byte("false"), Offset: 19, Width: 5},
			want: `<Token::False text="false", offset=19, width=5>`,
		},
		{
			name: "class",
			tok:  token.Token{Kind: token.Class, Text: []byte("class"), Offset: 21, Width: 5},
			want: `<Token::Class text="class", offset=21, width=5>`,
		},
		{
			name: "super",
			tok:  token.Token{Kind: token.Super, Text: []byte("super"), Offset: 67, Width: 5},
			want: `<Token::Super text="super", offset=67, width=5>`,
		},
		{
			name: "this",
			tok:  token.Token{Kind: token.This, Text: []byte("this"), Offset: 2, Width: 4},
			want: `<Token::This text="this", offset=2, width=4>`,
		},
		{
			name: "fun",
			tok:  token.Token{Kind: token.Fun, Text: []byte("fun"), Offset: 0, Width: 3},
			want: `<Token::Fun text="fun", offset=0, width=3>`,
		},
		{
			name: "var",
			tok:  token.Token{Kind: token.Var, Text: []byte("var"), Offset: 73, Width: 3},
			want: `<Token::Var text="var", offset=73, width=3>`,
		},
		{
			name: "nil",
			tok:  token.Token{Kind: token.Nil, Text: []byte("nil"), Offset: 189, Width: 3},
			want: `<Token::Nil text="nil", offset=189, width=3>`,
		},
		{
			name: "print",
			tok:  token.Token{Kind: token.Print, Text: []byte("print"), Offset: 0, Width: 5},
			want: `<Token::Print text="print", offset=0, width=5>`,
		},
		{
			name: "return",
			tok:  token.Token{Kind: token.Return, Text: []byte("return"), Offset: 17, Width: 6},
			want: `<Token::Return text="return", offset=17, width=6>`,
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
		{
			name:  "empty",
			ident: "",
			want:  token.Error,
			ok:    false,
		},
		{
			name:  "garbage",
			ident: "$%^&*vhbjj",
			want:  token.Error,
			ok:    false,
		},
		{
			name:  "non keyword",
			ident: "some_variable",
			want:  token.Error,
			ok:    false,
		},
		{
			name:  "if",
			ident: "if",
			want:  token.If,
			ok:    true,
		},
		{
			name:  "while",
			ident: "while",
			want:  token.While,
			ok:    true,
		},
		{
			name:  "super",
			ident: "super",
			want:  token.Super,
			ok:    true,
		},
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
