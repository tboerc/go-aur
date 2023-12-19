package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name:    "go-aur",
		Version: "0.1.0",
		Usage:   "another tool to manage AUR packages",
		Authors: []*cli.Author{
			{Name: "tboerc", Email: "tiago.boer@proton.me"},
		},
		Commands: []*cli.Command{
			install, remove, update, list,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
