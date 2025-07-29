// Package config provides configuration management for the cc-quick-profile application
package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ipfans/cc-quick-profile/autostart"
	"github.com/ipfans/cc-quick-profile/claude"
	"github.com/ipfans/cc-quick-profile/models"
)

// Manager handles loading and saving application settings
type Manager struct {
	configPath       string
	settings         *models.Settings
	claudeManager    *claude.Manager
	autostartManager autostart.Manager
}

// NewManager creates a new configuration manager
func NewManager() (*Manager, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("failed to get config path: %w", err)
	}

	// Initialize Claude manager
	claudeManager, err := claude.NewManager()
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Claude manager: %w", err)
	}

	// Initialize autostart manager
	executablePath, err := os.Executable()
	if err != nil {
		return nil, fmt.Errorf("failed to get executable path: %w", err)
	}

	autostartManager, err := autostart.NewManager("cc-quick-profile", executablePath)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize autostart manager: %w", err)
	}

	m := &Manager{
		configPath:       configPath,
		claudeManager:    claudeManager,
		autostartManager: autostartManager,
	}

	// Load existing config or create default
	if err := m.Load(); err != nil {
		// If config doesn't exist, create with defaults
		if os.IsNotExist(err) {
			m.settings = models.NewSettings()

			// Check if Claude has auth config to determine initial enabled state
			hasAuthConfig, err := claudeManager.HasAuthConfig()
			if err != nil {
				return nil, fmt.Errorf("failed to check Claude auth config: %w", err)
			}
			m.settings.Enabled = hasAuthConfig

			if err := m.Save(); err != nil {
				return nil, fmt.Errorf("failed to save default config: %w", err)
			}
		} else {
			return nil, fmt.Errorf("failed to load config: %w", err)
		}
	}

	// Always check Claude auth config and sync with current settings
	hasAuthConfig, err := claudeManager.HasAuthConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to check Claude auth config: %w", err)
	}

	// Check autostart status and sync with current settings
	autostartEnabled, err := autostartManager.IsEnabled()
	if err != nil {
		return nil, fmt.Errorf("failed to check autostart status: %w", err)
	}

	// If actual states don't match our settings, update and save
	needsSave := false
	if m.settings.Enabled != hasAuthConfig {
		m.settings.Enabled = hasAuthConfig
		needsSave = true
	}
	if m.settings.AutoStart != autostartEnabled {
		m.settings.AutoStart = autostartEnabled
		needsSave = true
	}

	if needsSave {
		if err := m.Save(); err != nil {
			return nil, fmt.Errorf("failed to save updated config: %w", err)
		}
	}

	return m, nil
}

// getConfigPath returns the platform-specific configuration file path
func getConfigPath() (string, error) {
	var configDir string

	switch runtime.GOOS {
	case "windows":
		configDir = os.Getenv("APPDATA")
		if configDir == "" {
			return "", fmt.Errorf("APPDATA environment variable not set")
		}
	case "darwin", "linux":
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		configDir = filepath.Join(homeDir, ".config")
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}

	appConfigDir := filepath.Join(configDir, "cc-quick-profile")

	// Ensure directory exists
	if err := os.MkdirAll(appConfigDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return filepath.Join(appConfigDir, "settings.json"), nil
}

// Load reads the configuration from disk
func (m *Manager) Load() error {
	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return err
	}

	settings := &models.Settings{}
	if err := json.Unmarshal(data, settings); err != nil {
		return fmt.Errorf("failed to parse config: %w", err)
	}

	m.settings = settings
	return nil
}

// Save writes the current configuration to disk
func (m *Manager) Save() error {
	data, err := json.MarshalIndent(m.settings, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(m.configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// GetSettings returns the current settings
func (m *Manager) GetSettings() *models.Settings {
	return m.settings
}

// AddProfile adds a new profile to the configuration
func (m *Manager) AddProfile(profile models.Profile) error {
	// Check if profile with same name already exists
	for _, p := range m.settings.Profiles {
		if p.Name == profile.Name {
			return fmt.Errorf("profile with name '%s' already exists", profile.Name)
		}
	}

	m.settings.Profiles = append(m.settings.Profiles, profile)
	return m.Save()
}

// DeleteProfile removes a profile by name
func (m *Manager) DeleteProfile(name string) error {
	profiles := []models.Profile{}
	found := false

	for _, p := range m.settings.Profiles {
		if p.Name != name {
			profiles = append(profiles, p)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("profile '%s' not found", name)
	}

	m.settings.Profiles = profiles
	return m.Save()
}

// UpdateProfile updates an existing profile
func (m *Manager) UpdateProfile(name string, updated models.Profile) error {
	found := false

	for i, p := range m.settings.Profiles {
		if p.Name == name {
			m.settings.Profiles[i] = updated
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("profile '%s' not found", name)
	}

	return m.Save()
}

// SetEnabled sets the global enabled state and updates Claude settings
func (m *Manager) SetEnabled(enabled bool) error {
	m.settings.Enabled = enabled

	if enabled {
		// If enabling and there's an active profile, apply it to Claude settings
		activeProfile := m.settings.GetActiveProfile()
		if activeProfile != nil {
			if err := m.claudeManager.SetAuthConfig(activeProfile.APIKey, activeProfile.APIURL); err != nil {
				return fmt.Errorf("failed to set Claude auth config: %w", err)
			}
		}
	} else {
		// If disabling, remove auth config from Claude settings
		if err := m.claudeManager.RemoveAuthConfig(); err != nil {
			return fmt.Errorf("failed to remove Claude auth config: %w", err)
		}
	}

	return m.Save()
}

// SetActiveProfile activates a profile by name and updates Claude settings
func (m *Manager) SetActiveProfile(name string) error {
	m.settings.SetActiveProfile(name)

	// If enabled, update Claude settings with the new active profile
	if m.settings.Enabled {
		activeProfile := m.settings.GetActiveProfile()
		if activeProfile != nil {
			if err := m.claudeManager.SetAuthConfig(activeProfile.APIKey, activeProfile.APIURL); err != nil {
				return fmt.Errorf("failed to set Claude auth config: %w", err)
			}
		}
	}

	return m.Save()
}

// SetAutoStart sets the auto-start state and updates system auto-start configuration
func (m *Manager) SetAutoStart(enabled bool) error {
	m.settings.AutoStart = enabled

	if enabled {
		if err := m.autostartManager.Enable(); err != nil {
			return fmt.Errorf("failed to enable autostart: %w", err)
		}
	} else {
		if err := m.autostartManager.Disable(); err != nil {
			return fmt.Errorf("failed to disable autostart: %w", err)
		}
	}

	return m.Save()
}
