package installer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/disturb16/sysup/internal/config"
)

func getInstalledFlatpaks() (map[string]bool, error) {
	cmd := exec.Command("flatpak", "list", "--app", "--columns=application")
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	installed := make(map[string]bool)
	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" {
			installed[trimmed] = true
		}
	}
	return installed, nil
}

func EnsureFlatpakInstalled() error {
	if commandExists("flatpak") {
		return nil
	}
	fmt.Println("Flatpak not found. Attempting to install...")
	if commandExists("dnf") {
		fmt.Println("Installing flatpak via DNF...")
		return InstallDNF([]string{"flatpak"})
	}
	if commandExists("apt-get") {
		fmt.Println("Installing flatpak via APT...")
		return InstallAPT([]string{"flatpak"})
	}
	return fmt.Errorf("flatpak not found and no supported package manager (dnf, apt) found")
}

func SetupFlatpakRemotes(remotes []config.FlatpakRemote) error {
	if len(remotes) == 0 {
		return nil
	}

	if err := EnsureFlatpakInstalled(); err != nil {
		return err
	}

	for _, remote := range remotes {
		fmt.Printf("Adding Flatpak remote: %s (%s)\n", remote.Name, remote.Url)
		cmd := exec.Command("flatpak", "remote-add", "--if-not-exists", remote.Name, remote.Url)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add flatpak remote %s: %v", remote.Name, err)
		}
	}
	return nil
}

func InstallFlatpak(apps []string) error {
	if len(apps) == 0 {
		return nil
	}

	if err := EnsureFlatpakInstalled(); err != nil {
		return err
	}

	installed, err := getInstalledFlatpaks()
	if err != nil {
		fmt.Printf("Warning: Failed to check installed flatpaks: %v\n", err)
		installed = make(map[string]bool)
	}

	var appsToInstall []string
	for _, app := range apps {
		// handle specific remote/app syntax like remote/app.id
		appName := app
		parts := strings.Split(app, "/")
		if len(parts) > 1 {
			appName = parts[len(parts)-1]
		}

		if installed[appName] {
			fmt.Printf("Flatpak app %s is already installed, skipping.\n", appName)
			continue
		}
		appsToInstall = append(appsToInstall, app)
	}

	if len(appsToInstall) == 0 {
		fmt.Println("All requested Flatpak apps are already installed.")
		return nil
	}

	fmt.Printf("Installing Flatpak apps: %s\n", strings.Join(appsToInstall, ", "))

	args := []string{"install", "-y", "--or-update"}
	args = append(args, appsToInstall...)

	cmd := exec.Command("flatpak", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install flatpak apps: %v", err)
	}

	return nil
}
