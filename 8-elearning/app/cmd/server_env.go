package cmd

import (
	"elearning-api/config"
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

func NewConfig(filepath string) (*config.Config, error) {
	if filepath == "" {
		return nil, fmt.Errorf("config file path is required")
	}
	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}
	var config config.Config
	if err := yaml.Unmarshal(bytes, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file: %w", err)
	}

	return &config, nil
}
