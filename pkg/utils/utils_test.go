package utils

import "testing"

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
			"empty string",
			args{str: " \t\n"},
			true,
		},
		{
			"not empty string",
			args{str: "kk"},
			false,
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
