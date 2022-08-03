package main

import (
	"strings"
)

func isFiltered(needle string, filters []string) bool {
	for _, filter := range filters {
		if strings.Contains(strings.ToLower(needle), strings.ToLower(filter)) {
			return true
		}
	}
	return false
}
