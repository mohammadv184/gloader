package main

import (
	"fmt"
	g "gloader"
	"gloader/driver"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gloader source destination [flags]",
	Short: "Migrate data from any source to any destination",
	Long: `Migrate data from any source to any destination in a single command.
           e.g. gloader mysql://root:root@localhost:3306/test cockroach://root:root@localhost:5432/test --filter version<3`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		gloader := g.NewGLoader()
		source := args[0]
		destination := args[1]
		fmt.Println("Migrating data from", source, "to", destination, "...")

		filter, _ := cmd.Flags().GetString("filter")
		sort, _ := cmd.Flags().GetString("sort")
		sortReverse, _ := cmd.Flags().GetString("sort-reverse")
		rowsPerBatch, _ := cmd.Flags().GetUint64("rows-per-batch")
		workers, _ := cmd.Flags().GetUint("workers")

		sourceDriver := regexp.MustCompile(`^([a-z]+)://`).FindStringSubmatch(source)[1]
		destinationDriver := regexp.MustCompile(`^([a-z]+)://`).FindStringSubmatch(destination)[1]

		sourceDSN := regexp.MustCompile(`^[a-z]+://(.*)`).FindStringSubmatch(source)[1]
		destinationDSN := regexp.MustCompile(`^[a-z]+://(.*)`).FindStringSubmatch(destination)[1]

		err := gloader.Source(sourceDriver, sourceDSN)
		if err != nil {
			log.Fatal(err)
		}
		err = gloader.Dest(destinationDriver, destinationDSN)
		if err != nil {
			log.Fatal(err)
		}

		if filter != "" {
			r := regexp.MustCompile(`^([^<>=]+)([<>=]+)(.*)$`)
			filterKey := r.FindStringSubmatch(filter)[1]
			filterOperator := r.FindStringSubmatch(filter)[2]
			filterValue := r.FindStringSubmatch(filter)[3]
			filterC := driver.GetConditionFromString(filterOperator)
			gloader.Filter(filterKey, filterC, filterValue)
		}

		if sort != "" {
			gloader.OrderBy(sort, driver.Asc)
		}
		if sortReverse != "" {
			gloader.OrderBy(sortReverse, driver.Desc)
		}

		if rowsPerBatch != 0 {
			gloader.SetRowsPerBatch(rowsPerBatch)
		}

		if workers != 0 {
			gloader.SetWorkers(workers)
		}

		err = gloader.Start()
		if err != nil {
			log.Fatal(err)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().StringP("filter", "f", "", "filter data to migrate")
	rootCmd.Flags().StringP("sort", "s", "", "sort data to migrate in ascending order")
	rootCmd.Flags().StringP("sort-reverse", "S", "", "sort data to migrate in descending order")
	rootCmd.Flags().Uint64P("rows-per-batch", "r", 100, "number of rows per batch")
	rootCmd.Flags().UintP("workers", "w", 5, "number of workers")
}

func main() {
	timeTracker := time.Now()
	Execute()
	println(time.Since(timeTracker).String())

}
