package assertion

import (
	"bytes"
	"fmt"

	"github.com/spf13/cast"
	"github.com/wklken/httptest/pkg/assert"
	"github.com/wklken/httptest/pkg/config"
	"github.com/wklken/httptest/pkg/util"
	"gopkg.in/xmlpath.v2"
)

func DoHTMLAssertions(body []byte, htmls []config.AssertHTML) (stats util.Stats) {
	root, err := xmlpath.ParseHTML(bytes.NewReader(body))
	if err != nil {
		stats.AddFailMessage("html parse fail: %s", err)
		stats.IncrFailAssertCountByN(int64(len(htmls)))
		return
	}

	// FIXME: same as xml, refactor it later
	for _, x := range htmls {
		path := x.Path
		expectedValue := x.Value
		stats.AddInfofMessage("assert.html.%s: ", path)

		p, err := xmlpath.Compile(path)
		if err != nil {
			message := fmt.Sprintf("wrong xpath %s, compile fail %s", path, err.Error())
			stats.AddFailMessage(message)
			stats.IncrFailAssertCount()
			continue
		}

		if actualValue, ok := p.String(root); ok {
			// convert expectedValue to String
			ev := cast.ToString(expectedValue)
			ok, message := assert.Equal(actualValue, ev)
			if ok {
				stats.AddPassMessage()
				stats.IncrOkAssertCount()
			} else {
				stats.AddFailMessage(message)
				stats.IncrFailAssertCount()
			}
		} else {
			stats.AddFailMessage("search html data fail, err=%s, path=%s, expected=%s", err, path, expectedValue)
			stats.IncrFailAssertCount()
		}

	}
	return
}
