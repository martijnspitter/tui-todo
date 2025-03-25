package logger

import (
	slog "log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	osoperations "github.com/martijnspitter/tui-todo/internal/os-operations"
)

func InitLogger(version string) *os.File {
	path := osoperations.GetFilePath("debug.log", version)
	// Open file with O_TRUNC flag to truncate it if it exists
	// This will create an empty file each time
	loggerFile, fileErr := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)

	if fileErr == nil {
		log.SetOutput(loggerFile)
		log.SetTimeFormat(time.StampMilli)
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
