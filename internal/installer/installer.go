package installer

import (
	"fmt"
	"os"
	"os/exec"
)

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
