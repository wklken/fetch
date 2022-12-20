package assertion

import (
	"fmt"
	"net/http"

	"github.com/wklken/httptest/pkg/assert"
	"github.com/wklken/httptest/pkg/config"
	"github.com/wklken/httptest/pkg/util"
)

func DoHeaderAssertions(c config.Case, respHeader http.Header) (stats util.Stats) {
	// key-value
	if len(c.Assert.Header) > 0 {
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
	}

	// header exists, all keys must present in response header
	if len(c.Assert.HeaderExists) > 0 {
		allOK := true
		stats.AddInfofMessage("assert.header_exists.%s: ", util.PrettyStringSlice(c.Assert.HeaderExists))
		for _, key := range c.Assert.HeaderExists {
			ok := respHeader.Get(key) != ""
			if !ok {
				message := fmt.Sprintf("header key `%s` not exists", key)
				lineNumber := c.GuessAssertLineNumber("header_exists")
				if lineNumber > 0 {
					message = fmt.Sprintf("line:%d | %s", lineNumber, message)
				}
				stats.AddFailMessage(message)
				stats.IncrFailAssertCount()

				allOK = false
				break
			}
		}
		if allOK {
			stats.AddPassMessage()
			stats.IncrOkAssertCount()
		}
	}

	return
}
