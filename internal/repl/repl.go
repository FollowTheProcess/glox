// Package repl implements the read, eval, print loop (REPL) for glox.
package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/FollowTheProcess/glox/internal/syntax/parser"
)

const prompt = "-> "

// Start starts the REPL, reading from in and printing to out.
func Start(in io.Reader, out io.Writer) error {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, prompt)

		scanned := scanner.Scan()
		if !scanned {
			return scanner.Err()
		}

		line := scanner.Text()
		p := parser.New("REPL", line)

		program, err := p.Parse()
		if err != nil {
			return err
		}

		fmt.Fprintf(out, "%#v\n", program)
	}
}
