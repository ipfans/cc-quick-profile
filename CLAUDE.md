# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a cross-platform system tray GUI application called "cc-quick-profile" that allows users to quickly switch between Claude Code profiles. The application is built using:
- **fyne.io/fyne/v2 v2.6.2**: Cross-platform GUI framework for Go
- **fyne.io/systray v1.11.0**: System tray functionality (via Fyne's bundled version)
- **Go 1.24.5**: Programming language
- **github.com/tidwall/gjson & sjson**: JSON manipulation for Claude settings
- **Task (Taskfile.dev)**: Build automation and task runner

## Configuration

The application manages two types of configuration:

### Application Settings
- **Windows**: `%APPDATA%\cc-quick-profile\settings.json`
- **macOS/Linux**: `$HOME/.config/cc-quick-profile/settings.json`

### Claude Code Settings
- **All platforms**: `$HOME/.claude/settings.json` (managed automatically)

## Build System

This project uses [Taskfile](https://taskfile.dev) for task automation. **IMPORTANT: Always use `task` commands for building, testing, and other development tasks.**

### Common Task Commands

```bash
# Show available tasks
task

# Build the application for current platform
task build

# Build and run the application
task run

# Build for all platforms
task build:all

# Run tests
task test

# Format code
task fmt

# Run linter
task lint

# Run all checks (fmt, lint, test)
task check

# Clean build artifacts
task clean

# Download and tidy dependencies
task deps
```

### Platform-Specific Builds

```bash
# Build for Windows
task build:windows

# Build for macOS
task build:darwin

# Build for Linux
task build:linux
```

## Architecture Notes

The application consists of several key components:

1. **System Tray Integration** (`main.go`):
   - Uses Fyne's desktop.App interface for system tray functionality
   - Implements dynamic menu generation based on saved profiles
   - Handles user interactions through event channels
   - Manages application lifecycle with hidden main window

2. **Configuration Management** (`config/config.go`):
   - Manages loading/saving of application settings.json
   - Handles platform-specific configuration paths
   - Provides CRUD operations for profiles
   - Integrates with Claude manager for settings synchronization

3. **Claude Integration** (`claude/claude.go`):
   - Manages Claude Code's `~/.claude/settings.json` file
   - Uses gjson/sjson for JSON manipulation
   - Handles ANTHROPIC_AUTH_TOKEN and ANTHROPIC_BASE_URL env variables
   - Ensures Claude directory and default settings exist

4. **Data Models** (`models/profile.go`):
   - Profile: Contains Name, APIURL, APIKey, and Active status
   - Settings: Contains global Enabled state and list of Profiles
   - Helper methods for active profile management

5. **UI Components** (`ui/modal.go`):
   - Fyne-based modal dialog for adding new profiles
   - Input validation for profile fields (name, URL format, API key)
   - Event communication with main application via channels

6. **Assets** (`assets/`):
   - Embedded application icon and Claude default settings template
   - Used for system tray icon and initial Claude settings creation

## Development Guidelines

- **Always use `task` commands** for building, testing, and running the application
- Follow standard Go project layout conventions
- Use Fyne's material design components for consistent UI
- Ensure proper error handling for file operations (settings.json and Claude settings)
- Implement graceful shutdown for the systray application
- Consider cross-platform compatibility (Windows, macOS, Linux)
- Run `task check` before committing to ensure code quality
- Use embedded assets for static resources (icons, templates)
- Maintain separation between application config and Claude settings management

## How It Works

1. **Application Startup**:
   - Initializes Fyne app and creates hidden main window
   - Sets up configuration manager and loads/creates settings
   - Checks Claude auth config and syncs with application enabled state
   - Creates system tray with dynamic menu based on profiles

2. **Profile Management**:
   - Users can add profiles via modal dialog with validation
   - Profiles are stored in application settings.json
   - Active profile selection updates both app settings and Claude settings

3. **Claude Integration**:
   - When enabled, active profile's API key and URL are written to Claude settings
   - When disabled, authentication is removed from Claude settings
   - Automatic sync ensures consistency between application and Claude state