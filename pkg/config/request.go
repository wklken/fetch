package config

import (
	"strings"

	"github.com/wklken/fetch/pkg/tpl"
)

type BasicAuth struct {
	Username string
	Password string
}

func (b *BasicAuth) Empty() bool {
	return b.Username == "" && b.Password == ""
}

type Request struct {
	Method string
	URL    string
	Body   string

	Header          map[string]string
	Cookie          string
	BasicAuth       BasicAuth `mapstructure:"basic_auth"`
	DisableRedirect bool      `mapstructure:"disable_redirect"`
	MaxRedirects    int       `mapstructure:"max_redirects"`
}

const TplBrace = "{{"

func (r *Request) Render(ctx map[string]interface{}) {
	if strings.Contains(r.Method, TplBrace) {
		r.Method = tpl.Render(r.Method, ctx)
	}

	if strings.Contains(r.URL, TplBrace) {
		r.URL = tpl.Render(r.URL, ctx)
	}

	if strings.Contains(r.Body, TplBrace) {
		r.Body = tpl.Render(r.Body, ctx)
	}

	for k, v := range r.Header {
		if strings.Contains(v, TplBrace) {
			r.Header[k] = tpl.Render(v, ctx)
		}
	}
}
