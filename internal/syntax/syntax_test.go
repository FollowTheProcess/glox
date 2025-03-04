package syntax_test

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/FollowTheProcess/glox/internal/syntax/lexer"
	"github.com/FollowTheProcess/glox/internal/syntax/parser"
	"github.com/FollowTheProcess/glox/internal/syntax/token"
	"github.com/FollowTheProcess/test"
	"github.com/FollowTheProcess/txtar"
	"github.com/kr/pretty"
)

var (
	debug  = flag.Bool("debug", false, "Emit parse traces during tests")
	update = flag.Bool("update", false, "Update snapshots and testdata")
)

func TestParseValid(t *testing.T) {
	validCaseGlob := filepath.Join("testdata", "valid", "*.txtar")

	validCases, err := filepath.Glob(validCaseGlob)
	test.Ok(t, err)

	for _, file := range validCases {
		name := filepath.Base(file)
		t.Run(name, func(t *testing.T) {
			archive, err := txtar.ParseFile(file)
			test.Ok(t, err)

			src, err := archive.Read("src.lox")
			test.Ok(t, err)

			expectedTokens, err := archive.Read("tokens.txt")
			test.Ok(t, err)

			expectedTree, err := archive.Read("parsed.txt")
			test.Ok(t, err)

			tokens := lex(src)
			tree := parse(t, name, src, *debug)

			if *update {
				// Update the expected with what's actually been seen
				err := archive.Write("tokens.txt", tokens)
				test.Ok(t, err)

				err = archive.Write("parsed.txt", tree)
				test.Ok(t, err)

				err = txtar.DumpFile(file, archive)
				test.Ok(t, err)
				return
			}

			// Test the tokens
			test.Diff(t, tokens, expectedTokens)

			// And the AST
			test.Diff(t, tree, expectedTree)
		})
	}
}

func TestParseInvalid(t *testing.T) {
	invalidCaseGlob := filepath.Join("testdata", "invalid", "*.txtar")

	invalidCases, err := filepath.Glob(invalidCaseGlob)
	test.Ok(t, err)

	for _, file := range invalidCases {
		name := filepath.Base(file)
		t.Run(name, func(t *testing.T) {
			archive, err := txtar.ParseFile(file)
			test.Ok(t, err)

			src, err := archive.Read("src.lox")
			test.Ok(t, err)

			want, err := archive.Read("want.txt")
			test.Ok(t, err)

			p := parser.New("src.lox", src, *debug, os.Stderr)
			_, err = p.Parse()
			test.Err(t, err, test.Context("invalid case must generate parse error"))

			if *update {
				err = archive.Write("want.txt", err.Error())
				test.Ok(t, err)

				err = txtar.DumpFile(file, archive)
				test.Ok(t, err)
				return
			}

			test.Diff(t, err.Error()+"\n", want)
		})
	}
}

// lex tokenises src, emitting a string of newline separated tokens.
func lex(src string) string {
	var tokens strings.Builder

	lexer := lexer.New(src)
	for {
		tok := lexer.NextToken()
		tokens.WriteString(tok.String() + "\n")
		if tok.Is(token.EOF) {
			break
		}
	}

	return tokens.String()
}

// parse parses a stream of tokens into an AST, returning it's formatted
// string representation.
func parse(t *testing.T, name, src string, debug bool) string {
	t.Helper()

	p := parser.New(name, src, debug, os.Stderr)
	program, err := p.Parse()
	test.Ok(t, err)

	return pretty.Sprintf("%# v\n", program)
}
