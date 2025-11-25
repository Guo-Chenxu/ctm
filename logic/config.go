package logic

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}
	return filepath.Join(homeDir, ".ctm", "config.json")
}

// LoadConfig loads the configuration from the config file
func LoadConfig() (*Config, error) {
	configPath := GetConfigPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

// SaveConfig saves the configuration to the config file
func SaveConfig(config *Config) error {
	configPath := GetConfigPath()

	// Ensure directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// ValidateIndex checks if the given index is valid
func ValidateIndex(config *Config, index int) error {
	if index < 0 || index >= len(config.Providers) {
		return fmt.Errorf("invalid index %d. Must be between 0 and %d", index, len(config.Providers)-1)
	}
	return nil
}

// GetAllEnvKeys collects all unique environment variable keys from all providers
func GetAllEnvKeys(config *Config) []string {
	keySet := make(map[string]bool)
	for _, provider := range config.Providers {
		for key := range provider.Env {
			keySet[key] = true
		}
	}

	keys := make([]string, 0, len(keySet))
	for key := range keySet {
		keys = append(keys, key)
	}
	return keys
}
