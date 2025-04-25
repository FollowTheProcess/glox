// Package lexer implement the lexical scanner for the lox language.
package lexer

import (
	"fmt"
	"unicode"
	"unicode/utf8"

	"github.com/FollowTheProcess/glox/internal/syntax"
	"github.com/FollowTheProcess/glox/internal/syntax/token"
)

const eof = rune(0) // eof signifies we have reached the end of the input.

// Lexer is the lexer.
type Lexer struct {
	handler   syntax.ErrorHandler // The error handler, if any
	name      string              // Filename
	src       string              // The raw source text
	line      int                 // Current line number
	lineStart int                 // Offset at which the current line started
	start     int                 // The starting offset of the current token
	pos       int                 // The current position in src
	width     int                 // The width of the last rune read, allows backup
}

// New returns a new [Lexer].
func New(name, src string, handler syntax.ErrorHandler) *Lexer {
	return &Lexer{
		handler: handler,
		src:     src,
		name:    name,
		line:    1,
	}
}

// NextToken returns the next token from the input stream.
func (l *Lexer) NextToken() token.Token {
	l.skip(unicode.IsSpace)

	switch char := l.next(); char {
	case eof:
		return l.emit(token.EOF)
	case '(':
		return l.emit(token.OpenParen)
	case ')':
		return l.emit(token.CloseParen)
	case '{':
		return l.emit(token.OpenBrace)
	case '}':
		return l.emit(token.CloseBrace)
	case ',':
		return l.emit(token.Comma)
	case '.':
		return l.emit(token.Dot)
	case '-':
		return l.emit(token.Minus)
	case '+':
		return l.emit(token.Plus)
	case ';':
		return l.emit(token.SemiColon)
	case '*':
		return l.emit(token.Star)
	case '/':
		return l.scanSlash()
	case '=':
		return l.scanEq()
	case '!':
		return l.scanBang()
	case '>':
		return l.scanGreater()
	case '<':
		return l.scanLess()
	default:
		l.errorf("unexpected character %q", char)
		return l.emit(token.Error)
	}
}

// next returns the next utf8 rune in the input, or [eof], and advances
// the lexer over that rune such that successive calls of [Lexer.next] iterate
// through the src one rune at a time.
func (l *Lexer) next() rune {
	if l.pos >= len(l.src) {
		return eof
	}

	char, width := utf8.DecodeRuneInString(l.src[l.pos:])
	l.width = width
	l.pos += width

	if char == '\n' {
		l.line++
		l.lineStart = l.pos
	}

	return char
}

// peek returns, but does not consume, the next utf8 rune in the input.
func (l *Lexer) peek() rune {
	if l.pos >= len(l.src) {
		return eof
	}

	char, _ := utf8.DecodeRuneInString(l.src[l.pos:])
	return char
}

// skip ignores any characters for which the predicate returns true, stopping at the first one
// that returns false such that after it returns, [Lexer.next] returns the first 'false' char.
func (l *Lexer) skip(predicate func(r rune) bool) {
	for predicate(l.peek()) {
		l.next()
	}
	l.start = l.pos
}

// emit returns a [token.Token] of the given kind using the lexer's internal
// state to fill in the position information.
func (l *Lexer) emit(kind token.Kind) token.Token {
	tok := token.Token{
		Kind:  kind,
		Line:  l.line,
		Start: l.start,
		End:   l.pos,
	}

	l.start = l.pos

	return tok
}

// error calculates the position information and calls the error handler.
func (l *Lexer) error(msg string) {
	if l.handler == nil {
		// Nothing to do
		return
	}

	// Column is the number of bytes between the last newline and the current position
	// +1 because columns are 1 indexed
	startCol := 1 + l.start - l.lineStart
	endCol := 1 + l.pos - l.lineStart

	position := syntax.Position{
		Name:     l.name,
		Offset:   l.pos,
		Line:     l.line,
		StartCol: startCol,
		EndCol:   endCol,
	}

	l.handler(position, msg)
}

// errorf calls error with a formatted message.
func (l *Lexer) errorf(format string, a ...any) {
	l.error(fmt.Sprintf(format, a...))
}

// scanSlash scans a '/' character.
func (l *Lexer) scanSlash() token.Token {
	if l.peek() == '/' {
		// Comments go to the end of the line
		for l.peek() != '\n' && l.peek() != eof {
			l.next()
		}
		// Return to the base state
		return l.NextToken()
	}
	return l.emit(token.Slash)
}

// scanEq scans a '=' character.
func (l *Lexer) scanEq() token.Token {
	if l.peek() == '=' {
		l.next()
		return l.emit(token.DoubleEq)
	}
	return l.emit(token.Eq)
}

// scanBang scans a '!' character.
func (l *Lexer) scanBang() token.Token {
	if l.peek() == '=' {
		l.next() // Consume the '='
		return l.emit(token.BangEq)
	}
	return l.emit(token.Bang)
}

// scanGreater scans a '>' character.
func (l *Lexer) scanGreater() token.Token {
	if l.peek() == '=' {
		l.next()
		return l.emit(token.GreaterEq)
	}
	return l.emit(token.Greater)
}

// scanLess scans a '<' character.
func (l *Lexer) scanLess() token.Token {
	if l.peek() == '=' {
		l.next()
		return l.emit(token.LessEq)
	}
	return l.emit(token.Less)
}
