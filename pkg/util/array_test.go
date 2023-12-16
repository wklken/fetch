package util_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/wklken/fetch/pkg/util"
)

func TestStringArrayToLower(t *testing.T) {
	type args struct {
		ss []string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{
			name: "test to lower",
			args: args{
				ss: []string{"HELLO", "WORLD"},
			},
			want: []string{"hello", "world"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := util.StringArrayToLower(tt.args.ss); !assert.Equal(t, got, tt.want) {
				t.Errorf("ToLower() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestItemInIntArray(t *testing.T) {
	type args struct {
		item  int
		array []int
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "test item in int array",
			args: args{
				item:  1,
				array: []int{1, 2, 3},
			},
			want: true,
		},
		{
			name: "test item not in int array",
			args: args{
				item:  4,
				array: []int{1, 2, 3},
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := util.ItemInIntArray(tt.args.item, tt.args.array); got != tt.want {
				t.Errorf("ItemInIntArray() = %v, want %v", got, tt.want)
			}
		})
	}
}
