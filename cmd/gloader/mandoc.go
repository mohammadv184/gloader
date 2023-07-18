package main

import (
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

var manDocCmd = &cobra.Command{
	Use:    "mandoc",
	Hidden: true,
	Run: func(cmd *cobra.Command, args []string) {
		rootc := cmd.Root()

		d := time.Now()
		if date != "" {
			var err error
			d, err = time.Parse(time.RFC3339, date)
			if err != nil {
				cmd.Println("failed to parse date:", err)
				return
			}
		}

		docH := &doc.GenManHeader{
			Title:   "GLoader",
			Section: "1",
			Date:    &d,
			Source:  "Database Migration Tool",
			Manual:  "GLoader Manual",
		}

		err := doc.GenMan(rootc, docH, cmd.OutOrStdout())
		if err != nil {
			cmd.Println("failed to generate man doc:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(manDocCmd)
}
