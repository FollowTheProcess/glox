package parser

import (
	"fmt"

	"github.com/FollowTheProcess/glox/internal/syntax/token"
)

// SyntaxError is a syntax error in Lox.
type SyntaxError struct {
	File  string      // The name of the source file or "stdin" if REPL
	Msg   string      // The error message
	Token token.Token // The offending token triggering the error
	Line  int         // Line number on which Token appears
	Col   int         // Column number on Line on which the first char of Token appears
}

// Error implements the error interface for a SyntaxError, enabling its use
// as a standard Go error.
func (s SyntaxError) Error() string {
	return fmt.Sprintf("%s:%d:%d: %s", s.File, s.Line, s.Col, s.Msg)
}
