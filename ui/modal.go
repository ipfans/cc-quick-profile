// Package ui provides user interface components for the application
package ui

import (
	"fmt"
	"net/url"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/ipfans/cc-quick-profile/config"
	"github.com/ipfans/cc-quick-profile/models"
)

// EventType represents UI events
type EventType string

const (
	// EventConfigUpdated is sent when configuration has been updated
	EventConfigUpdated EventType = "config_updated"
)

// Event represents a UI event
type Event struct {
	Type EventType
	Data interface{}
}

// ShowAddProfileModal displays the add profile dialog
func ShowAddProfileModal(app fyne.App, configManager *config.Manager, eventChan chan<- Event) {
	window := app.NewWindow("添加新配置")
	window.Resize(fyne.NewSize(400, 300))
	window.CenterOnScreen()

	// Create form fields
	nameEntry := widget.NewEntry()
	nameEntry.SetPlaceHolder("配置名称")

	apiURLEntry := widget.NewEntry()
	apiURLEntry.SetPlaceHolder("https://api.example.com")

	apiKeyEntry := widget.NewPasswordEntry()
	apiKeyEntry.SetPlaceHolder("API 密钥")

	// Validation error label
	errorLabel := widget.NewLabel("")
	errorLabel.Hide()
	errorLabel.TextStyle = fyne.TextStyle{Bold: true}

	// Validate input
	validate := func() error {
		// Clear previous error
		errorLabel.Hide()
		errorLabel.SetText("")

		// Check name
		name := strings.TrimSpace(nameEntry.Text)
		if name == "" {
			return fmt.Errorf("配置名称不能为空")
		}

		// Check API URL
		apiURL := strings.TrimSpace(apiURLEntry.Text)
		if apiURL == "" {
			return fmt.Errorf("API URL 不能为空")
		}

		// Validate URL format
		if _, err := url.Parse(apiURL); err != nil {
			return fmt.Errorf("API URL 格式无效")
		}

		// Check API key
		apiKey := strings.TrimSpace(apiKeyEntry.Text)
		if apiKey == "" {
			return fmt.Errorf("API 密钥不能为空")
		}

		return nil
	}

	// Save button handler
	onSave := func() {
		if err := validate(); err != nil {
			errorLabel.SetText(err.Error())
			errorLabel.Show()
			return
		}

		// Create new profile
		profile := models.Profile{
			Name:   strings.TrimSpace(nameEntry.Text),
			APIURL: strings.TrimSpace(apiURLEntry.Text),
			APIKey: strings.TrimSpace(apiKeyEntry.Text),
			Active: false,
		}

		// Add to config
		if err := configManager.AddProfile(profile); err != nil {
			errorLabel.SetText(fmt.Sprintf("保存失败: %v", err))
			errorLabel.Show()
			return
		}

		// Send update event
		eventChan <- Event{Type: EventConfigUpdated}

		// Show success and close
		dialog.ShowInformation("成功",
			fmt.Sprintf("配置 '%s' 已成功添加。", profile.Name),
			window)

		window.Close()
	}

	// Cancel button handler
	onCancel := func() {
		window.Close()
	}

	// Create form
	form := container.NewVBox(
		widget.NewLabel("添加新的 Claude Code 配置"),
		widget.NewSeparator(),
		container.NewGridWithColumns(1,
			widget.NewLabel("配置名称:"),
			nameEntry,
			widget.NewLabel("API URL:"),
			apiURLEntry,
			widget.NewLabel("API 密钥:"),
			apiKeyEntry,
		),
		errorLabel,
		widget.NewSeparator(),
		container.NewGridWithColumns(2,
			widget.NewButton("取消", onCancel),
			widget.NewButton("保存", onSave),
		),
	)

	// Set content and show
	window.SetContent(container.NewPadded(form))
	window.Show()

	// Set initial focus
	window.Canvas().Focus(nameEntry)
}
