// Package repl implements the read, eval, print loop (REPL) for glox.
package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/FollowTheProcess/glox/internal/eval"
	"github.com/FollowTheProcess/glox/internal/syntax/parser"
)

const prompt = "-> "

// Start starts the REPL, reading from in and printing to out.
func Start(in io.Reader, out io.Writer, trace bool) error {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, prompt)

		scanned := scanner.Scan()
		if !scanned {
			return scanner.Err()
		}

		line := scanner.Text()
		p := parser.New("REPL", line, trace, os.Stderr)

		program, err := p.Parse()
		if err != nil {
			return err
		}

		result, err := eval.Eval(program)
		if err != nil {
			return err
		}

		fmt.Fprintf(out, "%s\n", result)
	}
}
