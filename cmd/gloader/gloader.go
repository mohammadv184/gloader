package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gloader command [args] [flags]",
	Short: "Migrate data from any source to any destination",
	Long: `GLoader is a CLI tool for data migration between different databases. 
              It allows you to migrate your data from any source database to any destination database in a single command.`,
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
