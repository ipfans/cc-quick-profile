//go:build darwin

package autostart

import (
	"fmt"
	"os"
	"path/filepath"
)

// darwinManager implements auto-start functionality for macOS using Launch Agents
type darwinManager struct {
	appName        string
	executablePath string
	plistPath      string
}

// newManager creates a new macOS auto-start manager
func newManager(appName, executablePath string) (Manager, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get home directory: %w", err)
	}

	launchAgentsDir := filepath.Join(homeDir, "Library", "LaunchAgents")
	plistPath := filepath.Join(launchAgentsDir, fmt.Sprintf("com.ipfans.%s.plist", appName))

	// Ensure LaunchAgents directory exists
	if err := os.MkdirAll(launchAgentsDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create LaunchAgents directory: %w", err)
	}

	return &darwinManager{
		appName:        appName,
		executablePath: executablePath,
		plistPath:      plistPath,
	}, nil
}

// IsEnabled checks if the plist file exists
func (d *darwinManager) IsEnabled() (bool, error) {
	_, err := os.Stat(d.plistPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check plist file: %w", err)
	}
	return true, nil
}

// Enable creates a plist file in LaunchAgents directory
func (d *darwinManager) Enable() error {
	plistContent := fmt.Sprintf(`<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>Label</key>
	<string>com.ipfans.%s</string>
	<key>ProgramArguments</key>
	<array>
		<string>%s</string>
	</array>
	<key>RunAtLoad</key>
	<true/>
	<key>KeepAlive</key>
	<false/>
	<key>StandardOutPath</key>
	<string>/dev/null</string>
	<key>StandardErrorPath</key>
	<string>/dev/null</string>
</dict>
</plist>`, d.appName, d.executablePath)

	if err := os.WriteFile(d.plistPath, []byte(plistContent), 0644); err != nil {
		return fmt.Errorf("failed to write plist file: %w", err)
	}

	return nil
}

// Disable removes the plist file from LaunchAgents directory
func (d *darwinManager) Disable() error {
	err := os.Remove(d.plistPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove plist file: %w", err)
	}
	return nil
}
