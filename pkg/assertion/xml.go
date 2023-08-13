package assertion

import (
	"bytes"
	"encoding/xml"
	"fmt"

	"github.com/spf13/cast"
	"github.com/wklken/fetch/pkg/assert"
	"github.com/wklken/fetch/pkg/config"
	"github.com/wklken/fetch/pkg/util"
	"golang.org/x/net/html/charset"
	"gopkg.in/xmlpath.v2"
)

func DoXMLAssertions(body []byte, xmls []config.AssertXML) (stats util.Stats) {
	decoder := xml.NewDecoder(bytes.NewReader(body))
	decoder.CharsetReader = charset.NewReaderLabel
	root, err := xmlpath.ParseDecoder(decoder)
	// root, err := xmlpath.Parse(bytes.NewReader(body))
	if err != nil {
		stats.AddFailMessage("xml parse fail: %s", err)
		stats.IncrFailAssertCountByN(int64(len(xmls)))
		return
	}

	for _, x := range xmls {
		path := x.Path
		expectedValue := x.Value
		stats.AddInfofMessage("assert.xml.%s: ", path)

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
			stats.AddFailMessage("search xml data fail, err=%s, path=%s, expected=%s", err, path, expectedValue)
			stats.IncrFailAssertCount()
		}
	}

	return
}
