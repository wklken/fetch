package assertion

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin/binding"
	"github.com/pelletier/go-toml/v2"
	"github.com/wklken/fetch/pkg/config"
	"github.com/wklken/fetch/pkg/util"
)

func DoTOMLAssertions(body []byte, tomls []config.AssertTOML) (stats util.Stats) {
	var v interface{}

	d := toml.NewDecoder(bytes.NewReader(body))
	err := d.Decode(&v)
	if err != nil {
		var derr *toml.DecodeError
		if errors.As(err, &derr) {
			row, col := derr.Position()
			err = fmt.Errorf("%s\nerror occurred at row %d column %d", derr.String(), row, col)
		}
		stats.AddFailMessage("toml parse and convert to json fail: %s", err)
		stats.IncrFailAssertCountByN(int64(len(tomls)))
		return
	}

	jsonBody, err := json.Marshal(v)
	if err != nil {
		stats.AddFailMessage("toml parse and convert to json fail: %s", err)
		stats.IncrFailAssertCountByN(int64(len(tomls)))
		return
	}

	jsons := make([]config.AssertJSON, 0, len(tomls))
	for _, x := range tomls {
		jsons = append(jsons, config.AssertJSON{
			Path:  x.Path,
			Value: x.Value,
		})
	}
	var jsonData interface{}
	err = binding.JSON.BindBody(jsonBody, &jsonData)
	if err != nil {
		stats.AddFailMessage("binding.json fail: %s", err)
		stats.IncrFailAssertCountByN(int64(len(tomls)))
		return
	}

	return DoJSONAssertions(jsonData, jsons)
}
