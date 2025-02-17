// Package lexer implements the lexical scanner for glox.
package lexer

import (
	"fmt"
	"os"
	"unicode"
	"unicode/utf8"

	"github.com/FollowTheProcess/glox/internal/syntax/token"
)

// lexFn represents the state of the scanner as a function that returns the next state.
type lexFn func(*Lexer) lexFn

// Lexer is the lexical scanner.
type Lexer struct {
	tokens chan token.Token // Channel on which to emit lexed tokens
	src    []byte           // The src being scanned
	start  int              // Start position of the current token
	pos    int              // Current position in the input
	line   int              // Current line in the input
	width  int              // Width of the last rune read from input
}

// New creates a new lexer for the input string and sets it off in a goroutine.
func New(src []byte) *Lexer {
	l := &Lexer{
		tokens: make(chan token.Token),
		src:    src,
		start:  0,
		pos:    0,
		line:   1,
		width:  0,
	}
	go l.run()
	return l
}

// NextToken asks the lexer for the next token from the input.
func (l *Lexer) NextToken() token.Token {
	return <-l.tokens
}

// next returns, and consumes, the next rune in the input.
func (l *Lexer) next() rune {
	r, width := utf8.DecodeRune(l.rest())
	l.width = width
	l.pos += l.width
	if r == '\n' {
		l.line++
	}
	return r
}

// rest returns the string from the current lexer position to the end of the input.
func (l *Lexer) rest() []byte {
	if l.atEOF() {
		return nil
	}
	return l.src[l.pos:]
}

// atEOF returns whether or not the lexer is currently at the end of a file.
func (l *Lexer) atEOF() bool {
	return l.pos >= len(l.src)
}

// backup steps back one rune. Can only be called once per call of next.
func (l *Lexer) backup() {
	l.pos -= l.width
	if l.width == 1 && l.src[l.pos] == '\n' {
		l.line--
	}
}

// skipWhitespace consumes any utf-8 whitespace until something meaningful is hit.
func (l *Lexer) skipWhitespace() {
	for {
		r := l.next()
		if !unicode.IsSpace(r) {
			l.backup() // Go back to the last non-space
			l.start = l.pos
			break
		}
	}
}

// emit emits a token over the tokens channel.
func (l *Lexer) emit(kind token.Kind) {
	l.tokens <- token.Token{
		Text:   l.src[l.start:l.pos],
		Kind:   kind,
		Offset: l.pos,
	}
	l.start = l.pos
}

// run starts the state machine for the lexer.
func (l *Lexer) run() {
	for state := lexStart; state != nil; {
		state = state(l)
	}
	close(l.tokens)
}

// lexStart is the initial state of the lexer.
func lexStart(l *Lexer) lexFn {
	l.skipWhitespace()

	next := l.next()
	switch {
	case next == '(':
		return lexOpenParen
	case next == ')':
		return lexCloseParen
	case l.atEOF():
		l.tokens <- token.Token{Kind: token.EOF, Offset: l.pos}
		return nil
	default:
		fmt.Fprintf(os.Stderr, "unexpected char %q\n", next)
		return nil
	}
}

// lexOpenParen scans a '(' char.
func lexOpenParen(l *Lexer) lexFn {
	l.emit(token.OpenParen)
	return lexStart
}

// lexCloseParen scans a ')' char.
func lexCloseParen(l *Lexer) lexFn {
	l.emit(token.CloseParen)
	return lexStart
}
