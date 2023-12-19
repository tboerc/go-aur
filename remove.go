package main

import (
	"errors"

	"github.com/urfave/cli/v2"
)

var remove = &cli.Command{
	Name:    "remove",
	Aliases: []string{"rm"},
	Action: func(cCtx *cli.Context) error {
		pkg := cCtx.Args().Get(0)

		aur := &aur{}
		if err := aur.get(); err != nil {
			return err
		}

		entry := aur.entries[pkg]
		if entry == nil {
			return errors.New("package does not exist")
		}

		if err := entry.get(); err != nil {
			return err
		}

		if entry.installver == NOT_INSTALLED {
			return errors.New("package is not installed")
		}

		if err := entry.remove(); err != nil {
			return err
		}

		return nil
	},
}
