package main

import (
	"github.com/urfave/cli/v2"
)

var install = &cli.Command{
	Name:    "install",
	Aliases: []string{"it"},
	Action: func(cCtx *cli.Context) error {
		url := cCtx.Args().Get(0)

		aur := &aur{}
		if err := aur.get(); err != nil {
			return err
		}

		repo, err := aur.clone(url)
		if err != nil {
			return err
		}

		if err := aur.get(); err != nil {
			return err
		}

		entry := aur.entries[repo]
		if err := entry.get(); err != nil {
			return err
		}

		if err := entry.install(); err != nil {
			return err
		}

		return nil
	},
}
