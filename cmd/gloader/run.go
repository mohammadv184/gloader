package main

import (
	"fmt"
	"log"
	"regexp"

	g "github.com/mohammadv184/gloader"
	"github.com/mohammadv184/gloader/driver"

	"github.com/spf13/cobra"
)

var runCmd = &cobra.Command{
	Use:   "run source destination [options]",
	Short: "run a migration",
	Long: `Migrate data from any source to any destination in a single command.
           e.g. gloader run mysql://root:root@localhost:3306/tests cockroach://root:root@localhost:5432/tests --filter version<3`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		gloader := g.NewGLoader()
		source := args[0]
		destination := args[1]
		fmt.Println("Migrating data from", source, "to", destination, "...")

		filters, _ := cmd.Flags().GetStringToString("filter")
		filtersAll, _ := cmd.Flags().GetString("filter-all")
		sorts, _ := cmd.Flags().GetStringToString("sort")
		sortsAll, _ := cmd.Flags().GetString("sort-all")
		reverseSorts, _ := cmd.Flags().GetStringToString("sort-reverse")
		reverseSortsAll, _ := cmd.Flags().GetString("sort-reverse-all")
		includes, _ := cmd.Flags().GetStringSlice("table")
		excludes, _ := cmd.Flags().GetStringSlice("exclude")
		startOffsets, _ := cmd.Flags().GetStringToInt64("start-offset")
		endOffsets, _ := cmd.Flags().GetStringToInt64("end-offset")
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

		if filters != nil {
			for dc, filter := range filters {
				r := regexp.MustCompile(`^([^<>=]+)([<>=]+)(.*)$`)
				filterKey := r.FindStringSubmatch(filter)[1]
				filterOperator := r.FindStringSubmatch(filter)[2]
				filterValue := r.FindStringSubmatch(filter)[3]
				filterC := driver.GetConditionFromString(filterOperator)
				gloader.Filter(dc, filterKey, filterC, filterValue)
			}
		}
		if filtersAll != "" {
			r := regexp.MustCompile(`^([^<>=]+)([<>=]+)(.*)$`)
			filterKey := r.FindStringSubmatch(filtersAll)[1]
			filterOperator := r.FindStringSubmatch(filtersAll)[2]
			filterValue := r.FindStringSubmatch(filtersAll)[3]
			filterC := driver.GetConditionFromString(filterOperator)
			gloader.FilterAll(filterKey, filterC, filterValue)
		}

		if sorts != nil {
			for dc, sort := range sorts {
				gloader.OrderBy(dc, sort, driver.Asc)
			}
		}
		if sortsAll != "" {
			gloader.OrderByAll(sortsAll, driver.Asc)
		}
		if reverseSorts != nil {
			for dc, sort := range reverseSorts {
				gloader.OrderBy(dc, sort, driver.Desc)
			}
		}
		if reverseSortsAll != "" {
			gloader.OrderByAll(reverseSortsAll, driver.Desc)
		}
		if includes != nil {
			gloader.Include(includes...)
		}
		if excludes != nil {
			gloader.Exclude(excludes...)
		}
		if startOffsets != nil {
			for dc, offset := range startOffsets {
				gloader.SetStartOffset(dc, uint64(offset))
			}
		}
		if endOffsets != nil {
			for dc, offset := range endOffsets {
				gloader.SetEndOffset(dc, uint64(offset))
			}
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

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringToStringP("filter", "f", nil, "filter data to migrate")
	runCmd.Flags().String("filter-all", "", "filter data to migrate (all tables)")
	runCmd.Flags().StringToStringP("sort", "s", nil, "sort data to migrate in ascending order")
	runCmd.Flags().String("sort-all", "", "sort data to migrate in ascending order (all tables)")
	runCmd.Flags().StringToStringP("sort-reverse", "S", nil, "sort data to migrate in descending order")
	runCmd.Flags().String("sort-reverse-all", "", "sort data to migrate in descending order (all tables)")
	runCmd.Flags().StringSliceP("table", "t", nil, "migrate only these tables")
	runCmd.Flags().StringSliceP("exclude", "e", nil, "exclude tables from migration")
	runCmd.Flags().StringToInt64("start-offset", nil, "start offset for each table")
	runCmd.Flags().StringToInt64("end-offset", nil, "end offset for each table")
	runCmd.Flags().Uint64P("rows-per-batch", "r", g.DefaultRowsPerBatch, "number of rows per batch")
	runCmd.Flags().UintP("workers", "w", g.DefaultWorkers, "number of workers")
}
