package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	DNF          []string `yaml:"dnf"`
	Flatpak      []string `yaml:"flatpak"`
	Repositories []string `yaml:"repositories"`
	PostInstall  []string `yaml:"post_install"`
}

func LoadConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
