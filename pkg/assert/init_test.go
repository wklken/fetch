package assert

import (
	"reflect"
	"testing"
)

func TestPrettyLine(t *testing.T) {
	type args struct {
		s interface{}
	}
	tests := []struct {
		name string
		args args
		want interface{}
	}{
		{
			name: "test pretty line with string",
			args: args{
				s: "hello\nworld",
			},
			want: "hello\\nworld",
		},
		{
			name: "test pretty line with array",
			args: args{
				s: []int{1, 2, 3},
			},
			want: "[1, 2, 3]",
		},
		{
			name: "test pretty line with slice",
			args: args{
				s: []string{"hello", "world"},
			},
			want: `["hello", "world"]`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := prettyLine(tt.args.s); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("PrettyLine() = %v, want %v", got, tt.want)
			}
		})
	}
}
