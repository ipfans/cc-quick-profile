package main

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/driver/desktop"
	"github.com/ipfans/cc-quick-profile/assets"
	"github.com/ipfans/cc-quick-profile/config"
	"github.com/ipfans/cc-quick-profile/ui"
)

var (
	configManager  *config.Manager
	fyneApp        fyne.App
	mainWindow     fyne.Window
	systemTrayMenu *fyne.Menu
)

func main() {
	// Initialize Fyne app
	fyneApp = app.New()

	// Initialize config manager
	var err error
	configManager, err = config.NewManager()
	if err != nil {
		log.Fatalf("初始化配置管理器失败: %v", err)
	}

	// Create a hidden main window (required for app lifecycle)
	mainWindow = fyneApp.NewWindow("CC Quick Profile")
	mainWindow.SetCloseIntercept(func() {
		mainWindow.Hide()
	})

	// Set up system tray if supported
	if desk, ok := fyneApp.(desktop.App); ok {
		// Set system tray icon
		desk.SetSystemTrayIcon(fyne.NewStaticResource("icon", assets.Icon))

		// Build and set system tray menu
		updateSystemTrayMenu(desk)

		log.Println("系统托盘已初始化")
	} else {
		log.Println("此平台不支持系统托盘")
	}

	// Run the app
	mainWindow.ShowAndRun()
}

func updateSystemTrayMenu(desk desktop.App) {
	settings := configManager.GetSettings()

	menuItems := []*fyne.MenuItem{}

	// Enable/Disable toggle
	enabledItem := fyne.NewMenuItem("已禁用", func() {
		settings.Enabled = !settings.Enabled
		if err := configManager.SetEnabled(settings.Enabled); err != nil {
			log.Printf("更新启用状态失败: %v", err)
		} else {
			log.Printf("启用状态已更改为: %v", settings.Enabled)
			updateSystemTrayMenu(desk)
		}
	})
	if settings.Enabled {
		enabledItem.Label = "✓ 已启用"
	}
	menuItems = append(menuItems, enabledItem)

	menuItems = append(menuItems, fyne.NewMenuItemSeparator())

	// Profile list
	if len(settings.Profiles) == 0 {
		noProfilesItem := fyne.NewMenuItem("暂无配置文件", func() {})
		noProfilesItem.Disabled = true
		menuItems = append(menuItems, noProfilesItem)
	} else {
		for _, profile := range settings.Profiles {
			p := profile // capture for closure
			menuText := p.Name
			if p.Active {
				menuText = "✓ " + menuText
			}
			profileItem := fyne.NewMenuItem(menuText, func() {
				if err := configManager.SetActiveProfile(p.Name); err != nil {
					log.Printf("设置活动配置失败: %v", err)
				} else {
					log.Printf("已切换到配置: %s", p.Name)
					// Reload settings and update menu
					if err := configManager.Load(); err != nil {
						log.Printf("重新加载配置失败: %v", err)
					}
					updateSystemTrayMenu(desk)
				}
			})
			menuItems = append(menuItems, profileItem)
		}
	}

	menuItems = append(menuItems, fyne.NewMenuItemSeparator())

	// Add new profile
	menuItems = append(menuItems, fyne.NewMenuItem("添加新配置", func() {
		showAddProfileModal(desk)
	}))

	menuItems = append(menuItems, fyne.NewMenuItemSeparator())

	// Quit
	menuItems = append(menuItems, fyne.NewMenuItem("Quit", func() {
		fyneApp.Quit()
	}))

	// Create and set the menu
	systemTrayMenu = fyne.NewMenu("CC Quick Profile", menuItems...)
	desk.SetSystemTrayMenu(systemTrayMenu)
}

func showAddProfileModal(desk desktop.App) {
	// Create event channel for modal communication
	uiEventChan := make(chan ui.Event, 1)

	// Handle events from modal
	go func() {
		for event := range uiEventChan {
			if event.Type == ui.EventConfigUpdated {
				// Reload config and update menu
				if err := configManager.Load(); err != nil {
					log.Printf("重新加载配置失败: %v", err)
				}
				updateSystemTrayMenu(desk)
			}
		}
	}()

	// Show the modal
	ui.ShowAddProfileModal(fyneApp, configManager, uiEventChan)
}
