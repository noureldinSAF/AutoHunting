package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

/*
# File Syntax

config:
  shodan:
    - xxxxxxxxxxxx
    - zzzzzzzzzzzz
  whoisxmlapi:
    - yyyyyyyyyyyy
    - vvvvvvvvvvvv
*/

// ReadConfig reads (and creates if missing) a config.yaml located next to the binary.
// Example final path: /workspaces/AutoHunting/SubEnum/cmd/subenum/config.yaml
func ReadConfig() (*Config, error) {
	configPath := filepath.Join(os.Getenv("F"), "config.yaml")

	// Ensure directory exists
	dir := filepath.Dir(configPath)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	// Create file if it does not exist with an empty 'config' mapping
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		defaultConfig := []byte("config: {}\n")
		if err := os.WriteFile(configPath, defaultConfig, os.ModePerm); err != nil {
			return nil, err
		}
	} else if err != nil {
		// Stat returned an error other than NotExist
		return nil, err
	}

	// Read file
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	// Unmarshal YAML
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// Config represents the structure of the config.yaml file.
type Config struct {
	// Top-level key is `config`, which maps service names to lists of keys.
	// Example: config: { shodan: [key1, key2], whoisxmlapi: [key3] }
	Config map[string][]string `yaml:"config"`
}
