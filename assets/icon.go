package assets

import _ "embed"

// Icon is the embedded system tray icon data
//
//go:embed icon.png
var Icon []byte

// ClaudeDefaultSettings is the default Claude settings.json template
//
//go:embed claude_default_settings.json
var ClaudeDefaultSettings []byte
