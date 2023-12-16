package assert_test

import (
	"testing"

	"github.com/wklken/fetch/pkg/assert"
)

func TestEqual(t *testing.T) {
	type args struct {
		actual   interface{}
		expected interface{}
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "test equal",
			args: args{
				actual:   1,
				expected: 1,
			},
			want: true,
		},
		{
			name: "test not equal",
			args: args{
				actual:   1,
				expected: 2,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := assert.Equal(tt.args.actual, tt.args.expected); got != tt.want {
				t.Errorf("Equal() = %v, want %v", got, tt.want)
			}
		})
	}
}
