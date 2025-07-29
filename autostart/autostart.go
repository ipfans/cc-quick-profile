// Package autostart provides cross-platform functionality for managing application auto-start
package autostart

// Manager handles auto-start functionality across different platforms
type Manager interface {
	// IsEnabled checks if auto-start is currently enabled
	IsEnabled() (bool, error)
	// Enable enables auto-start for the application
	Enable() error
	// Disable disables auto-start for the application
	Disable() error
}

// NewManager creates a new platform-specific auto-start manager
func NewManager(appName, executablePath string) (Manager, error) {
	return newManager(appName, executablePath)
}
