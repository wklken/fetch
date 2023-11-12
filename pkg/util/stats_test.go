package util_test

import (
	"testing"

	"github.com/wklken/fetch/pkg/util"
)

func TestNewDummyStats(t *testing.T) {
	stats := util.NewDummyStats()

	if stats.GetOkAssertCount() != 0 {
		t.Errorf("Expected OkAssertCount to be 0, but got %d", stats.GetOkAssertCount())
	}

	if stats.GetFailAssertCount() != 0 {
		t.Errorf("Expected FailAssertCount to be 0, but got %d", stats.GetFailAssertCount())
	}
}
