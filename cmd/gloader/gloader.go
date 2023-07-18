package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "gloader command [args] [flags]",
	Version: fmt.Sprintf("gloader version %s %s %s", version, commit, date),
	Short:   "Migrate data from any source to any destination",
	Long: `Migrate data from any source to any destination in a single command.
           e.g. gloader run mysql://root:root@localhost:3306/tests cockroach://root:root@localhost:5432/tests --filter version<3`,
}

func main() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered Err:", r)
		}
	}()
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
