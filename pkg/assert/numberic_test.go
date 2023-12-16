package assert_test

import (
	"testing"

	"github.com/wklken/fetch/pkg/assert"
)

func TestLessOrEqual(t *testing.T) {
	type args struct {
		e1 interface{}
		e2 interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    bool
		wantErr string
	}{
		// TODO: Add test cases.
		{
			name: "test less or equal",
			args: args{
				e1: 1,
				e2: 2,
			},
			want:    true,
			wantErr: "OK",
		},
		{
			name: "test not less or equal",
			args: args{
				e1: 2,
				e2: 1,
			},
			want:    false,
			wantErr: "less_or_equal | `2` is not less than or equal to `1`",
		},
		{
			name: "test different types",
			args: args{
				e1: 1,
				e2: "1",
			},
			want:    false,
			wantErr: "less or equal error, elements should be the same type",
		},
		{
			name: "test uncomparable types",
			args: args{
				e1: []int{1, 2, 3},
				e2: []int{1, 2, 3},
			},
			want:    false,
			wantErr: "Can not compare type `[]int`",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := assert.LessOrEqual(tt.args.e1, tt.args.e2)
			if got != tt.want {
				t.Errorf("LessOrEqual() got = %v, want %v", got, tt.want)
			}
			if gotErr != tt.wantErr {
				t.Errorf("LessOrEqual() gotErr = %v, want %v", gotErr, tt.wantErr)
			}
		})
	}
}
