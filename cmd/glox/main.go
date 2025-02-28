package main

import (
	"os"

	"github.com/FollowTheProcess/glox/internal/cmd"
	"github.com/FollowTheProcess/msg"
)

func main() {
	if err := run(); err != nil {
		msg.Err(err)
		os.Exit(1)
	}
}

func run() error {
	cli, err := cmd.Build()
	if err != nil {
		return err
	}

	return cli.Execute()
}
