package assertion

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/goccy/go-yaml"
	"github.com/wklken/fetch/pkg/config"
	"github.com/wklken/fetch/pkg/util"
)

func DoYAMLAssertions(body []byte, yamls []config.AssertYAML) (stats util.Stats) {
	jsonBody, err := yaml.YAMLToJSON(body)
	if err != nil {
		stats.AddFailMessage("yaml parse and convert to json fail: %s", err)
		stats.IncrFailAssertCountByN(int64(len(yamls)))
		return
	}

	jsons := make([]config.AssertJSON, 0, len(yamls))
	for _, x := range yamls {
		jsons = append(jsons, config.AssertJSON{
			Path:  x.Path,
			Value: x.Value,
		})
	}
	var jsonData interface{}
	err = binding.JSON.BindBody(jsonBody, &jsonData)
	if err != nil {
		stats.AddFailMessage("binding.json fail: %s", err)
		stats.IncrFailAssertCountByN(int64(len(yamls)))
		return
	}

	return DoJSONAssertions(jsonData, jsons)
}
