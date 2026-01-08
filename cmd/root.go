// Package cmd holds all commands avaialable.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "sysup",
	Short: "A CLI tool to automate the installation of programs on Fedora",
	Long:  `sysup is a CLI tool designed for Fedora users who frequently reinstall their OS. It automates the installation of system packages via DNF and desktop applications via Flatpak based on a YAML configuration file.`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
}
