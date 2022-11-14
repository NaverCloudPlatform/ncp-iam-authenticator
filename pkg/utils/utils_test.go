package utils

import (
	"testing"
)

func TestIsEmptyString(t *testing.T) {
	type args struct {
		str string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			name: "empty string",
			args: args{str: "\t \n "},
			want: true,
		},
		{
			name: "not empty string",
			args: args{str: "abracadabra"},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsEmptyString(tt.args.str); got != tt.want {
				t.Errorf("IsEmptyString() = %v, want %v", got, tt.want)
			}
		})
	}
}
