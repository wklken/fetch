package assertion

import (
	"fmt"

	"github.com/wklken/fetch/pkg/assert"
	"github.com/wklken/fetch/pkg/config"
	"github.com/wklken/fetch/pkg/util"
)

func DoErrorAssertions(c *config.Case, err error) (stats util.Stats) {
	stats.AddInfofMessage("assert.error_contains: ")
	// ok, message := assert.Contains(err.Error(), c.Assert.ErrorContains)
	ok, message := assert.StringContainsAll(err.Error(), c.Assert.ErrorContains)
	if ok {
		stats.AddPassMessage()
		stats.IncrOkAssertCount()
	} else {
		// the ka.key is like assert.latency_lt
		lineNumber := c.GuessAssertLineNumber(c.Index, "assert.error_contains")
		if lineNumber > 0 {
			message = fmt.Sprintf("line:%d | %s", lineNumber, message)
		}
		stats.AddFailMessage(message)
		stats.IncrFailAssertCount()
	}

	return
}
