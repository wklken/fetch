package assert_test

import (
	"testing"

	"github.com/wklken/fetch/pkg/assert"
)

func TestStartsWith(t *testing.T) {
	type args struct {
		s      string
		prefix string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test starts with",
			args: args{
				s:      "hello world",
				prefix: "hello",
			},
			want: true,
		},
		{
			name: "test does not start with",
			args: args{
				s:      "hello world",
				prefix: "world",
			},
			want: false,
		},
		{
			name: "test empty string",
			args: args{
				s:      "",
				prefix: "",
			},
			want: true,
		},
		{
			name: "test prefix longer than string",
			args: args{
				s:      "hello",
				prefix: "hello world",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := assert.StartsWith(tt.args.s, tt.args.prefix); got != tt.want {
				t.Errorf("StartsWith() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEndsWith(t *testing.T) {
	type args struct {
		s      string
		suffix string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test ends with",
			args: args{
				s:      "hello world",
				suffix: "world",
			},
			want: true,
		},
		{
			name: "test does not end with",
			args: args{
				s:      "hello world",
				suffix: "hello",
			},
			want: false,
		},
		{
			name: "test empty string",
			args: args{
				s:      "",
				suffix: "",
			},
			want: true,
		},
		{
			name: "test suffix longer than string",
			args: args{
				s:      "hello",
				suffix: "hello world",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := assert.EndsWith(tt.args.s, tt.args.suffix); got != tt.want {
				t.Errorf("EndsWith() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotStartsWith(t *testing.T) {
	type args struct {
		s      string
		prefix string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test does not start with",
			args: args{
				s:      "hello world",
				prefix: "world",
			},
			want: true,
		},
		{
			name: "test starts with",
			args: args{
				s:      "hello world",
				prefix: "hello",
			},
			want: false,
		},
		{
			name: "test empty string",
			args: args{
				s:      "",
				prefix: "",
			},
			want: false,
		},
		{
			name: "test prefix longer than string",
			args: args{
				s:      "hello",
				prefix: "hello world",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := assert.NotStartsWith(tt.args.s, tt.args.prefix); got != tt.want {
				t.Errorf("NotStartsWith() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotEndsWith(t *testing.T) {
	type args struct {
		s      string
		suffix string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test does not end with",
			args: args{
				s:      "hello world",
				suffix: "hello",
			},
			want: true,
		},
		{
			name: "test ends with",
			args: args{
				s:      "hello world",
				suffix: "world",
			},
			want: false,
		},
		{
			name: "test empty string",
			args: args{
				s:      "",
				suffix: "",
			},
			want: false,
		},
		{
			name: "test suffix longer than string",
			args: args{
				s:      "hello",
				suffix: "hello world",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := assert.NotEndsWith(tt.args.s, tt.args.suffix); got != tt.want {
				t.Errorf("NotEndsWith() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringContainsAll(t *testing.T) {
	type args struct {
		s        string
		elements []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test contains all",
			args: args{
				s:        "hello world",
				elements: []string{"hello", "world"},
			},
			want: true,
		},
		{
			name: "test does not contain all",
			args: args{
				s:        "hello world",
				elements: []string{"hello", "world", "foo"},
			},
			want: false,
		},
		{
			name: "test empty string",
			args: args{
				s:        "",
				elements: []string{},
			},
			want: true,
		},
		{
			name: "test empty elements",
			args: args{
				s:        "hello world",
				elements: []string{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := assert.StringContainsAll(tt.args.s, tt.args.elements); got != tt.want {
				t.Errorf("StringContainsAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringNotContainsAll(t *testing.T) {
	type args struct {
		s        string
		elements []string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test does not contain all",
			args: args{
				s:        "hello world",
				elements: []string{"foo", "bar"},
			},
			want: true,
		},
		{
			name: "test contains all",
			args: args{
				s:        "hello world",
				elements: []string{"foo", "bar", "hello"},
			},
			want: false,
		},
		{
			name: "test empty string",
			args: args{
				s:        "",
				elements: []string{},
			},
			want: true,
		},
		{
			name: "test empty elements",
			args: args{
				s:        "hello world",
				elements: []string{},
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := assert.StringNotContainsAll(tt.args.s, tt.args.elements); got != tt.want {
				t.Errorf("StringNotContainsAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
