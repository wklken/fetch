package assertion

import (
	"fmt"
	"net/http"

	"github.com/wklken/httptest/pkg/assert"
	"github.com/wklken/httptest/pkg/config"
	"github.com/wklken/httptest/pkg/util"
)

func DoHeaderAssertions(c config.Case, respHeader http.Header) (stats util.Stats) {
	for key, value := range c.Assert.Header {
		stats.AddInfofMessage("assert.header.%s: ", key)
		ok, message := assert.Equal(respHeader.Get(key), value)
		if ok {
			stats.AddPassMessage()
			stats.IncrOkAssertCount()
		} else {
			// the ka.key is like assert.latency_lt
			lineNumber := c.GuessAssertLineNumber(key)
			if lineNumber > 0 {
				message = fmt.Sprintf("line:%d | %s", lineNumber, message)
			}
			stats.AddFailMessage(message)
			stats.IncrFailAssertCount()
		}
	}
	return
}
