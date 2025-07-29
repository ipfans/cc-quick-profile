//go:build linux

package autostart

import (
	"fmt"
	"os"
	"path/filepath"
)

// linuxManager implements auto-start functionality for Linux using desktop files
type linuxManager struct {
	appName        string
	executablePath string
	desktopPath    string
}

// newManager creates a new Linux auto-start manager
func newManager(appName, executablePath string) (Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	autostartDir := filepath.Join(homeDir, ".config", "autostart")
	desktopPath := filepath.Join(autostartDir, fmt.Sprintf("%s.desktop", appName))

	// Ensure autostart directory exists
	if err := os.MkdirAll(autostartDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create autostart directory: %w", err)
	}

	return &linuxManager{
		appName:        appName,
		executablePath: executablePath,
		desktopPath:    desktopPath,
	}, nil
}

// IsEnabled checks if the desktop file exists
func (l *linuxManager) IsEnabled() (bool, error) {
	_, err := os.Stat(l.desktopPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check desktop file: %w", err)
	}
	return true, nil
}

// Enable creates a desktop file in autostart directory
func (l *linuxManager) Enable() error {
	desktopContent := fmt.Sprintf(`[Desktop Entry]
Type=Application
Name=%s
Exec=%s
Hidden=false
NoDisplay=false
X-GNOME-Autostart-enabled=true
Comment=CC Quick Profile - Claude Code profile switcher
`, l.appName, l.executablePath)

	if err := os.WriteFile(l.desktopPath, []byte(desktopContent), 0644); err != nil {
		return fmt.Errorf("failed to write desktop file: %w", err)
	}

	return nil
}

// Disable removes the desktop file from autostart directory
func (l *linuxManager) Disable() error {
	err := os.Remove(l.desktopPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove desktop file: %w", err)
	}
	return nil
}
