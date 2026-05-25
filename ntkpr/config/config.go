package config

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	// program state storage
	StateFilePath string
	DataFilePath  string
}

func generateDefault() Config {
	dataFilePath := DataFilePathDefault()
	stateFilePath := StateFilePathDefault()

	cfg := Config{
		StateFilePath: stateFilePath,
		DataFilePath:  dataFilePath,
	}
	return cfg
}

func LoadOrCreateConfig() Config {
	path := ConfigPath()

	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		// Config file doesn't exist, create directory and config file with defaults
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating config directory: %v, using default\n", err)
			return generateDefault()
		}

		cfg := generateDefault()
		yamlData, err := yaml.Marshal(&cfg)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error marshaling default config: %v, using default\n", err)
			return cfg
		}

		if err := os.WriteFile(path, yamlData, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing config file: %v, using default\n", err)
			return cfg
		}

		return cfg
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading config file: %v, using default\n", err)
		return generateDefault()
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing config file: %v, using default\n", err)
		return generateDefault()
	}

	return cfg
}

// We need to provide handles for the users to change their configs.
