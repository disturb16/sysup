package cmd

import (
	"log"

	"github.com/disturb16/sysup/internal/config"
	"github.com/disturb16/sysup/internal/installer"
	"github.com/spf13/cobra"
)

var debianCmd = &cobra.Command{
	Use:   "debian",
	Short: "Install programs for Debian/Ubuntu (APT, Flatpak)",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		installAll := !onlyAPT && !onlyFlatpak

		// Setup Flatpak remotes
		if installAll || onlyFlatpak {
			setupFlatpakRemotes(cfg)
		}

		if onlyAPT || installAll {
			if err := installer.InstallAPT(cfg.APT); err != nil {
				log.Fatalf("APT installation failed: %v", err)
			}
		}

		if onlyFlatpak || installAll {
			if err := installer.InstallFlatpak(cfg.Flatpak); err != nil {
				log.Fatalf("Flatpak installation failed: %v", err)
			}
		}

		if installAll {
			if err := installer.RunScripts(cfg.PostInstall); err != nil {
				log.Fatalf("Post-install scripts failed: %v", err)
			}
		}
	},
}

func init() {
	installCmd.AddCommand(debianCmd)

	debianCmd.Flags().BoolVar(&onlyAPT, "apt", false, "Only install APT packages")
	debianCmd.Flags().BoolVar(&onlyFlatpak, "flatpak", false, "Only install Flatpak apps")
}
