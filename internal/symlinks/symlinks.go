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

	var created int
	var skippedExisting int
	var alreadyCorrect int
	var aborted bool

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
						// Destination exists: handle symlink collisions separately so we don't prompt twice.
						if SymlinkExists(destination) {
							if skip, didAbort := SymlinkExistsHandler(source, destination, flags); didAbort {
								aborted = true
								break
							} else if skip {
								if SymlinkPointsToSource(destination, source) {
									alreadyCorrect++
								} else {
									skippedExisting++
								}
								continue
							}
						} else {
							if skip, didAbort := DestinationExistsHandler(source, destination, flags); didAbort {
								aborted = true
								break
							} else if skip {
								skippedExisting++
								continue
							}
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
					created++
					log.Info("Symlink created", "source", source, "destination", destination)
				}
			}
			if aborted {
				break
			}
		}
		if aborted {
			break
		}
	}

	if aborted {
		log.Warn("Symlink run aborted", "created", created, "skipped_existing", skippedExisting, "already_correct", alreadyCorrect)
		return
	}
	log.Info("Symlink run finished", "created", created, "skipped_existing", skippedExisting, "already_correct", alreadyCorrect)
}
