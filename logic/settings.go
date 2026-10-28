package logic

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// GetClaudeSettingsPath returns the path to Claude Code's user settings file
func GetClaudeSettingsPath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error getting home directory: %v\n", err)
		os.Exit(1)
	}
	return filepath.Join(homeDir, ".claude", "settings.json")
}

// ClaudeSettings is a Claude Code settings document; top-level fields other than env are kept verbatim
type ClaudeSettings map[string]json.RawMessage

// LoadClaudeSettings loads the settings file, returning an empty document when it does not exist
func LoadClaudeSettings() (ClaudeSettings, error) {
	path := GetClaudeSettingsPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return ClaudeSettings{}, nil
		}
		return nil, fmt.Errorf("failed to read settings file: %v", err)
	}

	settings := ClaudeSettings{}
	if len(data) == 0 {
		return settings, nil
	}
	if err := json.Unmarshal(data, &settings); err != nil {
		return nil, fmt.Errorf("failed to parse %s: %v", path, err)
	}
	return settings, nil
}

// Env returns the environment variables defined in the settings file
func (s ClaudeSettings) Env() (map[string]string, error) {
	raw, exists := s["env"]
	if !exists {
		return map[string]string{}, nil
	}

	var env map[string]string
	if err := json.Unmarshal(raw, &env); err != nil {
		return nil, fmt.Errorf("failed to parse env in settings file: %v", err)
	}
	if env == nil {
		env = map[string]string{}
	}
	return env, nil
}

// SetEnv replaces the env section; an empty env removes the section
func (s ClaudeSettings) SetEnv(env map[string]string) error {
	if len(env) == 0 {
		delete(s, "env")
		return nil
	}

	raw, err := json.Marshal(env)
	if err != nil {
		return fmt.Errorf("failed to marshal env: %v", err)
	}
	s["env"] = raw
	return nil
}

// SaveClaudeSettings writes the settings file atomically, keeping the existing file mode
func SaveClaudeSettings(settings ClaudeSettings) error {
	path := GetClaudeSettingsPath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create settings directory: %v", err)
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal settings: %v", err)
	}
	data = append(data, '\n')

	mode := os.FileMode(0600)
	if info, statErr := os.Stat(path); statErr == nil {
		mode = info.Mode().Perm()
	}

	tmpFile, err := os.CreateTemp(dir, ".settings-*.json")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %v", err)
	}
	tmpPath := tmpFile.Name()
	defer os.Remove(tmpPath)

	if _, err := tmpFile.Write(data); err != nil {
		tmpFile.Close()
		return fmt.Errorf("failed to write temp file: %v", err)
	}
	if err := tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %v", err)
	}
	if err := os.Chmod(tmpPath, mode); err != nil {
		return fmt.Errorf("failed to set file mode: %v", err)
	}
	if err := os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("failed to replace settings file: %v", err)
	}
	return nil
}
