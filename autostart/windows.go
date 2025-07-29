//go:build windows

package autostart

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"unsafe"
)

// windowsManager implements auto-start functionality for Windows using shortcuts in Startup folder
type windowsManager struct {
	appName        string
	executablePath string
	shortcutPath   string
}

// newManager creates a new Windows auto-start manager
func newManager(appName, executablePath string) (Manager, error) {
	// Get the Startup folder path
	startupDir, err := getStartupDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get startup directory: %w", err)
	}

	shortcutPath := filepath.Join(startupDir, appName+".lnk")

	return &windowsManager{
		appName:        appName,
		executablePath: executablePath,
		shortcutPath:   shortcutPath,
	}, nil
}

// IsEnabled checks if the shortcut exists in the Startup folder
func (w *windowsManager) IsEnabled() (bool, error) {
	_, err := os.Stat(w.shortcutPath)
	if err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, fmt.Errorf("failed to check shortcut: %w", err)
	}
	return true, nil
}

// Enable creates a shortcut in the Startup folder
func (w *windowsManager) Enable() error {
	return w.createShortcut()
}

// Disable removes the shortcut from the Startup folder
func (w *windowsManager) Disable() error {
	err := os.Remove(w.shortcutPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove shortcut: %w", err)
	}
	return nil
}

// getStartupDir returns the Windows Startup folder path
func getStartupDir() (string, error) {
	appData := os.Getenv("APPDATA")
	if appData == "" {
		return "", fmt.Errorf("APPDATA environment variable not set")
	}
	return filepath.Join(appData, "Microsoft", "Windows", "Start Menu", "Programs", "Startup"), nil
}

// createShortcut creates a Windows shortcut (.lnk file)
func (w *windowsManager) createShortcut() error {
	// Load required DLLs
	ole32 := syscall.NewLazyDLL("ole32.dll")
	coInitialize := ole32.NewProc("CoInitialize")
	coCreateInstance := ole32.NewProc("CoCreateInstance")
	coUninitialize := ole32.NewProc("CoUninitialize")

	// Initialize COM
	coInitialize.Call(0)
	defer coUninitialize.Call()

	// Create IShellLink instance
	var shellLink uintptr
	clsidShellLink := &syscall.GUID{0x00021401, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
	iidIShellLink := &syscall.GUID{0x000214F9, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}

	ret, _, _ := coCreateInstance.Call(
		uintptr(unsafe.Pointer(clsidShellLink)),
		0,
		1, // CLSCTX_INPROC_SERVER
		uintptr(unsafe.Pointer(iidIShellLink)),
		uintptr(unsafe.Pointer(&shellLink)),
	)

	if ret != 0 {
		return fmt.Errorf("failed to create IShellLink instance: %x", ret)
	}

	// Get the vtable
	vtable := *(**[16]uintptr)(unsafe.Pointer(shellLink))

	// Set path (index 20 in vtable for SetPath)
	execPathPtr, _ := syscall.UTF16PtrFromString(w.executablePath)
	syscall.Syscall(vtable[20], 2, shellLink, uintptr(unsafe.Pointer(execPathPtr)), 0)

	// Set working directory
	workDir := filepath.Dir(w.executablePath)
	workDirPtr, _ := syscall.UTF16PtrFromString(workDir)
	syscall.Syscall(vtable[9], 2, shellLink, uintptr(unsafe.Pointer(workDirPtr)), 0)

	// Query for IPersistFile interface
	var persistFile uintptr
	iidIPersistFile := &syscall.GUID{0x0000010B, 0x0000, 0x0000, [8]byte{0xC0, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x46}}
	syscall.Syscall(vtable[0], 3, shellLink, uintptr(unsafe.Pointer(iidIPersistFile)), uintptr(unsafe.Pointer(&persistFile)))

	// Save the shortcut
	shortcutPathPtr, _ := syscall.UTF16PtrFromString(w.shortcutPath)
	persistVtable := *(**[16]uintptr)(unsafe.Pointer(persistFile))
	ret, _, _ = syscall.Syscall(persistVtable[6], 3, persistFile, uintptr(unsafe.Pointer(shortcutPathPtr)), 1)

	// Release interfaces
	syscall.Syscall(persistVtable[2], 1, persistFile, 0, 0) // Release IPersistFile
	syscall.Syscall(vtable[2], 1, shellLink, 0, 0)          // Release IShellLink

	if ret != 0 {
		return fmt.Errorf("failed to save shortcut: %x", ret)
	}

	return nil
}
