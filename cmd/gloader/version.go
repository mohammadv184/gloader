package main

import "github.com/spf13/cobra"

// it's populated during build time by -ldflags.
var (
	version = "dev"
	commit  = "?"
	date    = ""
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print the version number of gloader",
	Long: `Print full version information of gloader
            output format: gloader version <version> <commit> <date>`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Println("gloader version", version, commit, date)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
