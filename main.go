package main

import (
	"fmt"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

var args config

func main() {
	arg.MustParse(&args)
	args.Provider = strings.ToLower(args.Provider)

	semaphore := make(chan struct{}, args.Concurrent)
	var wg sync.WaitGroup

	for key, provider := range providers {
		if args.All || key == args.Provider {
			for _, r := range provider.regions {
				if (len(args.FilterRegion) > 0 && !isFiltered(r.name, args.FilterRegion)) ||
					(len(args.FilterLocation) > 0 && !isFiltered(r.location, args.FilterLocation)) {
					continue
				}
				if r.endpoint == "" {
					code := r.name
					if r.code != "" {
						code = r.code
					}
					r.endpoint = fmt.Sprintf(provider.hostTemplate, code)
				}
				semaphore <- struct{}{}
				wg.Add(1)
				go func(region *region) {
					defer wg.Done()
					defer func() { <-semaphore }()
					region.rtt, region.err = endpointPing(region.endpoint)
					if region.rtt == 0 && region.err == nil {
						region.err = fmt.Errorf("timeout")
					}
					if region.err != nil {
						fmt.Print(".")
					} else {
						fmt.Print("!")
					}
				}(r)
				time.Sleep(time.Millisecond)
			}
		}
	}
	wg.Wait()

	clearLine := "\033[2K\r"
	if runtime.GOOS == "windows" {
		clearLine = "\r"
	}
	fmt.Print(clearLine)

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
	if args.Provider == "dns" && !args.All {
		t.SetTitle("Public DNS Servers")
		t.Style().Title = table.TitleOptions{Align: text.AlignCenter}
		t.AppendHeader(table.Row{"Provider", "Type", "Address", "RTT", "Status"})
	} else {
		t.AppendHeader(table.Row{"Provider", "Region", "Location", "RTT", "Status"})
	}
	var minRTT time.Duration
	for _, region := range results {
		status := text.FgGreen.Sprint("OK")
		if region.err != nil {
			if args.HideErrors {
				continue
			}
			status = text.FgRed.Sprintf("ERROR: %s", region.err)
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
		if args.Provider == "dns" && !args.All {
			t.AppendRow(table.Row{region.name, region.location, region.endpoint, rtt, status})
		} else {
			t.AppendRow(table.Row{region.code, region.name, region.location, rtt, status})
		}
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
			Name:     "Status",
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
