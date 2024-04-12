package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database        string `yaml:"database"`
	Host            string `yaml:"host"`
	SessionLifetime int    `yaml:"session_lifetime"` // in minutes
	CleanupInterval int    `yaml:"cleanup_interval"` // in minutes
}

func LoadConfig(path string) (*Config, error) {
	contents, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	config := new(Config)
	err = yaml.Unmarshal(contents, config)
	return config, err
}
