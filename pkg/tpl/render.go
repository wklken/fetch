package tpl

import (
	"bytes"
	"text/template"

	"github.com/wklken/fetch/pkg/log"
)

func Render(s string, ctx map[string]interface{}) string {
	t, err := template.New("tmp").Parse(s)
	if err != nil {
		log.Warning("render string `%s` fail", s)
		return s
	}

	// TODO: what todo if render fail??????

	var rs bytes.Buffer
	err = t.Execute(&rs, ctx)
	if err != nil {
		log.Warning("render string `%s` fail", s)
		return s
	}

	return rs.String()
}
