package main

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

var (
	args config
	wg   sync.WaitGroup
)

func main() {
	arg.MustParse(&args)

	for key, provider := range providers {
		if args.All || key == args.Provider {
			for _, region := range provider.regions {
				if (len(args.FilterRegion) > 0 && !isFiltered(region.name, args.FilterRegion)) ||
					(len(args.FilterLocation) > 0 && !isFiltered(region.location, args.FilterLocation)) {
					continue
				}
				wg.Add(1)
				go endpointPing(provider.hostTemplate, region)
			}
		}
	}
	wg.Wait()

	results := []*region{}
	for key, provider := range providers {
		if args.All || key == args.Provider {
			for _, region := range provider.regions {
				region.code = provider.name
				results = append(results, region)
			}
		}
	}
	if len(results) == 0 {
		fmt.Printf("Provider %s not found\n", args.Provider)
		os.Exit(1)
	}
	sort.Slice(results, func(i, j int) bool {
		return results[i].rtt < results[j].rtt
	})

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.SetAutoIndex(true)
	t.AppendHeader(table.Row{"Cloud", "Region", "Location", "RTT", "Status"})
	var minRTT time.Duration
	for _, region := range results {
		status := "OK"
		if region.err != nil {
			if args.HideErrors {
				continue
			}
			status = fmt.Sprintf("ERROR: %s", region.err)
		} else if region.rtt == 0 {
			continue
		}
		rtt := region.rtt.Round(time.Microsecond)
		if rtt > time.Second {
			rtt = rtt.Round(time.Millisecond)
		}
		if rtt < minRTT || minRTT == 0 {
			minRTT = rtt
		}
		t.AppendRow(table.Row{region.code, region.name, region.location, rtt, status})
	}
	t.SetColumnConfigs([]table.ColumnConfig{
		{
			Name:  "RTT",
			Align: text.AlignRight,
			Transformer: func(val interface{}) string {
				if val.(time.Duration) > time.Second {
					return text.FgRed.Sprintf("%s", val)
				}
				if val.(time.Duration) == minRTT {
					return text.Bold.Sprintf("%s", val)
				}
				return val.(time.Duration).String()
			},
		},
		{
			Name: "Status",
			Transformer: text.Transformer(func(val interface{}) string {
				color := text.FgGreen
				str := val.(string)
				if strings.Contains(str, "ERROR") {
					color = text.FgRed
				}
				return color.Sprintf("%s", val)
			}),
			WidthMax: 55,
		},
	})
	if t.Length() > 0 {
		t.Render()
	} else {
		fmt.Println("No results found")
		os.Exit(2)
	}
}
