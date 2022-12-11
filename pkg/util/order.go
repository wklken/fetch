package util

import (
	"path/filepath"

	"github.com/wklken/httptest/pkg/config"
)

func GetRunningOrderedFiles(pathes []string, orders []config.Order) (files []string, err error) {
	hits := NewStringSet()
	for _, order := range orders {
		// pattern := order.Pattern
		var matches []string
		matches, err = filepath.Glob(order.Pattern)
		if err != nil {
			return
		}

		for _, m := range matches {
			if !hits.Has(m) {
				files = append(files, m)
				hits.Add(m)
			}
		}
	}

	// fmt.Println("the hits:", files)

	for _, p := range pathes {
		if !hits.Has(p) {
			files = append(files, p)
		}
	}

	return
}
