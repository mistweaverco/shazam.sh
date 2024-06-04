package shell

import (
	"os"
	"testing"
)

func TestGetCurrentShell(t *testing.T) {
	tests := []struct {
		got  string
		name string
		want string
	}{
		{
			"bash",
			"should return bash when the $SHELL environment variable is set to bash",
			"bash",
		},
		{
			"zsh",
			"should return zsh when the $SHELL environment variable is set to zsh",
			"zsh",
		},
		{
			"",
			"should return sh when the $SHELL environment variable is set to an empty string",
			"sh",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Setenv("SHELL", tt.got)
			got := GetCurrentShell()
			if got != tt.want {
				t.Errorf("GetCurrentShell() = %v, want %v", got, tt.want)
			}
		})
	}
}