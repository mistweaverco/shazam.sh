package symlinks

import (
	"os"

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

	for rootName := range cfg {
		for _, node := range cfg[rootName] {
			if flags.Root != "" && flags.Root != rootName {
				continue
			}
			if flags.Only != "" && flags.Only != node.Name {
				continue
			}
			for _, file := range node.Files {
				if flags.Path != "" && !PathMatchesSymlink(rootName, node.Name, file.Source, flagPath) {
					continue
				}
				source, err := GetExpandedSource(dotfilesPath, rootName, node.Name, file.Source)
				if err != nil {
					continue
				}
				if _, err := os.Lstat(source); err != nil {
					log.Warn("Source file does not exist", "source", source)
					continue
				}
				destination, err := GetExpandedDestination(file.Destination)
				if err != nil {
					continue
				}
				if _, err := os.Lstat(destination); err == nil {
					if flags.PullInExisting {
						if flags.DryRun {
							log.Info("Dry run", "destination, would be pulled in", destination, "to", source)
							continue
						} else {
							if !SymlinkPointsToSource(destination, source) {
								err = os.Rename(destination, source)
								if err != nil {
									log.Error("Error moving existing path", "source", destination, "destination", source, "error", err)
									continue
								}
							}
						}
					} else {
						// Destination probably exists, check if it's a symlink
						if SymlinkExistsHandler(source, destination, flags) {
							continue
						} else if DestinationExistsHandler(destination, flags) {
							continue
						}
					}
				}
				if flags.DryRun {
					log.Info("Dry run", "source", source, "destination", destination)
				} else {
					err = os.Symlink(source, destination)
					if err != nil {
						if SymlinkCreationErrorHandler(source, destination, flags) {
							continue
						}
					}
					log.Info("Symlink created", "source", source, "destination", destination)
				}
			}
		}
	}
}