package installer

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func InstallAPT(packages []string) error {
	if len(packages) == 0 {
		return nil
	}

	if !commandExists("apt-get") {
		return fmt.Errorf("apt-get command not found")
	}

	fmt.Println("Updating APT repositories...")
	updateCmd := exec.Command("sudo", "apt-get", "update")
	updateCmd.Stdout = os.Stdout
	updateCmd.Stderr = os.Stderr
	if err := updateCmd.Run(); err != nil {
		return fmt.Errorf("failed to update apt: %v", err)
	}

	fmt.Printf("Installing APT packages: %s\n", strings.Join(packages, ", "))

	args := append([]string{"apt-get", "install", "-y"}, packages...)
	cmd := exec.Command("sudo", args...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	return cmd.Run()
}
