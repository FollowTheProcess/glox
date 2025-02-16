package main

import (
	"fmt"
	"os"

	"github.com/FollowTheProcess/cli"
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
			fmt.Printf("Lox, args: %v\n", args)
			return nil
		}),
	)
	if err != nil {
		return err
	}

	return cmd.Execute()
}
