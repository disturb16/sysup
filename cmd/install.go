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
	onlyFlatpak  bool
	repositories bool
)

// installCmd represents the install command
var installCmd = &cobra.Command{
	Use:   "install",
	Short: "Install programs from the configuration file",
	Long:  `This command reads the configuration file and installs listed DNF and Flatpak programs. By default, it installs everything.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			log.Fatalf("Failed to load config: %v", err)
		}

		// If no specific flags are set, install everything
		installAll := !onlyDNF && !onlyFlatpak && !repositories

		if repositories || installAll {
			if err := installer.SetupRepositories(cfg.Repositories); err != nil {
				log.Fatalf("Failed to setup repositories: %v", err)
			}
		}

		// Setup Flatpak remotes
		if repositories || installAll || onlyFlatpak {
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
	rootCmd.AddCommand(installCmd)

	installCmd.Flags().StringVarP(&configPath, "config", "c", "config.yaml", "Path to the configuration file")
	installCmd.Flags().BoolVar(&onlyDNF, "dnf", false, "Only install DNF packages")
	installCmd.Flags().BoolVar(&onlyFlatpak, "flatpak", false, "Only install Flatpak apps")
	installCmd.Flags().BoolVar(&repositories, "repos", false, "Only setup repositories")
}
