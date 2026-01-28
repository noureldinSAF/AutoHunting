package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

/*
#File Syntax

config:
  shodan:
    - xxxxxxxxxxxx
    - zzzzzzzzzzzz
  whoisxmlapi:
    - yyyyyyyyyyyy
    - vvvvvvvvvvvv
*/

// Config: /root/.config/subenum/config.yaml
var configFile string = filepath.Join(os.Getenv("HOME"), ".config", "subenum", "config.yaml")

type Config struct {
	Format map[string][]string `yaml:"config"`
}

func ReadConfig() (*Config, error) {
	// 1. Ensure directory exists
	dir := filepath.Dir(configFile)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}
	// 2. Create file if it does not exist
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		defaultConfig := []byte("config: {}\n")
		if err := os.WriteFile(configFile, defaultConfig, os.ModePerm); err != nil {
			return nil, err
		}
	}

	// 3. Read file
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	// 4. Unmarshal YAML
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
