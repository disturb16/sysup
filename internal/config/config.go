package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type FlatpakRemote struct {
	Name string `yaml:"name"`
	Url  string `yaml:"url"`
}

type Config struct {
	DNF            []string        `yaml:"dnf"`
	APT            []string        `yaml:"apt"`
	Flatpak        []string        `yaml:"flatpak"`
	FlatpakRemotes []FlatpakRemote `yaml:"flatpak_remotes"`
	Repositories   []string        `yaml:"repositories"`
	PostInstall    []string        `yaml:"post_install"`
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
