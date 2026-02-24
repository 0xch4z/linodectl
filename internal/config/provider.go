package config

//go:generate go run go.uber.org/mock/mockgen -destination mock/mock_provider.go -package mock github.com/0xch4z/linodectl/internal/config Provider

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const fileName = ".linodectl"

// Provider manages the config file.
type Provider interface {
	Load() (*Config, error)
	Save(config *Config) error
}

type provider struct{}

// NewProvider builds a new config provider.
func NewProvider() Provider {
	return provider{}
}

func getPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home dir: %w", err)
	}
	return filepath.Join(homeDir, fileName), nil
}

// Save saves the config to disk.
func (provider) Save(config *Config) error {
	configBytes, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	configPath, err := getPath()
	if err != nil {
		return err
	}

	// open the file for writing, create it if it doesn't exist
	file, err := os.OpenFile(configPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		return err
	}

	_, err = file.Write(configBytes)
	return err
}

// Load loads the config from disk. If the config does not exist, a default one will
// be created.
func (p provider) Load() (*Config, error) {
	configPath, err := getPath()
	if err != nil {
		return nil, err
	}

	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			// the config file does not exist; create it
			config := DefaultConfig()
			return config, p.Save(config)
		}
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(configBytes, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return &config, nil
}

// DefaultConfig returns a default, minimal config.
func DefaultConfig() *Config {
	return &Config{
		Profile: "default",
		Profiles: map[string]Profile{
			"default": {
				APIVersion: "v4",
				APIBaseURL: "api.linode.com",
			},
		},
	}
}
