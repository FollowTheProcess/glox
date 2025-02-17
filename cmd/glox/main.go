package main

import (
	"errors"
	"os"

	"github.com/FollowTheProcess/cli"
	"github.com/FollowTheProcess/glox/internal/repl"
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
		cli.Run(func(cmd *cli.Command, args []string) error {
			if len(args) != 0 {
				return errors.New("only the REPL is supported at the moment")
			}

			return repl.Start(os.Stdin, os.Stdout)
		}),
	)
	if err != nil {
		return err
	}

	return cmd.Execute()
}
