// Package lexer implements the lexical scanner for glox.
package lexer

import (
	"errors"
	"unicode"
	"unicode/utf8"

	"github.com/FollowTheProcess/glox/internal/syntax/token"
)

// lexFn represents the state of the scanner as a function that returns the next state.
type lexFn func(*Lexer) lexFn

const eof = rune(0)

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

// peek returns, but does not consume, the next rune in the input.
func (l *Lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// current returns the rune the lexer is currently sat on.
func (l *Lexer) current() rune {
	if l.atEOF() {
		return eof
	}
	return rune(l.src[l.pos])
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
	if l.width == 1 && l.current() == '\n' {
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
	src := l.src[l.start:l.pos]
	width := len(src)
	l.tokens <- token.Token{
		Text:   src,
		Kind:   kind,
		Offset: l.pos - width, // The start position of this token
		Width:  width,
	}
	l.start = l.pos
}

// run starts the state machine for the lexer.
func (l *Lexer) run() {
	for state := lexStart; state != nil; {
		state = state(l)
	}
	l.tokens <- token.Token{Kind: token.EOF, Offset: l.pos}
	close(l.tokens)
}

// lexStart is the initial state of the lexer.
func lexStart(l *Lexer) lexFn { //nolint: cyclop // Technically yes, but it's trivial really
	l.skipWhitespace()

	char := l.current()
	switch {
	case char == '(':
		return lexOpenParen
	case char == ')':
		return lexCloseParen
	case char == '{':
		return lexOpenBrace
	case char == '}':
		return lexCloseBrace
	case char == ',':
		return lexComma
	case char == '.':
		return lexDot
	case char == '-':
		return lexMinus
	case char == '+':
		return lexPlus
	case char == ';':
		return lexSemiColon
	case char == '/':
		return lexForwardSlash
	case char == '*':
		return lexStar
	case char == '!':
		return lexBang
	case char == '=':
		return lexEqual
	case char == '>':
		return lexGreaterThan
	case char == '<':
		return lexLessThan
	case char == '"':
		return lexString
	case unicode.IsDigit(char):
		return lexNumber
	case isAlpha(char):
		return lexIdent
	case char == eof:
		return nil
	default:
		return lexUnexpectedChar
	}
}

// lexOpenParen scans a '(' char.
func lexOpenParen(l *Lexer) lexFn {
	l.pos++
	l.emit(token.OpenParen)
	return lexStart
}

// lexCloseParen scans a ')' char.
func lexCloseParen(l *Lexer) lexFn {
	l.pos++
	l.emit(token.CloseParen)
	return lexStart
}

// lexOpenBrace scans a '{' char.
func lexOpenBrace(l *Lexer) lexFn {
	l.pos++
	l.emit(token.OpenBrace)
	return lexStart
}

// lexCloseBrace scans a '}' char.
func lexCloseBrace(l *Lexer) lexFn {
	l.pos++
	l.emit(token.CloseBrace)
	return lexStart
}

// lexComma scans a ',' char.
func lexComma(l *Lexer) lexFn {
	l.pos++
	l.emit(token.Comma)
	return lexStart
}

// lexDot scans a '.' char.
func lexDot(l *Lexer) lexFn {
	l.pos++
	l.emit(token.Dot)
	return lexStart
}

// lexMinus scans a '-' char.
func lexMinus(l *Lexer) lexFn {
	l.pos++
	l.emit(token.Minus)
	return lexStart
}

// lexPlus scans a '+' char.
func lexPlus(l *Lexer) lexFn {
	l.pos++
	l.emit(token.Plus)
	return lexStart
}

// lexSemiColon scans a ';' char.
func lexSemiColon(l *Lexer) lexFn {
	l.pos++
	l.emit(token.SemiColon)
	return lexStart
}

// lexForwardSlash scans a '/' char.
func lexForwardSlash(l *Lexer) lexFn {
	l.pos++ // Consume the first '/'
	if l.peek() == '/' {
		// It's a comment, absorb the whole line
		for l.peek() != '\n' && !l.atEOF() {
			l.next()
		}
		return lexStart
	}
	l.emit(token.ForwardSlash)
	return lexStart
}

// lexStar scans a '*' char.
func lexStar(l *Lexer) lexFn {
	l.pos++
	l.emit(token.Star)
	return lexStart
}

// lexBang scans a '!' char.
func lexBang(l *Lexer) lexFn {
	l.pos++ // Consume the '!'
	if l.peek() == '=' {
		return lexBangEqual
	}
	l.emit(token.Bang)
	return lexStart
}

// lexBangEqual scans the '!=' lexeme.
//
// The '!' has already been consumed by lexBang.
func lexBangEqual(l *Lexer) lexFn {
	l.pos++ // Consume the remaining '='
	l.emit(token.BangEqual)
	return lexStart
}

// lexEqual scans a '=' char.
func lexEqual(l *Lexer) lexFn {
	l.pos++ // Consume the '='
	if l.peek() == '=' {
		// Its a '=='
		return lexDoubleEqual
	}
	l.emit(token.Equal)
	return lexStart
}

// lexDoubleEqual scans the '==' lexeme.
//
// The first '=' has already been consumed.
func lexDoubleEqual(l *Lexer) lexFn {
	l.pos++ // Consume the remaining '='
	l.emit(token.DoubleEqual)
	return lexStart
}

// lexGreaterThan scans the '>' char.
func lexGreaterThan(l *Lexer) lexFn {
	l.pos++ // Consume the '>'
	if l.peek() == '=' {
		// Its a '>='
		return lexGreaterThanEqual
	}
	l.emit(token.GreaterThan)
	return lexStart
}

// lexGreaterThanEqual scans the '>=' lexeme.
//
// The initial '>' has already been consumed.
func lexGreaterThanEqual(l *Lexer) lexFn {
	l.pos++ // Consume the remaining '='
	l.emit(token.GreaterThanEqual)
	return lexStart
}

// lexLessThan scans the '>' char.
func lexLessThan(l *Lexer) lexFn {
	l.pos++ // Consume the '<'
	if l.peek() == '=' {
		// Its a '<='
		return lexLessThanEqual
	}
	l.emit(token.LessThan)
	return lexStart
}

// lexLessThanEqual scans the '<=' lexeme.
//
// The initial '<' has already been consumed.
func lexLessThanEqual(l *Lexer) lexFn {
	l.pos++ // Consume the remaining '='
	l.emit(token.LessThanEqual)
	return lexStart
}

// lexString scans a quoted string literal.
func lexString(l *Lexer) lexFn {
	l.pos++ // Consume the opening '"'
	for l.peek() != '"' && !l.atEOF() {
		l.next() // Consume everything until the next quote
	}

	if l.atEOF() {
		return l.error(errors.New("unterminated string literal"))
	}

	l.pos++ // Consume the closing '"'

	// TODO(@FollowTheProcess): Do we need to strip the quotes off?
	l.emit(token.String)
	return lexStart
}

// lexNumber scans a number literal.
func lexNumber(l *Lexer) lexFn {
	for unicode.IsDigit(l.peek()) {
		l.next()

		if l.next() == '.' && unicode.IsDigit(l.peek()) {
			l.next() // Consume the '.'
			for unicode.IsDigit(l.peek()) {
				l.next() // Absorb digits
			}
		}
		l.backup() // Undo the l.next() call in the if above
	}

	l.emit(token.Number)
	return lexStart
}

// lexIdent scans an identifier.
func lexIdent(l *Lexer) lexFn {
	l.pos++ // Absorb the first ident char
	for isAlphaNumeric(l.peek()) {
		l.next() // Absorb any alphanumeric characters
	}

	l.emit(token.Ident)
	return lexStart
}

// error emits an error token and terminates the scan by returning nil, halting
// the state machine in l.run().
func (l *Lexer) error(err error) lexFn {
	l.tokens <- token.Token{
		Text:   []byte(err.Error()),
		Kind:   token.Error,
		Offset: l.pos,
	}
	return nil
}

// lexUnexpectedChar handles any unrecognised char in the input by
// emitting an error token with the information and returning nil
// to halt the state machine.
func lexUnexpectedChar(l *Lexer) lexFn {
	cur := string(l.current())
	l.tokens <- token.Token{
		Text:   []byte("unexpected char " + cur),
		Kind:   token.Error,
		Offset: l.pos,
	}
	return nil
}

// isAlpha reports whether the rune is a valid identifier.
func isAlpha(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r == '_')
}

// isAlphaNumeric reports whether the rune is an alpha numeric character.
func isAlphaNumeric(r rune) bool {
	return isAlpha(r) || unicode.IsDigit(r)
}
