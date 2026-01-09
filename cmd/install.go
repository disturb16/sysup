package cmd

import (
	"log"

	"github.com/disturb16/sysup/internal/config"
	"github.com/disturb16/sysup/internal/installer"
	"github.com/spf13/cobra"
)

var (
	configPath   string
	onlyDNF      bool
	onlyAPT      bool
	onlyFlatpak  bool
	repositories bool
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install programs from the configuration file",
	Long:  `Parent command for installation operations. Use subcommands 'fedora' or 'debian' to target specific distributions.`,
}

func setupFlatpakRemotes(cfg *config.Config) {
	remotes := cfg.FlatpakRemotes
	if len(remotes) == 0 {
		remotes = []config.FlatpakRemote{
			{Name: "flathub", Url: "https://flathub.org/repo/flathub.flatpakrepo"},
		}
	}

	if err := installer.SetupFlatpakRemotes(remotes); err != nil {
		log.Fatalf("Failed to setup flatpak remotes: %v", err)
	}
}

func init() {
	rootCmd.AddCommand(installCmd)

	installCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yaml", "Path to the configuration file")
}
