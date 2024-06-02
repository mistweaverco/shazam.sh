package shell

import (
	"os"

	"github.com/charmbracelet/log"
)

func GetCurrentShell() string {
	shell := os.Getenv("SHELL")
	if shell == "" {
		log.Warn("SHELL environment variable not set")
		return "sh"
	}
	return shell
}