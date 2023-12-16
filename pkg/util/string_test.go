package util_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wklken/fetch/pkg/util"
)

func TestTruncateBytes(t *testing.T) {
	type args struct {
		content []byte
		length  int
	}
	tests := []struct {
		name string
		args args
		want []byte
	}{
		{
			name: "truncate bytes to shorter length",
			args: args{
				content: []byte("hello world"),
				length:  5,
			},
			want: []byte("hello"),
		},
		{
			name: "truncate bytes to same length",
			args: args{
				content: []byte("hello world"),
				length:  11,
			},
			want: []byte("hello world"),
		},
		{
			name: "truncate bytes to longer length",
			args: args{
				content: []byte("hello"),
				length:  10,
			},
			want: []byte("hello"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := util.TruncateBytes(tt.args.content, tt.args.length); !assert.Equal(t, got, tt.want) {
				t.Errorf("TruncateBytes() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTruncateBytesToString(t *testing.T) {
	type args struct {
		content []byte
		length  int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "truncate bytes to shorter length",
			args: args{
				content: []byte("hello world"),
				length:  5,
			},
			want: "hello",
		},
		{
			name: "truncate bytes to same length",
			args: args{
				content: []byte("hello world"),
				length:  11,
			},
			want: "hello world",
		},
		{
			name: "truncate bytes to longer length",
			args: args{
				content: []byte("hello"),
				length:  10,
			},
			want: "hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := util.TruncateBytesToString(tt.args.content, tt.args.length); !assert.Equal(t, got, tt.want) {
				t.Errorf("TruncateBytesToString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestTruncateString(t *testing.T) {
	type args struct {
		s string
		n int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "truncate string to shorter length",
			args: args{
				s: "hello world",
				n: 5,
			},
			want: "hello",
		},
		{
			name: "truncate string to same length",
			args: args{
				s: "hello world",
				n: 11,
			},
			want: "hello world",
		},
		{
			name: "truncate string to longer length",
			args: args{
				s: "hello",
				n: 10,
			},
			want: "hello",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := util.TruncateString(tt.args.s, tt.args.n); !assert.Equal(t, got, tt.want) {
				t.Errorf("TruncateString() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestOmitMiddle(t *testing.T) {
	type args struct {
		s    string
		head int
		tail int
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "omit middle of string",
			args: args{
				s:    "hello world",
				head: 2,
				tail: 2,
			},
			want: "he...ld",
		},
		{
			name: "omit middle of string with no middle",
			args: args{
				s:    "hello",
				head: 2,
				tail: 2,
			},
			want: "he...lo",
		},
		{
			name: "omit middle of string with no head",
			args: args{
				s:    "hello world",
				head: 0,
				tail: 2,
			},
			want: "...ld",
		},
		{
			name: "omit middle of string with no tail",
			args: args{
				s:    "hello world",
				head: 2,
				tail: 0,
			},
			want: "he...",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := util.OmitMiddle(tt.args.s, tt.args.head, tt.args.tail); !assert.Equal(t, got, tt.want) {
				t.Errorf("OmitMiddle() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPrettyStringSlice(t *testing.T) {
	type args struct {
		s []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "pretty print string slice",
			args: args{
				s: []string{"hello", "world"},
			},
			want: "[hello, world]",
		},
		{
			name: "pretty print empty string slice",
			args: args{
				s: []string{},
			},
			want: "[]",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := util.PrettyStringSlice(tt.args.s); !assert.Equal(t, got, tt.want) {
				t.Errorf("PrettyStringSlice() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringArrayMapFunc(t *testing.T) {
	type args struct {
		elements []string
		f        func(string) string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "apply function to string array",
			args: args{
				elements: []string{"hello", "world"},
				f:        strings.ToUpper,
			},
			want: []string{"HELLO", "WORLD"},
		},
		{
			name: "apply function to empty string array",
			args: args{
				elements: []string{},
				f:        strings.ToUpper,
			},
			want: []string(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := util.StringArrayMapFunc(tt.args.elements, tt.args.f); !assert.Equal(t, got, tt.want) {
				t.Errorf("StringArrayMapFunc() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAnyStringHasPrefix(t *testing.T) {
	type args struct {
		l      []string
		prefix string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "prefix exists in one string",
			args: args{
				l:      []string{"hello", "world", "foo"},
				prefix: "w",
			},
			want: true,
		},
		{
			name: "prefix exists in multiple strings",
			args: args{
				l:      []string{"hello", "world", "foo"},
				prefix: "f",
			},
			want: true,
		},
		{
			name: "prefix does not exist in any string",
			args: args{
				l:      []string{"hello", "world", "foo"},
				prefix: "bar",
			},
			want: false,
		},
		{
			name: "empty slice",
			args: args{
				l:      []string{},
				prefix: "foo",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := util.AnyStringHasPrefix(tt.args.l, tt.args.prefix); got != tt.want {
				t.Errorf("AnyStringHasPrefix() = %v, want %v", got, tt.want)
			}
		})
	}
}
