package models

// Profile represents a Claude Code profile configuration
type Profile struct {
	Name   string `json:"name"`   // Configuration name for menu display
	APIURL string `json:"apiUrl"` // API endpoint URL
	APIKey string `json:"apiKey"` // API authentication key
	Active bool   `json:"active"` // Whether this is the currently active profile
}

// Settings represents the application settings
type Settings struct {
	Enabled   bool      `json:"enabled"`   // Global enable/disable state
	AutoStart bool      `json:"autoStart"` // Auto-start on system boot
	Profiles  []Profile `json:"profiles"`  // List of configured profiles
}

// NewSettings creates a new Settings instance with default values
func NewSettings() *Settings {
	return &Settings{
		Enabled:   true,
		AutoStart: false,
		Profiles:  []Profile{},
	}
}

// GetActiveProfile returns the currently active profile, or nil if none
func (s *Settings) GetActiveProfile() *Profile {
	for i := range s.Profiles {
		if s.Profiles[i].Active {
			return &s.Profiles[i]
		}
	}
	return nil
}

// SetActiveProfile sets the active profile by name and deactivates others
func (s *Settings) SetActiveProfile(name string) {
	for i := range s.Profiles {
		s.Profiles[i].Active = s.Profiles[i].Name == name
	}
}
