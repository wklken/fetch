package assertion

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/wklken/httptest/pkg/assert"
	"github.com/wklken/httptest/pkg/config"
	"github.com/wklken/httptest/pkg/util"
)

type Ctx struct {
	f        assert.AssertFunc
	element1 interface{}
	element2 interface{}
}

type keyAssert struct {
	key string
	ctx Ctx
}

func DoKeysAssertion(
	allKeys *util.StringSet,
	resp *http.Response,
	c config.Case,
	hasRedirect bool,
	latency int64,
	contentType string,
	body []byte,
) (stats util.Stats) {
	bodyStr := strings.TrimSuffix(string(body), "\n")

	// NOTE: the order
	keyAsserts := []keyAssert{
		// statuscode
		{
			key: "assert.statuscode",
			ctx: Ctx{
				f:        assert.Equal,
				element1: resp.StatusCode,
				element2: c.Assert.StatusCode,
			},
		},
		{
			key: "assert.statuscode_lt",
			ctx: Ctx{
				f:        assert.Less,
				element1: resp.StatusCode,
				element2: c.Assert.StatusCodeLt,
			},
		},
		{
			key: "assert.statuscode_lte",
			ctx: Ctx{
				f:        assert.LessOrEqual,
				element1: resp.StatusCode,
				element2: c.Assert.StatusCodeLte,
			},
		},
		{
			key: "assert.statuscode_gt",
			ctx: Ctx{
				f:        assert.Greater,
				element1: resp.StatusCode,
				element2: c.Assert.StatusCodeGt,
			},
		},
		{
			key: "assert.statuscode_gte",
			ctx: Ctx{
				f:        assert.GreaterOrEqual,
				element1: resp.StatusCode,
				element2: c.Assert.StatusCodeGte,
			},
		},
		{
			key: "assert.statuscode_in",
			ctx: Ctx{
				f:        assert.In,
				element1: resp.StatusCode,
				element2: c.Assert.StatusCodeIn,
			},
		},
		{
			key: "assert.statuscode_not_in",
			ctx: Ctx{
				f:        assert.NotIn,
				element1: resp.StatusCode,
				element2: c.Assert.StatusCodeNotIn,
			},
		},
		// status
		{
			key: "assert.status",
			ctx: Ctx{
				f:        assert.Equal,
				element1: strings.ToLower(http.StatusText(resp.StatusCode)),
				element2: strings.ToLower(c.Assert.Status),
			},
		},
		{
			key: "assert.status_in",
			ctx: Ctx{
				f:        assert.In,
				element1: strings.ToLower(http.StatusText(resp.StatusCode)),
				element2: util.ToLower(c.Assert.StatusIn),
			},
		},
		{
			key: "assert.status_not_in",
			ctx: Ctx{
				f:        assert.NotIn,
				element1: strings.ToLower(http.StatusText(resp.StatusCode)),
				element2: util.ToLower(c.Assert.StatusNotIn),
			},
		},
		{
			key: "assert.contenttype",
			ctx: Ctx{
				f:        assert.Equal,
				element1: strings.ToLower(contentType),
				element2: strings.ToLower(c.Assert.ContentType),
			},
		},
		{
			key: "assert.contenttype_in",
			ctx: Ctx{
				f:        assert.In,
				element1: strings.ToLower(contentType),
				element2: util.ToLower(c.Assert.ContentTypeIn),
			},
		},
		{
			key: "assert.contenttype_not_in",
			ctx: Ctx{
				f:        assert.NotIn,
				element1: strings.ToLower(contentType),
				element2: util.ToLower(c.Assert.ContentTypeNotIn),
			},
		},
		// contentlength
		{
			key: "assert.contentlength",
			ctx: Ctx{
				f:        assert.Equal,
				element1: resp.ContentLength,
				element2: c.Assert.ContentLength,
			},
		},
		{
			key: "assert.contentlength_lt",
			ctx: Ctx{
				f:        assert.Less,
				element1: resp.ContentLength,
				element2: c.Assert.ContentLengthLt,
			},
		},
		{
			key: "assert.contentlength_lte",
			ctx: Ctx{
				f:        assert.LessOrEqual,
				element1: resp.ContentLength,
				element2: c.Assert.ContentLengthLte,
			},
		},
		{
			key: "assert.contentlength_gt",
			ctx: Ctx{
				f:        assert.Greater,
				element1: resp.ContentLength,
				element2: c.Assert.ContentLengthGt,
			},
		},
		{
			key: "assert.contentlength_gte",
			ctx: Ctx{
				f:        assert.GreaterOrEqual,
				element1: resp.ContentLength,
				element2: c.Assert.ContentLengthGte,
			},
		},
		// latency
		{
			key: "assert.latency_lt",
			ctx: Ctx{
				f:        assert.Less,
				element1: latency,
				element2: c.Assert.LatencyLt,
			},
		},
		{
			key: "assert.latency_lte",
			ctx: Ctx{
				f:        assert.LessOrEqual,
				element1: latency,
				element2: c.Assert.LatencyLte,
			},
		},
		{
			key: "assert.latency_gt",
			ctx: Ctx{
				f:        assert.Greater,
				element1: latency,
				element2: c.Assert.LatencyGt,
			},
		},
		{
			key: "assert.latency_gte",
			ctx: Ctx{
				f:        assert.GreaterOrEqual,
				element1: latency,
				element2: c.Assert.LatencyGte,
			},
		},
		// body
		{
			key: "assert.body",
			ctx: Ctx{
				f:        assert.Equal,
				element1: bodyStr,
				element2: c.Assert.Body,
			},
		},
		{
			key: "assert.body_contains",
			ctx: Ctx{
				f:        assert.Contains,
				element1: bodyStr,
				element2: c.Assert.BodyContains,
			},
		},
		{
			key: "assert.body_not_contains",
			ctx: Ctx{
				f:        assert.NotContains,
				element1: bodyStr,
				element2: c.Assert.BodyNotContains,
			},
		},
		{
			key: "assert.body_startswith",
			ctx: Ctx{
				f:        assert.StartsWith,
				element1: bodyStr,
				element2: c.Assert.BodyStartsWith,
			},
		},
		{
			key: "assert.body_endswith",
			ctx: Ctx{
				f:        assert.EndsWith,
				element1: bodyStr,
				element2: c.Assert.BodyEndsWith,
			},
		},
		{
			key: "assert.body_not_startswith",
			ctx: Ctx{
				f:        assert.NotStartsWith,
				element1: bodyStr,
				element2: c.Assert.BodyNotStartsWith,
			},
		},
		{
			key: "assert.body_not_endswith",
			ctx: Ctx{
				f:        assert.NotEndsWith,
				element1: bodyStr,
				element2: c.Assert.BodyNotEndsWith,
			},
		},
		{
			key: "assert.hasredirect",
			ctx: Ctx{
				f:        assert.Equal,
				element1: hasRedirect,
				element2: c.Assert.HasRedirect,
			},
		},
		{
			key: "assert.proto",
			ctx: Ctx{
				f:        assert.Equal,
				element1: resp.Proto,
				element2: c.Assert.Proto,
			},
		},
		{
			key: "assert.protomajor",
			ctx: Ctx{
				f:        assert.Equal,
				element1: resp.ProtoMajor,
				element2: c.Assert.ProtoMajor,
			},
		},
		{
			key: "assert.protominor",
			ctx: Ctx{
				f:        assert.Equal,
				element1: resp.ProtoMinor,
				element2: c.Assert.ProtoMinor,
			},
		},
	}

	for _, ka := range keyAsserts {
		if allKeys.Has(ka.key) {
			stats.AddInfofMessage("%s: ", ka.key)
			ok, message := ka.ctx.f(ka.ctx.element1, ka.ctx.element2)
			if ok {
				stats.AddPassMessage()
				stats.IncrOkAssertCount()
			} else {
				// the ka.key is like assert.latency_lt
				lineNumber := c.GuessAssertLineNumber(ka.key)
				if lineNumber > 0 {
					message = fmt.Sprintf("line:%d | %s", lineNumber, message)
				}
				stats.AddFailMessage(message)
				stats.IncrFailAssertCount()
			}
		}
	}

	return stats
}