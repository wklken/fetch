package assertion

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/wklken/httptest/pkg/config"
	"github.com/wklken/httptest/pkg/util"
)

func anyStringHasPrefix(l []string, prefix string) bool {
	for _, x := range l {
		if strings.HasPrefix(x, prefix) {
			return true
		}
	}
	return false
}

func DoCookieAssertions(
	c config.Case,
	cookies []*http.Cookie,
) (stats util.Stats) {
	cookieUniqueKeys := []string{}
	cookieNames := util.NewFixedLengthStringSet(len(cookies))
	for _, c := range cookies {
		cookieUniqueKeys = append(cookieUniqueKeys, c.String())
		cookieNames.Add(c.Name)
	}

	// NOTE: currently only match the simplest case,
	// should be refactor to support complex case, like max-age
	for _, x := range c.Assert.Cookie {
		var cookieKey string
		if x.Domain != "" {
			cookieKey = fmt.Sprintf("%s=%s; Domain=%s", x.Name, x.Value, x.Domain)
		} else if x.Path != "" {
			cookieKey = fmt.Sprintf("%s=%s; Path=%s", x.Name, x.Value, x.Path)
		} else {
			cookieKey = fmt.Sprintf("%s=%s", x.Name, x.Value)
		}

		stats.AddInfofMessage("assert.cookie.[%s]: ", cookieKey)

		if anyStringHasPrefix(cookieUniqueKeys, cookieKey) {
			stats.AddPassMessage()
			stats.IncrOkAssertCount()
		} else {
			stats.AddFailMessage(fmt.Sprintf("no cookie equals to `%s`", cookieKey))
			stats.IncrFailAssertCount()
		}
	}

	if len(c.Assert.CookieExists) > 0 {
		allOK := true
		stats.AddInfofMessage("assert.cookie_exists.%s: ", util.PrettyStringSlice(c.Assert.CookieExists))
		for _, name := range c.Assert.CookieExists {
			if !cookieNames.Has(name) {
				message := fmt.Sprintf("cookie name `%s` not exists", name)
				lineNumber := c.GuessAssertLineNumber("cookie_exists")
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

	return stats
}
