package config

import (
	"fmt"
	"strings"

	"github.com/wklken/httptest/pkg/tpl"
)

// TODO: use reflect instead hard code
type Hook struct {
	SaveCookie string `mapstructure:"save_cookie"`
}

type Case struct {
	Title       string
	Description string

	Config CaseConfig `mapstructure:"config"`
	Env    map[string]interface{}

	Request Request
	Assert  Assert
	Hook    Hook

	FileLines []string
}

func (c *Case) Render(ctx map[string]interface{}) {
	if strings.Contains(c.Title, TplBrace) {
		c.Title = tpl.Render(c.Title, ctx)
	}
	if strings.Contains(c.Description, TplBrace) {
		c.Description = tpl.Render(c.Description, ctx)
	}

	c.Request.Render(ctx)
	c.Assert.Render(ctx)
}

// GetAssertLineNumber will guess the assertion line number
// toml: status = "ok"
// json: "status": "ok"
// yaml: status: ok
// ini: status=ok
func (c *Case) GuessAssertLineNumber(key string) int {
	parts := strings.Split(key, ".")
	if len(parts) > 0 {
		key = parts[len(parts)-1]
	}

	keys := []string{
		fmt.Sprintf("%s=", key),
		fmt.Sprintf("%s =", key),
		fmt.Sprintf(`"%s":`, key),
		fmt.Sprintf(`"%s" :`, key),
		fmt.Sprintf(`%s:`, key),
		fmt.Sprintf(`%s :`, key),
		fmt.Sprintf(`%s=`, key),
		fmt.Sprintf(`%s =`, key),
	}

	// NOTE: maybe get the wrong line number!
	// scan from the end of the file
	count := len(c.FileLines)
	for i := count - 1; i >= 0; i-- {
		for _, k := range keys {
			if strings.Contains(c.FileLines[i], k) {
				return i + 1
			}
		}
	}

	return -1
}
