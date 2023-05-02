package main

import (
	"fmt"
	"sort"
	"time"
)

type config struct {
	All            bool          `arg:"--all,-a" help:"Scan all providers"`
	AltPing        bool          `arg:"--alt-ping" help:"Use alternative ICMP ping method"`
	Count          int           `arg:"--count,-c" default:"4" help:"Number of pings to send"`
	HideErrors     bool          `arg:"--hide-errors,-e" help:"Hide errors from results"`
	FilterRegion   []string      `arg:"--region,-r,separate" help:"Filter by regions, can be specified multiple times"`
	FilterLocation []string      `arg:"--location,-l,separate" help:"Filter by location, can be specified multiple times"`
	Provider       string        `arg:"--provider,-p" default:"aws" help:"Choose provider for ping"`
	Timeout        time.Duration `arg:"--timeout,-t" default:"5s" help:"Timeout before ping ends"`
}

func (c *config) Description() string {
	keys := make([]string, 0, len(providers))
	for k := range providers {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	desc := "Supported providers for cloud ping:\n"
	for _, k := range keys {
		desc += fmt.Sprintf("%s - %s\n", k, providers[k].name)
	}
	return desc
}
