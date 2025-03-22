package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/martijnspitter/tui-todo/internal/logger"
	"github.com/martijnspitter/tui-todo/internal/models"
)

func main() {
	logger := logger.InitLogger()
	if logger != nil {
		defer logger.Close()
	}

	baseModel := models.NewBaseModel()

	// Initialize TUI with endpoints as options
	p := tea.NewProgram(
		baseModel,
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
