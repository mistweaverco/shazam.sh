package shazam

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/mistweaverco/shazam.sh/internal/config"
	"github.com/mistweaverco/shazam.sh/internal/symlinks"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "shazam",
	Short: "Dotfiles 🗃️ manager on steroids ⚡.",
	Long:  "Dotfiles 🗃️ manager on steroids ⚡. Makes managing your dotfiles a breeze.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Info("Starting shazam.sh 🚀", "version", VERSION)

			cfg := config.NewConfig(config.Config{
				DataReader: os.ReadFile,
			})
			symlinks.CreateSymlinks(cfg.File, config.Flags)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&config.ConfigPath, "config", "shazam.yml", "config file")
	rootCmd.PersistentFlags().BoolVar(&config.Flags.DryRun, "dry-run", false, "dry run")
	rootCmd.PersistentFlags().BoolVar(&config.Flags.PullInExisting, "pull-in-existing", false, "pull in existing files")
	rootCmd.PersistentFlags().StringVar(&config.Flags.Root, "root", "", "root workspace")
	rootCmd.PersistentFlags().StringVar(&config.Flags.Only, "only", "", "only specific nodes matching a name")
	rootCmd.PersistentFlags().StringVar(&config.Flags.DotfilesPath, "dotfiles-path", "", "dotfiles path")
	rootCmd.PersistentFlags().StringVar(&config.Flags.Path, "path", "", "path to config file or dir")
}