# cc-quick-profile

A cross-platform system tray GUI application that allows you to quickly switch between Claude Code profiles. Easily manage multiple API configurations and switch between them with a single click.

## Features

- 🔄 **Quick Profile Switching**: Switch Claude Code profiles instantly from system tray
- 🌐 **Cross-Platform**: Full support for Windows, macOS, and Linux
- ⚡ **Minimal Resource Usage**: Lightweight system tray application
- 🔧 **Easy Management**: Add, activate, and manage profiles through intuitive UI
- 🔐 **Secure Storage**: Profile configurations stored locally and securely
- 🎯 **Smart Sync**: Automatic synchronization with Claude Code settings
- ✨ **Modern UI**: Clean interface built with Fyne framework

## Supported Platforms

- ✅ Windows (x64, ARM64)
- ✅ macOS (Intel, Apple Silicon)  
- ✅ Linux (x64, ARM64)

## Requirements

- Go 1.24.5 or later
- [Task](https://taskfile.dev) for development (optional)

## Installation

### From Source
```bash
go install github.com/ipfans/cc-quick-profile@latest
```

### From Releases
Download pre-built binaries from the [releases page](https://github.com/ipfans/cc-quick-profile/releases).

## Quick Start

1. **Launch the application** - The system tray icon will appear
2. **Add your first profile** - Right-click the tray icon and select "添加新配置"
3. **Fill in the details**:
   - **Profile Name**: A friendly name for your configuration
   - **API URL**: Your Claude API endpoint
   - **API Key**: Your authentication token
4. **Activate the profile** - Click on the profile name in the tray menu
5. **Start using Claude Code** with your selected profile!

## Configuration

The application manages two types of settings:

### Application Settings
- **Windows**: `%APPDATA%\cc-quick-profile\settings.json`
- **macOS/Linux**: `$HOME/.config/cc-quick-profile/settings.json`

### Claude Code Integration
- **All platforms**: `$HOME/.claude/settings.json` (managed automatically)
- Sets `ANTHROPIC_AUTH_TOKEN` and `ANTHROPIC_BASE_URL` environment variables

## Development

This project uses [Task](https://taskfile.dev) for build automation.

### Common Commands
```bash
# Show available tasks
task

# Build and run
task run

# Build for current platform
task build

# Build for all platforms
task build:all

# Run tests and checks
task check
```

### Project Structure
```
├── main.go              # Application entry point and system tray logic
├── assets/              # Embedded resources (icons, templates)
├── claude/              # Claude Code settings management
├── config/              # Application configuration management
├── models/              # Data structures (Profile, Settings)
├── ui/                  # User interface components (modal dialogs)
├── Taskfile.yml         # Build automation tasks
└── go.mod              # Go module dependencies
```

## How It Works

1. **Profile Storage**: Your profiles are stored locally in a JSON configuration file
2. **System Tray Integration**: The application runs in the background with a system tray icon
3. **Claude Integration**: When you activate a profile, the app updates Claude Code's settings automatically
4. **Cross-Platform**: Built with Fyne, ensuring consistent behavior across all supported platforms

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Run `task check` to ensure code quality
5. Commit your changes (`git commit -am 'Add amazing feature'`)
6. Push to the branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.