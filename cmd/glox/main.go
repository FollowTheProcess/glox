package main

import (
	"fmt"
	"os"
	"time"

	"github.com/FollowTheProcess/cli"
	"github.com/FollowTheProcess/glox/internal/repl"
	"github.com/FollowTheProcess/glox/internal/syntax/lexer"
	"github.com/FollowTheProcess/glox/internal/syntax/token"
	"github.com/FollowTheProcess/msg"
)

const bytesPerMegaByte = 1024 * 1024

var (
	version = "dev"
	commit  = ""
	date    = ""
)

func main() {
	if err := run(); err != nil {
		msg.Err(err)
		os.Exit(1)
	}
}

func run() error {
	cmd, err := cli.New(
		"glox",
		cli.Short("Go implementation of the Lox language from Crafting Interpreters"),
		cli.Version(version),
		cli.Commit(commit),
		cli.BuildDate(date),
		cli.OptionalArg("file", "Path to a .lox source file to execute", "repl"),
		cli.Run(func(cmd *cli.Command, args []string) error {
			file := cmd.Arg("file")
			if file == "repl" {
				return repl.Start(os.Stdin, os.Stdout)
			}

			contents, err := os.ReadFile(file)
			if err != nil {
				return err
			}

			src := string(contents)

			tokenCount := 0

			start := time.Now()
			lex := lexer.New(src)
			for tok := lex.NextToken(); tok.Kind != token.EOF; tok = lex.NextToken() {
				tokenCount++
				fmt.Fprintf(os.Stdout, "%s\n", tok)
			}

			duration := time.Since(start)
			size := len(contents)
			throughput := (float64(size) / duration.Seconds()) / bytesPerMegaByte

			fmt.Printf(
				"\nLexed %d tokens (%d bytes of source code) in %s (%.2f MB/s)",
				tokenCount,
				size,
				duration,
				throughput,
			)

			return nil
		}),
	)
	if err != nil {
		return err
	}

	return cmd.Execute()
}
