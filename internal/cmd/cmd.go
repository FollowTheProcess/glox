// Package cmd implements the CLI for glox.
package cmd

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/FollowTheProcess/cli"
	"github.com/FollowTheProcess/glox/internal/repl"
	"github.com/FollowTheProcess/glox/internal/syntax/lexer"
	"github.com/FollowTheProcess/glox/internal/syntax/parser"
	"github.com/FollowTheProcess/glox/internal/syntax/token"
)

var (
	version = "dev"
	commit  = ""
	date    = ""
)

// Build builds and returns the glox CLI.
func Build() (*cli.Command, error) {
	return cli.New(
		"glox",
		cli.Short("Go implementation of the Lox language from Crafting Interpreters"),
		cli.Version(version),
		cli.Commit(commit),
		cli.BuildDate(date),
		cli.Allow(cli.NoArgs()),
		cli.SubCommands(runCommand, replCommand),
	)
}

type runOptions struct {
	tokenise bool // Tokenise the source file only, emitting the stream of tokens
	debug    bool // Emit debug information
	timings  bool // Emit performance information
}

// runCommand returns the run subcommand.
func runCommand() (*cli.Command, error) {
	var options runOptions
	return cli.New(
		"run",
		cli.Short("Run the Lox interpreter on a source file"),
		cli.RequiredArg("src", "Lox source file"),
		cli.Flag(
			&options.tokenise,
			"tokenise",
			cli.NoShortHand,
			false,
			"Tokenise the source file only, emitting the stream of tokens",
		),
		cli.Flag(&options.debug, "debug", cli.NoShortHand, false, "Emit debug information"),
		cli.Flag(&options.timings, "timings", cli.NoShortHand, false, "Emit performance information"),
		cli.Run(doRun(&options)),
	)
}

// replCommand returns the repl subcommand.
func replCommand() (*cli.Command, error) {
	return cli.New(
		"repl",
		cli.Short("Start an interactive REPL for Lox"),
		cli.Allow(cli.NoArgs()),
		cli.Run(func(cmd *cli.Command, args []string) error {
			return repl.Start(os.Stdin, os.Stdout)
		}),
	)
}

// doRun actually runs glox.
func doRun(options *runOptions) func(cmd *cli.Command, args []string) error {
	return func(cmd *cli.Command, args []string) error {
		start := time.Now()
		defer func() {
			fmt.Fprintf(cmd.Stderr(), "\nTook %v\n", time.Since(start))
		}()

		src := cmd.Arg("src")
		contents, err := os.ReadFile(src)
		if err != nil {
			return err
		}

		// Tokenise only
		if options.tokenise {
			stdoutWriter := bufio.NewWriter(cmd.Stdout())
			l := lexer.New(string(contents))
			for {
				tok := l.NextToken()
				fmt.Fprintln(stdoutWriter, tok.String())
				if tok.Is(token.EOF) {
					break
				}
			}
			return stdoutWriter.Flush()
		}

		parser := parser.New(src, string(contents))
		program, err := parser.Parse()
		if err != nil {
			return err
		}

		fmt.Fprintf(cmd.Stdout(), "%#v\n", program)
		return nil
	}
}
