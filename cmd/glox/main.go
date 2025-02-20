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

			start := time.Now()
			contents, err := os.ReadFile(file)
			if err != nil {
				return err
			}

			tokenCount := 0

			lex := lexer.New(contents)
			for tok := lex.NextToken(); tok.Kind != token.EOF; tok = lex.NextToken() {
				tokenCount++
				fmt.Fprintf(os.Stdout, "%s\n", tok)
			}

			fmt.Printf("\nLexed %d tokens in %s", tokenCount, time.Since(start))

			return nil
		}),
	)
	if err != nil {
		return err
	}

	return cmd.Execute()
}
