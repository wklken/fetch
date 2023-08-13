package config

import (
	"fmt"
	"strings"

	"github.com/wklken/fetch/pkg/tpl"
)

// TODO: use reflect instead hard code
type Hook struct {
	SaveCookie   string `mapstructure:"save_cookie"`
	SaveResponse string `mapstructure:"save_response"`
	Exec         string `mapstructure:"exec"`
	Sleep        int    `mapstructure:"sleep"`
}

type Parse struct {
	Key    string `yaml:"key" mapstructure:"key"`
	Source string `yaml:"source" mapstructure:"source"`
	// body: json
	Jmespath string `yaml:"jmespath" mapstructure:"jmespath"`
	// header name
	Header string `yaml:"header" mapstructure:"header"`
}

type Case struct {
	Title       string
	Description string
	Path        string
	Index       int

	Config CaseConfig `mapstructure:"config"`
	Env    map[string]interface{}

	Request Request
	Assert  Assert
	Hook    Hook

	Parse []Parse

	// caseIndex => {lineNo: lineContent}
	FileLines map[int]map[int]string
	AllKeys   []string
}

func (c *Case) ID() string {
	if c.Index == 1 {
		if c.Title != "" {
			return fmt.Sprintf("%s | %s", c.Path, c.Title)
		} else {
			return fmt.Sprintf("%s | -", c.Path)
		}
	}

	if c.Title != "" {
		return fmt.Sprintf("%s[%d] | %s", c.Path, c.Index, c.Title)
	} else {
		return fmt.Sprintf("%s[%d] | -", c.Path, c.Index)
	}
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
func (c *Case) GuessAssertLineNumber(caseIndex int, key string) int {
	// fmt.Println("c.guess", caseIndex, key, c.FileLines[caseIndex])
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

	// protect
	if caseIndex > len(c.FileLines) {
		caseIndex = len(c.FileLines)
	}
	// caseIndex => {lineNo: lineContent}

	linesMapping := c.FileLines[caseIndex]
	// {lineNo: lineContent}

	for lineNo, lineContent := range linesMapping {
		for _, k := range keys {
			if strings.Contains(lineContent, k) {
				return lineNo
			}
		}
	}

	return -1
}
