package assert_test

import (
	"testing"

	"github.com/wklken/fetch/pkg/assert"
)

func TestRegexpMatch(t *testing.T) {
	type args struct {
		text string
		expr string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "test matches",
			args: args{
				text: "hello world",
				expr: "hello",
			},
			want: true,
		},
		{
			name: "test not matches",
			args: args{
				text: "hello world",
				expr: "goodbye",
			},
			want: false,
		},
		{
			name: "test empty string",
			args: args{
				text: "",
				expr: "hello",
			},
			want: false,
		},
		{
			name: "test empty expression",
			args: args{
				text: "hello world",
				expr: "",
			},
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := assert.RegexpMatch(tt.args.text, tt.args.expr); got != tt.want {
				t.Errorf("RegexpMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMatches(t *testing.T) {
	type args struct {
		text interface{}
		expr interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "test matches",
			args: args{
				text: "hello world",
				expr: "hello",
			},
			want: true,
		},
		{
			name: "test not matches",
			args: args{
				text: "hello world",
				expr: "goodbye",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := assert.Matches(tt.args.text, tt.args.expr); got != tt.want {
				t.Errorf("Matches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNotMatches(t *testing.T) {
	type args struct {
		text interface{}
		expr interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "test not matches",
			args: args{
				text: "hello world",
				expr: "goodbye",
			},
			want: true,
		},
		{
			name: "test matches",
			args: args{
				text: "hello world",
				expr: "hello",
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := assert.NotMatches(tt.args.text, tt.args.expr); got != tt.want {
				t.Errorf("NotMatches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringMatchesAll(t *testing.T) {
	type args struct {
		text  interface{}
		exprs interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "test matches all",
			args: args{
				text:  "hello world",
				exprs: []string{"hello", "world"},
			},
			want: true,
		},
		{
			name: "test not matches all",
			args: args{
				text:  "hello world",
				exprs: []string{"hello", "goodbye"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := assert.StringMatchesAll(tt.args.text, tt.args.exprs); got != tt.want {
				t.Errorf("StringMatchesAll() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestStringNotMatchesAll(t *testing.T) {
	type args struct {
		text  interface{}
		exprs interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "test not matches all",
			args: args{
				text:  "hello world",
				exprs: []string{"goodbye", "farewell"},
			},
			want: true,
		},
		{
			name: "test matches all",
			args: args{
				text:  "hello world",
				exprs: []string{"hello", "world"},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := assert.StringNotMatchesAll(tt.args.text, tt.args.exprs); got != tt.want {
				t.Errorf("StringNotMatchesAll() = %v, want %v", got, tt.want)
			}
		})
	}
}
