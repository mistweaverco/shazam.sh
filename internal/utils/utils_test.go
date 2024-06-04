package utils

import (
	"testing"
)

func TestGetOperatingSystem(t *testing.T) {
	tests := []struct {
		got  string
		name string
		want string
	}{
		{
			"darwin",
			"should return macOS when the operating system is darwin",
			"macOS",
		},
		{
			"linux",
			"should return Linux when the operating system is linux",
			"Linux",
		},
		{
			"windows",
			"should return Windows when the operating system is windows",
			"Windows",
		},
		{
			"freebsd",
			"should return Other when the operating system is not darwin, linux, or windows",
			"Other",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetOperatingSystem(tt.got); got != tt.want {
				t.Errorf("GetOperatingSystem() = %v, want %v", got, tt.want)
			}
		})
	}
}