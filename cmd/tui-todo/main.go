package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	_ "github.com/mattn/go-sqlite3"

	"github.com/martijnspitter/tui-todo/internal/logger"
	"github.com/martijnspitter/tui-todo/internal/repository"
	"github.com/martijnspitter/tui-todo/internal/service"
	"github.com/martijnspitter/tui-todo/internal/ui"
)

func main() {
	logger := logger.InitLogger()
	if logger != nil {
		defer logger.Close()
	}

	todoRepo, err := repository.NewSQLiteTodoRepository()
	if err != nil {
		log.Error("Failed to start db", err)
		os.Exit(1)
	}
	defer todoRepo.Close()

	service := service.NewAppService(todoRepo)

	baseModel := ui.NewBaseModel(service)

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
