package main

import (
	"context"
	"fmt"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/log"
	_ "modernc.org/sqlite"

	"github.com/martijnspitter/tui-todo/internal/i18n"
	"github.com/martijnspitter/tui-todo/internal/logger"
	"github.com/martijnspitter/tui-todo/internal/repository"
	"github.com/martijnspitter/tui-todo/internal/service"
	"github.com/martijnspitter/tui-todo/internal/ui"
	"github.com/martijnspitter/tui-todo/internal/version"
)

func main() {
	appVersion := version.GetVersion()
	if len(os.Args) > 1 && (os.Args[1] == "--version" || os.Args[1] == "-v") {
		fmt.Printf("todo version %s\n", appVersion)
		os.Exit(0)
	}

	logger := logger.InitLogger(appVersion)
	if logger != nil {
		defer logger.Close()
	}

	todoRepo, err := repository.NewSQLiteTodoRepository(appVersion)
	if err != nil {
		log.Error("Failed to start db", err)
		os.Exit(1)
	}
	defer todoRepo.Close()

	translationService, err := i18n.NewTranslationService("en")
	if err != nil {
		log.Fatal(err)
	}
	service := service.NewAppService(todoRepo)

	baseModel := ui.NewBaseModel(service, translationService)

	// Initialize TUI with endpoints as options
	p := tea.NewProgram(
		baseModel,
		tea.WithAltScreen(),       // use the full size of the terminal in its "alternate screen buffer"
		tea.WithMouseCellMotion(), // turn on mouse support so we can track the mouse wheel
	)

	go func() {
		// Wait a short period to let the UI initialize
		time.Sleep(1 * time.Second)

		ctx := context.Background()
		updateInfo, err := version.CheckForUpdates(ctx, appVersion)
		if err != nil {
			log.Error("Failed to check for updates", "error", err)
			return
		}

		if updateInfo != nil {
			// Store the update info in the service
			service.SetUpdateInfo(
				updateInfo.Version,
				updateInfo.ReleaseURL,
				updateInfo.ReleaseNotes,
				updateInfo.ForceUpdate,
				updateInfo.HasUpdate,
			)
			// Send a message to the program to notify about the update
			p.Send(ui.UpdateCheckCompletedMsg{
				ForceUpdate: updateInfo.ForceUpdate,
			})
		}
	}()

	if _, err := p.Run(); err != nil {
		fmt.Printf("Alas, there's been an error: %v", err)
		os.Exit(1)
	}
}
