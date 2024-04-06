package app

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Database string `yaml:"database"`
	Host     string `yaml:"host"`
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
