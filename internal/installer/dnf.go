package installer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func SetupDNFRepositories(repos []string) error {
	if len(repos) == 0 {
		return nil
	}

	// Check if dnf-plugins-core is installed
	checkCmd := exec.Command("rpm", "-q", "dnf-plugins-core")
	if err := checkCmd.Run(); err != nil {
		fmt.Println("dnf-plugins-core not found. Installing...")
		if err := InstallDNF([]string{"dnf-plugins-core"}); err != nil {
			return fmt.Errorf("failed to install dnf-plugins-core: %v", err)
		}
	}

	for _, repo := range repos {
		// Check if repository already exists to ensure idempotency
		// repoName := repo
		// if strings.HasSuffix(repo, ".repo") {
		// 	parts := strings.Split(repo, "/")
		// 	repoName = parts[len(parts)-1]
		// }

		// fmt.Printf("Checking repository: %s\n", repo)
		// // We use dnf repolist to see if it's already known
		// // This is a bit slow but reliable
		// checkRepo := exec.Command("dnf", "repolist", "--enabled")
		// output, _ := checkRepo.CombinedOutput()
		//
		// // If the repo URL or the file name is already in the repolist, skip
		// if strings.Contains(string(output), repoName) || strings.Contains(string(output), repo) {
		// 	fmt.Printf("Repository %s already enabled, skipping.\n", repo)
		// 	continue
		// }

		// fmt.Printf("Adding repository: %s\n", repo)
		cmd := exec.Command("sudo", "dnf", "config-manager", "addrepo", "--from-repofile="+repo, "--overwrite")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			return fmt.Errorf("failed to add repository %s: %v", repo, err)
		}
	}
	return nil
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
