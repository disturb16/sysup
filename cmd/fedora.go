package cmd

import (
	"log"

	"github.com/disturb16/sysup/internal/config"
	"github.com/disturb16/sysup/internal/installer"
	"github.com/spf13/cobra"
)

var fedoraCmd = &cobra.Command{
	Use:   "fedora",
	Short: "Install programs for Fedora (DNF, Flatpak)",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		installAll := !onlyDNF && !onlyFlatpak && !repositories

		if repositories || installAll {
			if err := installer.SetupDNFRepositories(cfg.Repositories); err != nil {
				log.Fatalf("Failed to setup DNF repositories: %v", err)
			}
		}

		// Setup Flatpak remotes
		if repositories || installAll || onlyFlatpak {
			setupFlatpakRemotes(cfg)
		}

		if onlyDNF || installAll {
			if err := installer.InstallDNF(cfg.DNF); err != nil {
				log.Fatalf("DNF installation failed: %v", err)
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
	installCmd.AddCommand(fedoraCmd)

	fedoraCmd.Flags().BoolVar(&onlyDNF, "dnf", false, "Only install DNF packages")
	fedoraCmd.Flags().BoolVar(&onlyFlatpak, "flatpak", false, "Only install Flatpak apps")
	fedoraCmd.Flags().BoolVar(&repositories, "repos", false, "Only setup repositories")
}
