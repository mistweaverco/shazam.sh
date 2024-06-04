package symlinks

import (
	"os"
	"testing"
)

func TestPathMatchesSymlink(t *testing.T) {
	var userName = os.Getenv("USER")
	var ps = string(os.PathSeparator)
	type testInput = struct {
		rootName   string
		nodeName   string
		fileSource string
		flagPath   string
	}
	var tests = []struct {
		name  string
		input testInput
		want  bool
	}{
		{
			"should match when flagPath is equal to the source",
			testInput{
				rootName:   "configurations",
				nodeName:   "neovim",
				fileSource: "nvim",
				flagPath:   "configurations" + ps + "neovim" + ps + "nvim",
			},
			true,
		},
		{
			"should match when flagPath is a subset of the source",
			testInput{
				rootName:   "configurations",
				nodeName:   "neovim",
				fileSource: "nvim",
				flagPath:   "configurations" + ps + "neovim",
			},
			true,
		},
		{
			"should match when rootName is a user expansion",
			testInput{
				rootName:   "configurations-$USER",
				nodeName:   "neovim",
				fileSource: "nvim",
				flagPath:   "configurations-" + userName + ps + "neovim",
			},
			true,
		},
		{
			"should match when rootName is a user expansion and flagPath is a subset of the source",
			testInput{
				rootName:   "configurations-$USER",
				nodeName:   "neovim",
				fileSource: "nvim",
				flagPath:   "configurations-$USER" + ps + "neovim",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := PathMatchesSymlink(tt.input.rootName, tt.input.nodeName, tt.input.fileSource, tt.input.flagPath)
			if res != tt.want {
				t.Errorf("got %t, want %t", res, tt.want)
			}
		})
	}
}