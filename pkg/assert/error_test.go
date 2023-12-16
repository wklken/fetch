package assert_test

import (
	"fmt"
	"testing"

	"github.com/wklken/fetch/pkg/assert"
)

func TestNoError(t *testing.T) {
	type args struct {
		err error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		// TODO: Add test cases.
		{
			name: "test no error",
			args: args{
				err: nil,
			},
			want: true,
		},
		{
			name: "test with error",
			args: args{
				err: fmt.Errorf("some error"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got, _ := assert.NoError(tt.args.err); got != tt.want {
				t.Errorf("NoError() = %v, want %v", got, tt.want)
			}
		})
	}
}
