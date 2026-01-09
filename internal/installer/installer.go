package installer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/disturb16/sysup/internal/config"
)

func SetupRepositories(repos []string) error {
	if len(repos) == 0 {
		return nil
	}

	// Check if dnf-plugins-core is installed
	fmt.Println("Checking for dnf-plugins-core...")
	checkCmd := exec.Command("rpm", "-q", "dnf-plugins-core")
	if err := checkCmd.Run(); err != nil {
		fmt.Println("dnf-plugins-core not found. Installing...")
		if err := InstallDNF([]string{"dnf-plugins-core"}); err != nil {
			return fmt.Errorf("failed to install dnf-plugins-core: %v", err)
		}
	}

	for _, repo := range repos {
		// Check if repository already exists to ensure idempotency
		repoName := repo
		if strings.HasSuffix(repo, ".repo") {
			parts := strings.Split(repo, "/")
			repoName = parts[len(parts)-1]
		}

		fmt.Printf("Checking repository: %s\n", repo)
		// We use dnf repolist to see if it's already known
		// This is a bit slow but reliable
		checkRepo := exec.Command("dnf", "repolist", "--enabled")
		output, _ := checkRepo.CombinedOutput()

		// If the repo URL or the file name is already in the repolist, skip
		if strings.Contains(string(output), repoName) || strings.Contains(string(output), repo) {
			fmt.Printf("Repository %s already enabled, skipping.\n", repo)
			continue
		}

		fmt.Printf("Adding repository: %s\n", repo)
		cmd := exec.Command("sudo", "dnf", "config-manager", "addrepo", "--from-repofile="+repo, "--overwrite")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add repository %s: %v", repo, err)
		}
	}
	return nil
}

func RunScripts(scripts []string) error {
	for _, script := range scripts {
		fmt.Printf("Running script: %s\n", script)
		cmd := exec.Command("sh", "-c", script)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("script failed: %v", err)
		}
	}
	return nil
}

func commandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

func InstallDNF(packages []string) error {
	if len(packages) == 0 {
		return nil
	}

	if !commandExists("dnf") {
		return fmt.Errorf("dnf command not found")
	}

	fmt.Printf("Installing DNF packages: %s\n", strings.Join(packages, ", "))

	args := append([]string{"dnf", "install", "--skip-unavailable"}, packages...)
	cmd := exec.Command("sudo", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}

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

func SetupFlatpakRemotes(remotes []config.FlatpakRemote) error {
	if len(remotes) == 0 {
		return nil
	}

	if !commandExists("flatpak") {
		fmt.Println("Flatpak not found. Installing flatpak via DNF...")
		if err := InstallDNF([]string{"flatpak"}); err != nil {
			return fmt.Errorf("failed to install flatpak: %v", err)
		}
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

	if !commandExists("flatpak") {
		fmt.Println("Flatpak not found. Installing flatpak via DNF...")
		if err := InstallDNF([]string{"flatpak"}); err != nil {
			return fmt.Errorf("failed to install flatpak: %v", err)
		}
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
