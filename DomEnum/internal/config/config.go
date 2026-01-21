package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

// We declare the variable here, but we will initialize it properly in a function
var configFile string

type Config struct {
	Format map[string][]string `yaml:"config"`
}

func ReadConfig() (*Config, error) {
	// Initialize the path inside the function
	configFile = filepath.Join(os.Getenv("F"), "internal", "config", "config.yaml")

	dir := filepath.Dir(configFile)
	
	// 1. Create directory if it doesn't exist
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return nil, err
	}

	// 2. Create default file if it doesn't exist
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		defaultConfig := []byte("config: {}\n")
		if err := os.WriteFile(configFile, defaultConfig, 0644); err != nil {
			return nil, err
		}
	}

	// 3. Read file
	dat, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	// 4. Unmarshal YAML
	var cfg Config
	if err := yaml.Unmarshal(dat, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}