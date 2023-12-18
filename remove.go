package main

import (
	"github.com/urfave/cli/v2"
)

var remove = &cli.Command{
	Name:    "remove",
	Aliases: []string{"rm"},
}
