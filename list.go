package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/urfave/cli/v2"
)

var list = &cli.Command{
	Name:    "list",
	Aliases: []string{"ls"},
	Action: func(cCtx *cli.Context) error {
		aur := &aur{}
		if err := aur.get(); err != nil {
			return err
		}

		w := tabwriter.NewWriter(os.Stdout, 1, 1, 5, ' ', 0)
		defer w.Flush()

		fmt.Fprintln(w, "Package\tInstalled version\t")

		for _, entry := range aur.entries {
			if err := entry.get(); err != nil {
				return err
			}

			fmt.Fprintf(w, "%v\t%v\t\n", entry.pkg.name(), entry.installver)
		}

		return nil
	},
}
