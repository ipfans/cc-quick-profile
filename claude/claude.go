// Package claude provides management for Claude Code settings.json file
package claude

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/ipfans/cc-quick-profile/assets"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

// Manager handles Claude settings.json file operations
type Manager struct {
	settingsPath string
}

// NewManager creates a new Claude settings manager
func NewManager() (*Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	claudeDir := filepath.Join(homeDir, ".claude")
	settingsPath := filepath.Join(claudeDir, "settings.json")

	m := &Manager{
		settingsPath: settingsPath,
	}

	// Ensure Claude directory exists and settings file is initialized
	if err := m.ensureSettingsFile(); err != nil {
		return nil, fmt.Errorf("failed to initialize Claude settings: %w", err)
	}

	return m, nil
}

// ensureSettingsFile creates the Claude settings file if it doesn't exist
func (m *Manager) ensureSettingsFile() error {
	// Create .claude directory if it doesn't exist
	claudeDir := filepath.Dir(m.settingsPath)
	if err := os.MkdirAll(claudeDir, 0755); err != nil {
		return fmt.Errorf("failed to create Claude directory: %w", err)
	}

	// Check if settings file exists
	if _, err := os.Stat(m.settingsPath); os.IsNotExist(err) {
		// Create default settings file using embedded template
		if err := os.WriteFile(m.settingsPath, assets.ClaudeDefaultSettings, 0644); err != nil {
			return fmt.Errorf("failed to create default settings file: %w", err)
		}
	}

	return nil
}

// HasAuthToken checks if ANTHROPIC_AUTH_TOKEN exists in env
func (m *Manager) HasAuthToken() (bool, error) {
	data, err := os.ReadFile(m.settingsPath)
	if err != nil {
		return false, fmt.Errorf("failed to read settings file: %w", err)
	}

	token := gjson.GetBytes(data, "env.ANTHROPIC_AUTH_TOKEN")
	return token.Exists(), nil
}

// HasBaseURL checks if ANTHROPIC_BASE_URL exists in env
func (m *Manager) HasBaseURL() (bool, error) {
	data, err := os.ReadFile(m.settingsPath)
	if err != nil {
		return false, fmt.Errorf("failed to read settings file: %w", err)
	}

	baseURL := gjson.GetBytes(data, "env.ANTHROPIC_BASE_URL")
	return baseURL.Exists(), nil
}

// HasAuthConfig checks if both ANTHROPIC_AUTH_TOKEN and ANTHROPIC_BASE_URL exist
func (m *Manager) HasAuthConfig() (bool, error) {
	hasToken, err := m.HasAuthToken()
	if err != nil {
		return false, err
	}

	hasURL, err := m.HasBaseURL()
	if err != nil {
		return false, err
	}

	return hasToken && hasURL, nil
}

// SetAuthConfig sets both ANTHROPIC_AUTH_TOKEN and ANTHROPIC_BASE_URL
func (m *Manager) SetAuthConfig(apiKey, apiURL string) error {
	data, err := os.ReadFile(m.settingsPath)
	if err != nil {
		return fmt.Errorf("failed to read settings file: %w", err)
	}

	// Set ANTHROPIC_AUTH_TOKEN
	updatedData, err := sjson.SetBytes(data, "env.ANTHROPIC_AUTH_TOKEN", apiKey)
	if err != nil {
		return fmt.Errorf("failed to set auth token: %w", err)
	}

	// Set ANTHROPIC_BASE_URL
	updatedData, err = sjson.SetBytes(updatedData, "env.ANTHROPIC_BASE_URL", apiURL)
	if err != nil {
		return fmt.Errorf("failed to set base URL: %w", err)
	}

	// Write back to file
	if err := os.WriteFile(m.settingsPath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write settings file: %w", err)
	}

	return nil
}

// RemoveAuthConfig removes both ANTHROPIC_AUTH_TOKEN and ANTHROPIC_BASE_URL
func (m *Manager) RemoveAuthConfig() error {
	data, err := os.ReadFile(m.settingsPath)
	if err != nil {
		return fmt.Errorf("failed to read settings file: %w", err)
	}

	// Remove ANTHROPIC_AUTH_TOKEN
	updatedData, err := sjson.DeleteBytes(data, "env.ANTHROPIC_AUTH_TOKEN")
	if err != nil {
		return fmt.Errorf("failed to remove auth token: %w", err)
	}

	// Remove ANTHROPIC_BASE_URL
	updatedData, err = sjson.DeleteBytes(updatedData, "env.ANTHROPIC_BASE_URL")
	if err != nil {
		return fmt.Errorf("failed to remove base URL: %w", err)
	}

	// Write back to file
	if err := os.WriteFile(m.settingsPath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write settings file: %w", err)
	}

	return nil
}
