package assertion

import (
	"fmt"
	"net/http"

	"github.com/wklken/fetch/pkg/assert"
	"github.com/wklken/fetch/pkg/config"
	"github.com/wklken/fetch/pkg/util"
)

func DoHeaderAssertions(c *config.Case, respHeader http.Header) (stats util.Stats) {
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
				lineNumber := c.GuessAssertLineNumber(c.Index, key)
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
				lineNumber := c.GuessAssertLineNumber(c.Index, "header_exists")
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

	// header value matches
	if len(c.Assert.HeaderValueMatches) > 0 {
		for key, valueRegex := range c.Assert.HeaderValueMatches {
			stats.AddInfofMessage("assert.header_value_matches.%s: ", key)

			ok, message := assert.Matches(respHeader.Get(key), valueRegex)
			if ok {
				stats.AddPassMessage()
				stats.IncrOkAssertCount()
			} else {
				// the ka.key is like assert.latency_lt
				lineNumber := c.GuessAssertLineNumber(c.Index, key)
				if lineNumber > 0 {
					message = fmt.Sprintf("line:%d | %s", lineNumber, message)
				}
				stats.AddFailMessage(message)
				stats.IncrFailAssertCount()
			}
		}
	}

	// key-value
	if len(c.Assert.HeaderValueContains) > 0 {
		for key, value := range c.Assert.HeaderValueContains {
			stats.AddInfofMessage("assert.header_value_contains.%s: ", key)
			ok, message := assert.Contains(respHeader.Get(key), value)
			if ok {
				stats.AddPassMessage()
				stats.IncrOkAssertCount()
			} else {
				// the ka.key is like assert.latency_lt
				lineNumber := c.GuessAssertLineNumber(c.Index, key)
				if lineNumber > 0 {
					message = fmt.Sprintf("line:%d | %s", lineNumber, message)
				}
				stats.AddFailMessage(message)
				stats.IncrFailAssertCount()
			}
		}
	}

	return
}
