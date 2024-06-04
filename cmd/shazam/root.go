package shazam

import (
	"os"

	"github.com/charmbracelet/log"
	"github.com/mistweaverco/shazam.sh/internal/config"
	"github.com/mistweaverco/shazam.sh/internal/symlinks"
	"github.com/spf13/cobra"
)

var cfg = config.NewConfig(config.Config{
	DataReader: os.ReadFile,
})

var rootCmd = &cobra.Command{
	Use:   "shazam",
	Short: "Dotfiles 🗃️ manager on steroids ⚡.",
	Long:  "Dotfiles 🗃️ manager on steroids ⚡. Makes managing your dotfiles a breeze.",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Info("Starting shazam.sh 🚀", "version", VERSION)
			symlinks.CreateSymlinks(cfg.GetConfigFile(), cfg.GetConfigFlags())
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
	rootCmd.PersistentFlags().StringVar(&cfg.ConfigPath, "config", "shazam.yml", "config file")
	rootCmd.PersistentFlags().BoolVar(&cfg.Flags.DryRun, "dry-run", false, "dry run")
	rootCmd.PersistentFlags().BoolVar(&cfg.Flags.PullInExisting, "pull-in-existing", false, "pull in existing files")
	rootCmd.PersistentFlags().StringVar(&cfg.Flags.Root, "root", "", "root workspace")
	rootCmd.PersistentFlags().StringVar(&cfg.Flags.Only, "only", "", "only specific nodes matching a name")
	rootCmd.PersistentFlags().StringVar(&cfg.Flags.DotfilesPath, "dotfiles-path", "", "dotfiles path")
	rootCmd.PersistentFlags().StringVar(&cfg.Flags.Path, "path", "", "path to config file or dir")
}