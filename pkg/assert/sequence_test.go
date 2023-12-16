package assert_test

import (
	"testing"

	"github.com/wklken/fetch/pkg/assert"
)

func TestContains(t *testing.T) {
	type args struct {
		s       interface{}
		element interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantMsg string
	}{
		{
			name: "test contains",
			args: args{
				s:       []int{1, 2, 3},
				element: 2,
			},
			want:    true,
			wantMsg: "OK",
		},
		{
			name: "test does not contain",
			args: args{
				s:       []int{1, 2, 3},
				element: 4,
			},
			want:    false,
			wantMsg: "contains | `[1, 2, 3]` does not contain`4`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, gotMsg := assert.Contains(tt.args.s, tt.args.element); got != tt.want || gotMsg != tt.wantMsg {
				t.Errorf("Contains() = (%v, %v), want (%v, %v)", got, gotMsg, tt.want, tt.wantMsg)
			}
		})
	}
}

func TestNotContains(t *testing.T) {
	type args struct {
		s       interface{}
		element interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantMsg string
	}{
		{
			name: "test not contains",
			args: args{
				s:       []int{1, 2, 3},
				element: 4,
			},
			want:    true,
			wantMsg: "OK",
		},
		{
			name: "test contains",
			args: args{
				s:       []int{1, 2, 3},
				element: 2,
			},
			want:    false,
			wantMsg: "not_contains | `[1, 2, 3]` should not contain `2`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, gotMsg := assert.NotContains(tt.args.s, tt.args.element); got != tt.want || gotMsg != tt.wantMsg {
				t.Errorf("NotContains() = (%v, %v), want (%v, %v)", got, gotMsg, tt.want, tt.wantMsg)
			}
		})
	}
}

func TestIn(t *testing.T) {
	type args struct {
		element interface{}
		s       interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantMsg string
	}{
		{
			name: "test in",
			args: args{
				element: 2,
				s:       []int{1, 2, 3},
			},
			want:    true,
			wantMsg: "OK",
		},
		{
			name: "test not in",
			args: args{
				element: 4,
				s:       []int{1, 2, 3},
			},
			want:    false,
			wantMsg: "in | `4` not in `[1, 2, 3]`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, gotMsg := assert.In(tt.args.element, tt.args.s); got != tt.want || gotMsg != tt.wantMsg {
				t.Errorf("In() = (%v, %v), want (%v, %v)", got, gotMsg, tt.want, tt.wantMsg)
			}
		})
	}
}

func TestNotIn(t *testing.T) {
	type args struct {
		element interface{}
		s       interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantMsg string
	}{
		{
			name: "test not in",
			args: args{
				element: 4,
				s:       []int{1, 2, 3},
			},
			want:    true,
			wantMsg: "OK",
		},
		{
			name: "test in",
			args: args{
				element: 2,
				s:       []int{1, 2, 3},
			},
			want:    false,
			wantMsg: "not_in | `2` should not in `[1, 2, 3]`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, gotMsg := assert.NotIn(tt.args.element, tt.args.s); got != tt.want || gotMsg != tt.wantMsg {
				t.Errorf("NotIn() = (%v, %v), want (%v, %v)", got, gotMsg, tt.want, tt.wantMsg)
			}
		})
	}
}
