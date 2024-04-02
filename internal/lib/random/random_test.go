package random

import (
	"testing"

	"github.com/go-playground/assert/v2"
)

func TestRandomString(t *testing.T) {
	tests := []struct {
		name string
		arg  int
		want int
	}{
		{
			name: "should return string len 5",
			arg:  5,
			want: 5,
		},
		{
			name: "should return string len 10",
			arg:  10,
			want: 10,
		},
		{
			name: "should return string len 0",
			arg:  0,
			want: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewRandomString(tt.arg)
			assert.Equal(t, tt.want, len(got))
		})
	}
}
