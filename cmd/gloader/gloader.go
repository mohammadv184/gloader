package main

import (
	"os"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gloader command [args] [flags]",
	Short: "Migrate data from any source to any destination",
	Long: `Migrate data from any source to any destination in a single command.
           e.g. gloader run mysql://root:root@localhost:3306/tests cockroach://root:root@localhost:5432/tests --filter version<3`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func main() {
	timeTracker := time.Now()
	Execute()
	println(time.Since(timeTracker).String())
}
