package config

import (
	"strings"

	"github.com/wklken/fetch/pkg/tpl"
)

type Assert struct {
	Status      string   `yaml:"status"`
	StatusIn    []string `yaml:"status_in" mapstructure:"status_in"`
	StatusNotIn []string `yaml:"status_not_in" mapstructure:"status_not_in"`

	StatusCode      int   `yaml:"statusCode"`
	StatusCodeIn    []int `yaml:"statusCode_in" mapstructure:"statusCode_in"`
	StatusCodeNotIn []int `yaml:"statusCode_not_in" mapstructure:"statusCode_not_in"`
	StatusCodeLt    int   `yaml:"statusCode_lt" mapstructure:"statusCode_lt"`
	StatusCodeLte   int   `yaml:"statusCode_lte" mapstructure:"statusCode_lte"`
	StatusCodeGt    int   `yaml:"statusCode_gt" mapstructure:"statusCode_gt"`
	StatusCodeGte   int   `yaml:"statusCode_gte" mapstructure:"statusCode_gte"`

	ContentLength    int64 `yaml:"contentLength"`
	ContentLengthLt  int64 `yaml:"contentLength_lt" mapstructure:"contentLength_lt"`
	ContentLengthLte int64 `yaml:"contentLength_lte" mapstructure:"contentLength_lte"`
	ContentLengthGt  int64 `yaml:"contentLength_gt" mapstructure:"contentLength_gt"`
	ContentLengthGte int64 `yaml:"contentLength_gte" mapstructure:"contentLength_gte"`

	ContentType      string   `yaml:"contentType"`
	ContentTypeIn    []string `yaml:"contentType_in" mapstructure:"contentType_in"`
	ContentTypeNotIn []string `yaml:"contentType_not_in" mapstructure:"contentType_not_in"`

	// latency
	LatencyLt  int64 `yaml:"latency_lt" mapstructure:"latency_lt"`
	LatencyLte int64 `yaml:"latency_lte" mapstructure:"latency_lte"`
	LatencyGt  int64 `yaml:"latency_gt" mapstructure:"latency_gt"`
	LatencyGte int64 `yaml:"latency_gte" mapstructure:"latency_gte"`

	Body string `yaml:"body"`

	BodyContains      string `yaml:"body_contains" mapstructure:"body_contains"`
	BodyNotContains   string `yaml:"body_not_contains" mapstructure:"body_not_contains"`
	BodyIContains     string `yaml:"body_icontains" mapstructure:"body_icontains"`
	BodyStartsWith    string `yaml:"body_startswith" mapstructure:"body_startswith"`
	BodyEndsWith      string `yaml:"body_endswith" mapstructure:"body_endswith"`
	BodyNotStartsWith string `yaml:"body_not_startswith" mapstructure:"body_not_startswith"`
	BodyNotEndsWith   string `yaml:"body_not_endswith" mapstructure:"body_not_endswith"`
	BodyMatches       string `yaml:"body_matches" mapstructure:"body_matches"`
	BodyNotMatches    string `yaml:"body_not_matches" mapstructure:"body_not_matches"`

	Header              map[string]interface{} `yaml:"header"`
	HeaderExists        []string               `yaml:"header_exists" mapstructure:"header_exists"`
	HeaderValueMatches  map[string]string      `yaml:"header_value_matches" mapstructure:"header_value_matches"`
	HeaderValueContains map[string]string      `yaml:"header_value_contains" mapstructure:"header_value_contains"`

	JSON []AssertJSON `yaml:"json" mapstructure:"json"`
	XML  []AssertXML  `yaml:"xml" mapstructure:"xml"`
	HTML []AssertHTML `yaml:"html" mapstructure:"html"`
	YAML []AssertYAML `yaml:"yaml" mapstructure:"yaml"`
	TOML []AssertTOML `yaml:"toml" mapstructure:"toml"`

	Cookie       []AssertCookie `yaml:"cookie" mapstructure:"cookie"`
	CookieExists []string       `yaml:"cookie_exists" mapstructure:"cookie_exists"`

	// FIXME: cookie assert, should set_cookie key exist or key-value

	// if request fail like dial fail/context deadline exceeded, will do assert error_contains only,
	// will pass if the error message contains the string
	ErrorContains string `yaml:"error_contains" mapstructure:"error_contains"`

	HasRedirect      bool  `yaml:"has_redirect" mapstructure:"has_redirect"`
	RedirectCountLt  int64 `yaml:"redirectCount_lt" mapstructure:"redirectCount_lt"`
	RedirectCountLte int64 `yaml:"redirectCount_lte" mapstructure:"redirectCount_lte"`
	RedirectCountGt  int64 `yaml:"redirectCount_gt" mapstructure:"redirectCount_gt"`
	RedirectCountGte int64 `yaml:"redirectCount_gte" mapstructure:"redirectCount_gte"`

	Proto      string `yaml:"proto"`      // e.g. "HTTP/1.0"
	ProtoMajor int    `yaml:"protoMajor"` // e.g. 1
	ProtoMinor int    `yaml:"protoMinor"` // e.g. 0
	// TODO: gt/gte/lt/lte
}

func (a *Assert) Render(ctx map[string]interface{}) {
	a.Status = tpl.Render(a.Status, ctx)

	a.ContentType = tpl.Render(a.ContentType, ctx)

	if len(a.StatusIn) > 0 {
		n := make([]string, 0, len(a.StatusIn))
		for _, s := range a.StatusIn {
			n = append(n, tpl.Render(s, ctx))
		}
		a.StatusIn = n
	}
	if len(a.StatusNotIn) > 0 {
		n := make([]string, 0, len(a.StatusNotIn))
		for _, s := range a.StatusNotIn {
			n = append(n, tpl.Render(s, ctx))
		}
		a.StatusNotIn = n
	}
	if len(a.ContentTypeIn) > 0 {
		n := make([]string, 0, len(a.ContentTypeIn))
		for _, s := range a.ContentTypeIn {
			n = append(n, tpl.Render(s, ctx))
		}
		a.ContentTypeIn = n
	}
	if len(a.ContentTypeNotIn) > 0 {
		n := make([]string, 0, len(a.ContentTypeNotIn))
		for _, s := range a.ContentTypeNotIn {
			n = append(n, tpl.Render(s, ctx))
		}
		a.ContentTypeNotIn = n
	}

	a.Body = tpl.Render(a.Body, ctx)

	a.BodyContains = tpl.Render(a.BodyContains, ctx)
	a.BodyNotContains = tpl.Render(a.BodyNotContains, ctx)
	a.BodyStartsWith = tpl.Render(a.BodyStartsWith, ctx)
	a.BodyEndsWith = tpl.Render(a.BodyEndsWith, ctx)
	a.BodyNotStartsWith = tpl.Render(a.BodyNotStartsWith, ctx)
	a.BodyNotEndsWith = tpl.Render(a.BodyNotEndsWith, ctx)

	for _, j := range a.JSON {
		j.Render(ctx)
	}
}

type AssertJSON struct {
	Path  string
	Value interface{}
}

type AssertXML struct {
	Path  string
	Value interface{}
}

type AssertHTML struct {
	Path  string
	Value interface{}
}

type AssertYAML struct {
	Path  string
	Value interface{}
}

type AssertTOML struct {
	Path  string
	Value interface{}
}

type AssertCookie struct {
	// NOTE: https://go.dev/src/net/http/cookie_test.go
	// current only support check all equals(Name, Value, Domain, Path)
	Name   string
	Value  string
	Domain string
	Path   string
}

func (a *AssertJSON) Render(ctx map[string]interface{}) {
	if strings.Contains(a.Path, TplBrace) {
		a.Path = tpl.Render(a.Path, ctx)
	}
}
