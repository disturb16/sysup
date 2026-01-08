package installer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
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

	args := append([]string{"install", "-y"}, packages...)
	cmd := exec.Command("sudo", append([]string{"dnf"}, args...)...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
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

	// Add flathub repository if not present
	_ = exec.Command("flatpak", "remote-add", "--if-not-exists", "flathub", "https://flathub.org/repo/flathub.flatpakrepo").Run()

	fmt.Printf("Installing Flatpak apps: %s\n", strings.Join(apps, ", "))

	args := []string{"install", "-y", "--or-update", "flathub"}
	args = append(args, apps...)

	cmd := exec.Command("flatpak", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to install flatpak apps: %v", err)
	}

	return nil
}
