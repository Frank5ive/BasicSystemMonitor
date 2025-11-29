package main

import (
	"log"
	"os"
	"time"

	"gopkg.in/yaml.v2"
)

// Config holds the application's configuration settings.
type Config struct {
	RefreshInterval        string `yaml:"refreshInterval"`
	DiskPath               string `yaml:"diskPath"`
	ProcessRefreshInterval string `yaml:"processRefreshInterval"`
}

// DefaultConfig returns a Config struct with default values.
func DefaultConfig() Config {
	return Config{
		RefreshInterval:        "1s",
		DiskPath:               "/",
		ProcessRefreshInterval: "3s",
	}
}

// LoadConfig loads configuration from the specified file path.
// If the file doesn't exist or there's an error, it returns a default configuration.
func LoadConfig(configPath string) (Config, error) {
	config := DefaultConfig()

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			log.Printf("Config file '%s' not found, using default configuration.", configPath)
			return config, nil
		}
		return config, err
	}

	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// GetRefreshInterval parses the RefreshInterval string into a time.Duration.
func (c *Config) GetRefreshInterval() (time.Duration, error) {
	return time.ParseDuration(c.RefreshInterval)
}

// GetProcessRefreshInterval parses the ProcessRefreshInterval string into a time.Duration.
func (c *Config) GetProcessRefreshInterval() (time.Duration, error) {
	return time.ParseDuration(c.ProcessRefreshInterval)
}
