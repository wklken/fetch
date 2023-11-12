package config_test

import (
	"testing"

	"github.com/wklken/fetch/pkg/config"
)

func TestRequest_Render(t *testing.T) {
	type fields struct {
		Method string
		URL    string
		Body   string
		Header map[string]string
	}
	type args struct {
		ctx map[string]interface{}
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   fields
	}{
		{
			name: "no template",
			fields: fields{
				Method: "GET",
				URL:    "http://example.com",
				Body:   "",
				Header: map[string]string{"Content-Type": "application/json"},
			},
			args: args{ctx: map[string]interface{}{"foo": "bar"}},
			want: fields{
				Method: "GET",
				URL:    "http://example.com",
				Body:   "",
				Header: map[string]string{"Content-Type": "application/json"},
			},
		},
		{
			name: "with template",
			fields: fields{
				Method: "{{.method}}",
				URL:    "http://{{.host}}/path",
				Body:   "{\"foo\": \"{{.foo}}\"}",
				Header: map[string]string{"Content-Type": "{{.content_type}}"},
			},
			args: args{ctx: map[string]interface{}{
				"method":        "POST",
				"host":          "example.com",
				"foo":           "bar",
				"content_type":  "application/json",
				"another_field": "another_value",
			}},
			want: fields{
				Method: "POST",
				URL:    "http://example.com/path",
				Body:   "{\"foo\": \"bar\"}",
				Header: map[string]string{"Content-Type": "application/json"},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &config.Request{
				Method: tt.fields.Method,
				URL:    tt.fields.URL,
				Body:   tt.fields.Body,
				Header: tt.fields.Header,
			}
			r.Render(tt.args.ctx)
			if r.Method != tt.want.Method {
				t.Errorf("Request.Render() Method = %v, want %v", r.Method, tt.want.Method)
			}
			if r.URL != tt.want.URL {
				t.Errorf("Request.Render() URL = %v, want %v", r.URL, tt.want.URL)
			}
			if r.Body != tt.want.Body {
				t.Errorf("Request.Render() Body = %v, want %v", r.Body, tt.want.Body)
			}
			if len(r.Header) != len(tt.want.Header) {
				t.Errorf("Request.Render() Header length = %v, want %v", len(r.Header), len(tt.want.Header))
			} else {
				for k, v := range tt.want.Header {
					if r.Header[k] != v {
						t.Errorf("Request.Render() Header[%v] = %v, want %v", k, r.Header[k], v)
					}
				}
			}
		})
	}
}
