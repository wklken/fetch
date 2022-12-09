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
		files = append(files, matches...)
		hits.Append(matches...)
	}

	for _, p := range pathes {
		if !hits.Has(p) {
			files = append(files, p)
		}
	}

	return
}
