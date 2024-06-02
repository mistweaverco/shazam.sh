package symlinks

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/mistweaverco/shazam.sh/internal/config"
)

func CreateSymlinks(cfg config.ConfigFile, flags config.ConfigFlags) {
	var ps = string(os.PathSeparator)
	dotfilesPath := ""
	if flags.DotfilesPath != "" {
		dotfilesPath = flags.DotfilesPath + ps
	}
	flagPath := dotfilesPath + flags.Path

	log.Info("Creating symlinks 🔗")
	for rootName := range cfg {
		for _, node := range cfg[rootName] {
			if flags.Root != "" && flags.Root != rootName {
				continue
			}
			if flags.Only != "" && flags.Only != node.Name {
				continue
			}
			for _, file := range node.Files {
				expSource := os.ExpandEnv(rootName + ps + node.Name + ps + file.Source)
				if flags.Path != "" && !strings.HasPrefix(expSource, flagPath) {
					continue
				}
				source, err := filepath.Abs(dotfilesPath + expSource)
				if err != nil {
					log.Error("Error getting absolute path", "source", file.Source, "error", err)
					continue
				}
				if _, err := os.Lstat(source); err != nil {
					log.Warn("Source file does not exist", "source", source)
					continue
				}
				destination, err := filepath.Abs(os.ExpandEnv(file.Destination))
				if err != nil {
					log.Error("Error getting absolute path", "destination", file.Destination, "error", err)
					continue
				}
				if _, err := os.Lstat(destination); err == nil {
					if flags.PullInExisting && !flags.DryRun {
						err = os.Rename(destination, source)
						if err != nil {
							log.Error("Error moving existing path", "source", destination, "destination", source, "error", err)
							continue
						}
					} else {
						log.Warn("File already exists", "source", source, "destination", destination)
						continue
					}
				}
				if !flags.DryRun {
					err = os.Symlink(source, destination)
					if err != nil {
						log.Error("Error creating symlink", "source", source, "destination", destination, "error", err)
						continue
					}
					log.Info("Symlink created", "source", source, "destination", destination)
				} else {
					log.Info("Dry run", "source", source, "destination", destination)
				}
			}
		}
	}
}