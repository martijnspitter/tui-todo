package osoperations

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func GetFilePath(fileName string, isDev bool) string {
	if isDev {
		return fmt.Sprintf("./%s", fileName)
	}

	// Get the app data directory based on OS conventions
	appDir := getAppDataDir()

	// Create the directory if it doesn't exist
	os.MkdirAll(appDir, 0755)

	return filepath.Join(appDir, fileName)
}

// getAppDataDir returns the OS-specific directory for application data
func getAppDataDir() string {
	const appName = "tui-todo"

	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if we can't get home directory
		return "."
	}

	switch runtime.GOOS {
	case "windows":
		// Windows: %APPDATA%\tui-todo
		appData := os.Getenv("APPDATA")
		if appData != "" {
			return filepath.Join(appData, appName)
		}
		return filepath.Join(homeDir, "AppData", "Roaming", appName)

	case "darwin":
		// macOS: ~/Library/Application Support/tui-todo
		return filepath.Join(homeDir, "Library", "Application Support", appName)

	default:
		// Linux/Unix: ~/.local/share/tui-todo
		// First try XDG_DATA_HOME which is the standard in newer systems
		xdgDataHome := os.Getenv("XDG_DATA_HOME")
		if xdgDataHome != "" {
			return filepath.Join(xdgDataHome, appName)
		}
		// Fallback to ~/.local/share
		return filepath.Join(homeDir, ".local", "share", appName)
	}
}
