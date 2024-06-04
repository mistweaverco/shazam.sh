package symlinks

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/mistweaverco/shazam.sh/internal/config"
)

func SymlinkCreationErrorHandler(source string, destination string, flags config.ConfigFlags) bool {
	// Check if the creation failed,
	// because the parent directory does not exist
	dir := filepath.Dir(destination)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if !flags.DryRun {
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				log.Error("Error creating parent directory", "directory", dir, "error", err)
				return true
			}
			// Retry symlink creation
			err = os.Symlink(source, destination)
			if err != nil {
				log.Error("Error creating symlink", "source", source, "destination", destination, "error", err)
				return true
			}
		} else {
			log.Info("Dry run", "parent directory does not exist, would be created", dir)
			return true
		}
	}
	return false
}

func DestinationExistsHandler(destination string, flags config.ConfigFlags) bool {
	if flags.DryRun {
		log.Warn("Destination exists, skipping", "destination", destination)
		return true
	} else {
		log.Info("Dry run", "destination exists, would be deleted", destination)
		return true
	}
}

func SymlinkExistsHandler(source string, destination string, flags config.ConfigFlags) bool {
	if link, err := os.Readlink(destination); err == nil {
		// Symlink already exists and points to the correct source
		if link == source {
			return true
		} else {
			// Symlink already exists, but points to a different source
			if flags.DryRun {
				log.Info("Dry run", "destination exists as symlink, would be deleted", destination)
				return true
			} else {
				log.Info("Destination exists as symlink, skipping", "destination", destination)
				return true
			}
		}
	}
	return false
}