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

func DoCookieAssertions(cookies []*http.Cookie, assertCookies []config.AssertCookie) (stats util.Stats) {
	cookieUniqueKeys := []string{}
	for _, c := range cookies {
		cookieUniqueKeys = append(cookieUniqueKeys, c.String())
	}

	// NOTE: currently only match the simplest case,
	// should be refactor to support complex case, like max-age
	for _, x := range assertCookies {
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

	return stats
}
