// Package repl implements the read, eval, print loop (REPL) for glox.
package repl

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/FollowTheProcess/glox/internal/interpreter"
	"github.com/FollowTheProcess/glox/internal/syntax/parser"
)

const prompt = "-> "

// Start starts the REPL, reading from in and printing to out.
func Start(in io.Reader, out io.Writer, trace bool) error {
	scanner := bufio.NewScanner(in)
	interp := interpreter.New(out, os.Stderr) // Must be outside so state is preserved

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

		result, err := interp.Eval(program)
		if err != nil {
			return err
		}

		if result != nil {
			fmt.Fprintf(out, "%s\n", result)
		}
	}
}
