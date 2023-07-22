package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"sync"
	"syscall"
	"time"

	g "github.com/mohammadv184/gloader"
	"github.com/mohammadv184/gloader/driver"
	"github.com/spf13/cobra"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"golang.org/x/term"
)

var (
	flagFilterAll      []string
	flagSortAll        []string
	flagReverseSortAll []string
	flagTable          []string
	flagExclude        []string
	flagFilter         StringToStringSliceFlag
	flagSort           StringToStringSliceFlag
	flagReverseSort    StringToStringSliceFlag
	flagStartOffset    map[string]int64
	flagEndOffset      map[string]int64
	flagRowsPerBatch   uint64
	flagWorkers        uint
)

var runCmd = &cobra.Command{
	Use:   "run <source> <destination> [options]",
	Short: "run a migration",
	Long: `Migrate data from any source to any destination in a single command.
           e.g. gloader run mysql://root:root@localhost:3306/tests cockroach://root:root@localhost:5432/tests --filter version<3`,
	Args: cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		gloader := g.NewGLoader()
		source := args[0]
		destination := args[1]
		fmt.Println("Migrating data from", source, "to", destination, "...")

		sourceDriver := regexp.MustCompile(`^([a-z]+)://`).FindStringSubmatch(source)[1]
		destinationDriver := regexp.MustCompile(`^([a-z]+)://`).FindStringSubmatch(destination)[1]

		sourceDSN := regexp.MustCompile(`^[a-z]+://(.*)`).FindStringSubmatch(source)[1]
		destinationDSN := regexp.MustCompile(`^[a-z]+://(.*)`).FindStringSubmatch(destination)[1]

		err := gloader.Src(sourceDriver, sourceDSN)
		if err != nil {
			log.Fatal(err)
		}
		err = gloader.Dest(destinationDriver, destinationDSN)
		if err != nil {
			log.Fatal(err)
		}

		if flagFilter.Length() > 0 {
			for dc, filters := range flagFilter.Value() {
				for _, filter := range filters {
					r := regexp.MustCompile(`^([^<>=]+)([<>=]+)(.*)$`)
					filterKey := r.FindStringSubmatch(filter)[1]
					filterOperator := r.FindStringSubmatch(filter)[2]
					filterValue := r.FindStringSubmatch(filter)[3]
					filterC := driver.GetConditionFromString(filterOperator)
					gloader.Filter(dc, filterKey, filterC, filterValue)
				}
			}
		}
		if len(flagFilterAll) > 0 {
			for _, filtersAll := range flagFilterAll {
				r := regexp.MustCompile(`^([^<>=]+)([<>=]+)(.*)$`)
				filterKey := r.FindStringSubmatch(filtersAll)[1]
				filterOperator := r.FindStringSubmatch(filtersAll)[2]
				filterValue := r.FindStringSubmatch(filtersAll)[3]
				filterC := driver.GetConditionFromString(filterOperator)
				gloader.FilterAll(filterKey, filterC, filterValue)
			}
		}

		if flagSort.Length() > 0 {
			for dc, sorts := range flagSort.Value() {
				for _, sort := range sorts {
					gloader.OrderBy(dc, sort, driver.Asc)
				}
			}
		}

		if len(flagSortAll) > 0 {
			for _, sortsAll := range flagSortAll {
				gloader.OrderByAll(sortsAll, driver.Asc)
			}
		}

		if flagReverseSort.Length() > 0 {
			for dc, sorts := range flagReverseSort.Value() {
				for _, sort := range sorts {
					gloader.OrderBy(dc, sort, driver.Desc)
				}
			}
		}

		if len(flagReverseSortAll) > 0 {
			for _, sortsAll := range flagReverseSortAll {
				gloader.OrderByAll(sortsAll, driver.Desc)
			}
		}

		if len(flagTable) > 0 {
			gloader.Include(flagTable...)
		}

		if len(flagExclude) > 0 {
			gloader.Exclude(flagExclude...)
		}

		if len(flagStartOffset) > 0 {
			for dc, offset := range flagStartOffset {
				gloader.SetStartOffset(dc, uint64(offset))
			}
		}

		if len(flagEndOffset) > 0 {
			for dc, offset := range flagEndOffset {
				gloader.SetEndOffset(dc, uint64(offset))
			}
		}

		if flagRowsPerBatch != 0 {
			gloader.SetRowsPerBatch(flagRowsPerBatch)
		}

		if flagWorkers != 0 {
			gloader.SetWorkers(flagWorkers)
		}

		wg := &sync.WaitGroup{}

		ctx, cancelFunc := context.WithCancelCause(context.Background())

		srcDetails, err := gloader.GetSrcDetails(ctx)
		if err != nil {
			log.Fatal(err)
		}

		dataCollections := srcDetails.DataCollections
		if len(flagTable) > 0 {
			dataCollections = srcDetails.OnlyDataCollections(flagTable...)
		}
		if len(flagExclude) > 0 {
			dataCollections = srcDetails.AllDataCollectionsExcept(flagExclude...)
		}

		gStats := gloader.Stats()

		w := 80
		if term.IsTerminal(int(os.Stdout.Fd())) {
			width, _, err := term.GetSize(int(os.Stdout.Fd()))
			if err != nil {
				log.Println(err)
			} else {
				w = width
			}
		}

		pbars := mpb.New(
			mpb.WithWidth(w),
			mpb.WithContext(ctx),
			mpb.WithWaitGroup(wg),
		)

		for i, dc := range dataCollections {
			var maxDCNameWidth int
			for _, dc := range dataCollections {
				maxDCNameWidth = int(math.Max(float64(maxDCNameWidth), float64(len(dc.Name))))
			}

			b := pbars.AddBar(
				int64(dc.DataSetCount),
				mpb.BarPriority(i),
				mpb.PrependDecorators(
					decor.Name(dc.Name, decor.WC{W: maxDCNameWidth, C: decor.DidentRight}),
					decor.CountersNoUnit("%d Rows / %d Rows", decor.WCSyncSpace),
				),
				mpb.AppendDecorators(
					decor.Percentage(decor.WC{W: 6, C: decor.DidentRight}),
					decor.EwmaETA(decor.ET_STYLE_HHMMSS, float64(flagRowsPerBatch)*float64(flagWorkers), decor.WCSyncSpaceR),
					decor.AverageSpeed(0, "%.0f Rows/s", decor.WCSyncWidth),
				),
			)

			go func(b *mpb.Bar, dc driver.DataCollectionDetail) {
				m := gStats.MustGetSequentialCounter(g.MetricBufferTotalReadLengthRows.String())

				mChangeNotifier := make(chan any)
				m.NotifyOnChange(mChangeNotifier, dc.Name)

				for {
					lastReportT := time.Now()

					select {
					case <-ctx.Done():
						if !b.Completed() {
							wg.Done()
						}
						return
					case <-mChangeNotifier:
						b.IncrBy(int(m.Value(dc.Name)-b.Current()), time.Since(lastReportT))
					}
				}
			}(b, dc)

		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			err = gloader.StartWithContext(ctx)
			if err != nil {
				log.Fatal(err)
			}
			cancelFunc(errors.New("done"))
		}()

		closeSignal := make(chan os.Signal, 1)
		signal.Notify(closeSignal, syscall.SIGTERM, syscall.SIGINT, os.Interrupt)
		select {
		case <-closeSignal:
			fmt.Println("close signal received")
			cancelFunc(errors.New("close signal received"))
		case <-ctx.Done():
		}

		wg.Wait()
		fmt.Println("Done.")
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	runCmd.Flags().VarP(&flagFilter, "filter", "f", "filter data to migrate")
	runCmd.Flags().StringSliceVar(&flagFilterAll, "filter-all", nil, "filter data to migrate (all tables)")
	runCmd.Flags().VarP(&flagSort, "sort", "s", "sort data to migrate in ascending order")
	runCmd.Flags().StringSliceVar(&flagSortAll, "sort-all", nil, "sort data to migrate in ascending order (all tables)")
	runCmd.Flags().VarP(&flagReverseSort, "sort-reverse", "S", "sort data to migrate in descending order")
	runCmd.Flags().StringSliceVar(&flagReverseSortAll, "sort-reverse-all", nil, "sort data to migrate in descending order (all tables)")
	runCmd.Flags().StringSliceVarP(&flagTable, "table", "t", nil, "migrate only these tables")
	runCmd.Flags().StringSliceVarP(&flagExclude, "exclude", "e", nil, "exclude tables from migration")
	runCmd.Flags().StringToInt64Var(&flagStartOffset, "start-offset", nil, "start offset for each table")
	runCmd.Flags().StringToInt64Var(&flagEndOffset, "end-offset", nil, "end offset for each table")
	runCmd.Flags().Uint64VarP(&flagRowsPerBatch, "rows-per-batch", "r", g.DefaultRowsPerBatch, "number of rows per batch")
	runCmd.Flags().UintVarP(&flagWorkers, "workers", "w", g.DefaultWorkers, "number of workers")
}

// StringToStringSliceFlag is a custom flag type
// for example: --filter=clients=["id > 5","name = 'John'"],orders=["id > 10"].
type StringToStringSliceFlag struct {
	value map[string][]string
}

func (f *StringToStringSliceFlag) String() string {
	var s strings.Builder
	for k, v := range f.value {
		s.WriteString(fmt.Sprintf("%s=[%s],", k, strings.Join(v, ",")))
	}
	return strings.TrimSuffix(s.String(), ",")
}

func (f *StringToStringSliceFlag) Set(value string) error {
	if f.value == nil {
		f.value = make(map[string][]string)
	}

	// Regular expression pattern to match key-value pairs
	pattern := `([^=]+)=\[(.*?)\](?:,|$)`

	// Find all matches of key-value pairs in the input string
	matches := regexp.MustCompile(pattern).FindAllStringSubmatch(value, -1)

	for _, match := range matches {
		if len(match) != 3 {
			return fmt.Errorf("invalid key-value pair: %s", match[0])
		}

		key := strings.TrimSpace(match[1])

		// Split the values within the square brackets by comma
		values := strings.Split(match[2], ",")

		// Trim any leading/trailing spaces from each value
		for i := 0; i < len(values); i++ {
			values[i] = strings.TrimSpace(values[i])
			values[i] = strings.Trim(values[i], `"`)
		}

		// Append the values to the existing slice for the corresponding key
		f.value[key] = append(f.value[key], values...)
	}

	return nil
}

func (f *StringToStringSliceFlag) Type() string {
	return "stringToStringSlice"
}

func (f *StringToStringSliceFlag) Value() map[string][]string {
	return f.value
}

func (f *StringToStringSliceFlag) GetValuesOf(key string) []string {
	return f.value[key]
}

func (f *StringToStringSliceFlag) GetKeys() []string {
	keys := make([]string, 0, len(f.value))
	for k := range f.value {
		keys = append(keys, k)
	}
	return keys
}

func (f *StringToStringSliceFlag) HasKey(key string) bool {
	_, ok := f.value[key]
	return ok
}

func (f *StringToStringSliceFlag) Length() int {
	return len(f.value)
}
