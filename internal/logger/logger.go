package logger

import (
	slog "log"
	"os"
	"path/filepath"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
)

func InitLogger() *os.File {
	path := getDefaultPath(true)
	// Open file with O_TRUNC flag to truncate it if it exists
	// This will create an empty file each time
	loggerFile, fileErr := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	if fileErr == nil {
		log.SetOutput(loggerFile)
		log.SetTimeFormat(time.DateTime)
		log.SetReportCaller(true)
		log.SetLevel(log.DebugLevel)
		log.Debug("Logging initialized")
	} else {
		// Fallback to Bubbletea's built-in logging
		loggerFile, _ = tea.LogToFile("debug.log", "debug")
		slog.Print("Failed setting up custom logging", fileErr)
	}

	return loggerFile
}

func getDefaultPath(isDev bool) string {
	if isDev {
		return "./debug.log"
	}
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fall back to current directory if we can't get home dir
		return "./debug.log"
	}

	// Create .klartui directory if it doesn't exist
	configDir := filepath.Join(homeDir, ".tui-todo")
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return "./debug.log" // Fall back to current directory
	}

	return filepath.Join(configDir, "debug.log")
}
