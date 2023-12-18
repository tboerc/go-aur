package main

import (
	"fmt"
	"os"
	"text/tabwriter"

	"github.com/urfave/cli/v2"
)

var update = &cli.Command{
	Name:    "update",
	Aliases: []string{"up"},
	Action: func(cCtx *cli.Context) error {
		aur := &aur{}
		if err := aur.get(); err != nil {
			return err
		}

		outdated := []*entry{}

		for _, entry := range aur.entries {
			if err := entry.pull(); err != nil {
				return err
			}

			if err := entry.get(); err != nil {
				return err
			}

			if entry.installver != entry.pkg.version() {
				outdated = append(outdated, entry)
			}
		}

		if len(outdated) == 0 {
			fmt.Println("\nEverything is up-to-date")
			return nil
		}

		w := tabwriter.NewWriter(os.Stdout, 1, 1, 5, ' ', 0)

		fmt.Fprintln(w, "\nPackage\tInstalled version\tNew version\t")

		for _, entry := range outdated {
			fmt.Fprintf(w, "%v\t%v\t%v\t\n", entry.pkg.name(), entry.installver, entry.pkg.version())
		}

		w.Flush()

		for _, entry := range outdated {
			entry.update()
		}

		return nil
	},
}
