package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := new(Config)

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return config, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return config, err
	}

	return config, nil
}
