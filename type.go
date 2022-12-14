package main

import "time"

type (
	region struct {
		name     string
		location string
		endpoint string
		code     string
		rtt      time.Duration
		err      error
	}

	provider map[string]struct {
		name         string
		hostTemplate string
		regions      []*region
	}
)
