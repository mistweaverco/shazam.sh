package symlinks

import (
	"os"
	"path/filepath"

	"github.com/charmbracelet/log"
	"github.com/mistweaverco/shazam.sh/internal/config"
	"github.com/mistweaverco/shazam.sh/internal/ui"
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

func DestinationExistsHandler(source string, destination string, flags config.ConfigFlags) (skip bool, aborted bool) {
	if flags.DryRun {
		log.Info("Dry run", "destination exists, would be skipped", destination)
		return true, false
	}

	// Non-interactive environments should keep the current flow.
	if !ui.IsTTY() {
		log.Info("Destination exists, skipping", "destination", destination)
		return true, false
	}

	resolution, err := ui.PromptExistingPathResolution(ui.ExistingPathPromptInput{
		Source:      source,
		Destination: destination,
	})
	if err != nil {
		log.Error("Error prompting for resolution", "destination", destination, "error", err)
		return true, false
	}

	switch resolution {
	case ui.ResolutionSkip:
		log.Info("Destination exists, skipping", "destination", destination)
		return true, false
	case ui.ResolutionOverwrite:
		if err := os.RemoveAll(destination); err != nil {
			log.Error("Error deleting existing destination", "destination", destination, "error", err)
			return true, false
		}
		log.Info("Deleted existing destination", "destination", destination)
		return false, false
	case ui.ResolutionAbort:
		log.Warn("Aborted by user", "destination", destination)
		return true, true
	default:
		log.Info("Destination exists, skipping", "destination", destination)
		return true, false
	}
}

func SymlinkExistsHandler(source string, destination string, flags config.ConfigFlags) (skip bool, aborted bool) {
	if link, err := os.Readlink(destination); err == nil {
		// Symlink already exists and points to the correct source
		if link == source {
			return true, false
		} else {
			// Symlink already exists, but points to a different source: treat like any other collision.
			return DestinationExistsHandler(source, destination, flags)
		}
	}
	return false, false
}
