// Package lexer implement the lexical scanner for the lox language.
package lexer

import (
	"fmt"
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
	switch char := l.next(); char {
	case eof:
		return l.emit(token.EOF)
	case '(':
		return l.emit(token.OpenParen)
	default:
		l.errorf("invalid character %c", char)
		return l.emit(token.Invalid)
	}
}

// next returns the next utf8 rune in the input, or [eof], and advances
// the lexer over that rune such that successive calls of [Lexer.next] iterate
// through the src one rune at a time.
func (l *Lexer) next() rune {
	if l.pos >= len(l.src) {
		return eof
	}

	char, width := utf8.DecodeRuneInString(l.src)

	l.width = width
	l.pos += width

	if char == '\n' {
		l.line++
		l.lineStart = l.pos
	}

	return char
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
