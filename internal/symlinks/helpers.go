package symlinks

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
)

func GetExpandedSource(dotfilesPath string, rootName string, nodeName string, fileSource string) (string, error) {
	ps := string(os.PathSeparator)
	if dotfilesPath != "" {
		// Add path separator to the end of the path if it's not there
		if dotfilesPath[len(dotfilesPath)-1:] != ps {
			dotfilesPath += ps
		}
	}
	source, err := filepath.Abs(os.ExpandEnv(dotfilesPath + rootName + ps + nodeName + ps + fileSource))
	if err != nil {
		log.Error("Error getting absolute path", "source", fileSource, "error", err)
	}
	return source, err
}

func GetExpandedDestination(destination string) (string, error) {
	destination, err := filepath.Abs(os.ExpandEnv(destination))
	if err != nil {
		log.Error("Error getting absolute path", "destination", destination, "error", err)
	}
	return destination, err
}