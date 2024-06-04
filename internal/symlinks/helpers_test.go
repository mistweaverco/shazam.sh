package symlinks

import (
	"os"
	"testing"
)

func TestGetExpandedSource(t *testing.T) {
	var baseDir, _ = os.Getwd()
	var userHome = os.Getenv("HOME")
	var userName = os.Getenv("USER")
	var ps = string(os.PathSeparator)
	type testInput = struct {
		dotfilesPath string
		rootName     string
		nodeName     string
		fileSource   string
	}
	var tests = []struct {
		name  string
		input testInput
		want  string
	}{
		{
			"should return the correct path without dotfilesPath and env vars",
			testInput{
				dotfilesPath: "",
				rootName:     "configurations",
				nodeName:     "neovim",
				fileSource:   "nvim",
			},
			baseDir + ps + "configurations" + ps + "neovim" + ps + "nvim",
		},
		{
			"should return the correct path with dotfilesPath but without env vars",
			testInput{
				dotfilesPath: "/home/user/dotfiles",
				rootName:     "configurations",
				nodeName:     "neovim",
				fileSource:   "nvim",
			},
			"/home/user/dotfiles" + ps + "configurations" + ps + "neovim" + ps + "nvim",
		},
		{
			"should return the correct path with dotfilesPath and env vars",
			testInput{
				dotfilesPath: "$HOME/dotfiles",
				rootName:     "configurations-$USER",
				nodeName:     "neovim",
				fileSource:   "nvim",
			},
			userHome + "/dotfiles" + ps + "configurations-" + userName + ps + "neovim" + ps + "nvim",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, _ := GetExpandedSource(tt.input.dotfilesPath, tt.input.rootName, tt.input.nodeName, tt.input.fileSource)
			if res != tt.want {
				t.Errorf("got %s, want %s", res, tt.want)
			} else {
				t.Logf("got %s, want %s", res, tt.want)
			}
		})
	}
}

func TestGetExpandedDestination(t *testing.T) {
	var userHome = os.Getenv("HOME")
	var tests = []struct {
		name  string
		input string
		want  string
	}{
		{
			"should return the correct path without env var",
			"/home/marco/.config/neovim/nvim",
			"/home/marco/.config/neovim/nvim",
		},
		{
			"should return the correct path with env var $HOME",
			"$HOME/.config/neovim/nvim",
			userHome + "/.config/neovim/nvim",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, _ := GetExpandedDestination(tt.input)
			if res != tt.want {
				t.Errorf("got %s, want %s", res, tt.want)
			}
		})
	}
}