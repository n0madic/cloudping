package main

import (
	"fmt"
	"sort"
	"time"
)

type config struct {
	All            bool          `arg:"--all,-a" help:"Scan all providers"`
	AltPing        bool          `arg:"--alt-ping" help:"Use alternative ICMP ping method"`
	Concurrent     int           `arg:"--concurrent,-C" placeholder:"NUM" default:"10" help:"Number of concurrent pings"`
	Count          int           `arg:"--count,-c" placeholder:"NUM" default:"4" help:"Number of pings to send"`
	HideErrors     bool          `arg:"--hide-errors,-e" help:"Hide errors from results"`
	FilterRegion   []string      `arg:"--region,-r,separate" placeholder:"NAME" help:"Filter by regions, can be specified multiple times"`
	FilterLocation []string      `arg:"--location,-l,separate" placeholder:"NAME" help:"Filter by location, can be specified multiple times"`
	Provider       string        `arg:"--provider,-p" placeholder:"NAME" default:"aws" help:"Choose provider for ping"`
	Timeout        time.Duration `arg:"--timeout,-t" placeholder:"DURATION" default:"5s" help:"Timeout before ping ends"`
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
