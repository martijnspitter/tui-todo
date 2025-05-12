package service

import (
	"fmt"
	"sort"
	"time"

	"github.com/charmbracelet/log"
	"github.com/martijnspitter/tui-todo/internal/models"
	"github.com/martijnspitter/tui-todo/internal/repository"
	"github.com/martijnspitter/tui-todo/internal/utils"
)

type UpdateInfo struct {
	Version     string
	URL         string
	Notes       string
	ForceUpdate bool
	HasUpdate   bool
	CheckedAt   time.Time
}

type AppService struct {
	todoRepo   repository.TodoRepository
	updateInfo *UpdateInfo
}

func NewAppService(todoRepo repository.TodoRepository) *AppService {
	return &AppService{
		todoRepo:   todoRepo,
		updateInfo: &UpdateInfo{},
	}
}

// ===========================================================================
// Todo Methods
// ===========================================================================
func (s *AppService) SaveTodo(todo *models.Todo, tags []string) error {
	// Service decides whether to create or update based on ID or other criteria
	if todo.ID == 0 {
		// Create new
		return s.CreateTodo(todo.Title, todo.Description, todo.Priority, tags, todo.DueDate, todo.Status)
	} else {
		// Update existing
		return s.UpdateTodo(todo, tags)
	}
}

func (s *AppService) CreateTodo(title, description string, priority models.Priority, tags []string, dueDate *time.Time, status models.Status) error {
	todo := &models.Todo{
		Title:       title,
		Description: description,
		Status:      status,
		Priority:    priority,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		DueDate:     dueDate,
		TimeSpent:   0,
	}

	// If creating a task directly in Doing status, set time_started
	if status == models.Doing {
		now := time.Now()
		todo.TimeStarted = &now
	}

	err := s.todoRepo.Create(todo)
	if err != nil {
		log.Error("Failed to create todo", "error", err, "title", title)
		return fmt.Errorf("error.create_failed")
	}

	for _, tag := range tags {
		err := s.AddTagToTodo(todo.ID, tag)
		if err != nil {
			log.Error("Could not add tag: %s %w", tag, err)
			return fmt.Errorf("error.tag_add_failed")
		}
	}

	return nil
}

func (s *AppService) GetAllTodos(showArchived bool) ([]*models.Todo, error) {
	archivedFilter := repository.NotArchivedFilter()
	if showArchived {
		archivedFilter = repository.ArchivedFilter()
	}
	todos, err := s.todoRepo.GetAll(archivedFilter)
	if err != nil {
		log.Error("Failed to fetch todos", "error", err)
		return nil, fmt.Errorf("error.todos_not_found")
	}

	return sortTodos(todos), nil
}

func (s *AppService) GetTodo(id int64) (*models.Todo, error) {
	todo, err := s.todoRepo.GetByID(id)
	if err != nil {
		log.Error("Failed to fetch todo", "error", err, "id", id)
		return nil, fmt.Errorf("error.todo_not_found")
	}

	return todo, nil
}

func (s *AppService) UpdateTodo(todo *models.Todo, tags []string) error {
	todo.UpdatedAt = time.Now()
	err := s.todoRepo.Update(todo)
	if err != nil {
		log.Error("Failed to update todo", "error", err, "id", todo.ID)
		return fmt.Errorf("error.update_failed")
	}

	for _, tag := range tags {
		err := s.AddTagToTodo(todo.ID, tag)
		if err != nil {
			log.Error("Could not add tag: %s %w", tag, err)
			return fmt.Errorf("error.tag_add_failed")
		}
	}

	return nil
}

func (s *AppService) DeleteTodo(id int64) error {
	err := s.todoRepo.Delete(id)
	if err != nil {
		log.Error("Failed to delete todo", "error", err, "id", id)
		return fmt.Errorf("error.delete_failed")
	}

	return nil
}

func (s *AppService) MarkAsOpen(id int64) error {
	todo, err := s.todoRepo.GetByID(id)
	if err != nil {
		log.Error("Failed to fetch todo for status change", "error", err, "id", id)
		return fmt.Errorf("error.todo_not_found")
	}

	// If transitioning from Doing to Open, calculate elapsed time
	if todo.Status == models.Doing && todo.TimeStarted != nil {
		elapsed := time.Since(*todo.TimeStarted).Seconds()
		todo.TimeSpent += int64(elapsed)
		todo.TimeStarted = nil // Clear the start time when moving to Open
	}

	todo.Status = models.Open
	todo.UpdatedAt = time.Now()

	err = s.todoRepo.Update(todo)
	if err != nil {
		log.Error("Failed to mark todo as open", "error", err, "id", id)
		return fmt.Errorf("error.status_change_failed")
	}

	return nil
}

func (s *AppService) MarkAsDoing(id int64) error {
	todo, err := s.todoRepo.GetByID(id)
	if err != nil {
		log.Error("Failed to fetch todo for status change", "error", err, "id", id)
		return fmt.Errorf("error.todo_not_found")
	}

	// Set time_started only if not already in Doing status
	if todo.Status != models.Doing {
		now := time.Now()
		todo.TimeStarted = &now
	}

	todo.Status = models.Doing
	todo.UpdatedAt = time.Now()

	err = s.todoRepo.Update(todo)
	if err != nil {
		log.Error("Failed to mark todo as doing", "error", err, "id", id)
		return fmt.Errorf("error.status_change_failed")
	}

	return nil
}

func (s *AppService) MarkAsDone(id int64) error {
	todo, err := s.todoRepo.GetByID(id)
	if err != nil {
		log.Error("Failed to fetch todo for status change", "error", err, "id", id)
		return fmt.Errorf("error.todo_not_found")
	}

	// Calculate and accumulate time spent if task was in Doing status
	if todo.Status == models.Doing && todo.TimeStarted != nil {
		elapsed := time.Since(*todo.TimeStarted).Seconds()
		todo.TimeSpent += int64(elapsed)
		todo.TimeStarted = nil // Clear the start time
	}

	todo.Status = models.Done
	todo.UpdatedAt = time.Now()

	err = s.todoRepo.Update(todo)
	if err != nil {
		log.Error("Failed to mark todo as done", "error", err, "id", id)
		return fmt.Errorf("error.status_change_failed")
	}

	return nil
}

func (s *AppService) ArchiveTodo(todoID int64) error {
	todo, err := s.todoRepo.GetByID(todoID)
	if err != nil {
		log.Error("Failed to fetch todo for archiving", "error", err, "id", todoID)
		return fmt.Errorf("error.todo_not_found")
	}

	todo.Archived = true
	todo.UpdatedAt = time.Now()

	err = s.todoRepo.Update(todo)
	if err != nil {
		log.Error("Failed to archive todo", "error", err, "id", todoID)
		return fmt.Errorf("error.archive_failed")
	}

	return nil
}

func (s *AppService) UnarchiveTodo(todoID int64) error {
	todo, err := s.todoRepo.GetByID(todoID)
	if err != nil {
		log.Error("Failed to fetch todo for unarchiving", "error", err, "id", todoID)
		return fmt.Errorf("error.todo_not_found")
	}

	todo.Archived = false
	todo.UpdatedAt = time.Now()

	err = s.todoRepo.Update(todo)
	if err != nil {
		log.Error("Failed to unarchive todo", "error", err, "id", todoID)
		return fmt.Errorf("error.unarchive_failed")
	}

	return nil
}

// Filtered queries
func (s *AppService) GetOpenTodos() ([]*models.Todo, error) {
	todos, err := s.todoRepo.GetOpen()
	if err != nil {
		log.Error("Failed to fetch open todos", "error", err)
		return nil, fmt.Errorf("error.todos_not_found")
	}

	return sortTodos(todos), nil
}

func (s *AppService) GetActiveTodos() ([]*models.Todo, error) {
	todos, err := s.todoRepo.GetActive()
	if err != nil {
		log.Error("Failed to fetch active todos", "error", err)
		return nil, fmt.Errorf("error.todos_not_found")
	}

	return sortTodos(todos), nil
}

func (s *AppService) GetCompletedTodos() ([]*models.Todo, error) {
	todos, err := s.todoRepo.GetCompleted()
	if err != nil {
		log.Error("Failed to fetch completed todos", "error", err)
		return nil, fmt.Errorf("error.todos_not_found")
	}

	return sortTodos(todos), nil
}

func (s *AppService) GetTodosForToday() (highPrio, dueToday, inProgress, overDue, comingUp []*models.Todo, error error) {
	// Get todos that are due today and not archived
	highPrio, err := s.todoRepo.GetAll(repository.PrioAboveHighFilter(), repository.NotArchivedFilter())
	if err != nil {
		log.Error("Failed to fetch highPrio for today", "error", err)
		return nil, nil, nil, nil, nil, fmt.Errorf("error.todos_not_found")
	}
	dueToday, err = s.todoRepo.GetAll(repository.DueTodayFilter(), repository.NotArchivedFilter())
	if err != nil {
		log.Error("Failed to fetch dueToday for today", "error", err)
		return nil, nil, nil, nil, nil, fmt.Errorf("error.todos_not_found")
	}
	inProgress, err = s.todoRepo.GetActive()
	if err != nil {
		log.Error("Failed to fetch inProgress for today", "error", err)
		return nil, nil, nil, nil, nil, fmt.Errorf("error.todos_not_found")
	}
	overDue, err = s.todoRepo.GetAll(repository.OverDueFilter(), repository.NotArchivedFilter())
	if err != nil {
		log.Error("Failed to fetch overDue for today", "error", err)
		return nil, nil, nil, nil, nil, fmt.Errorf("error.todos_not_found")
	}
	comingUp, err = s.todoRepo.GetAll(repository.ComingUpFilter(), repository.NotArchivedFilter())
	if err != nil {
		log.Error("Failed to fetch comingUp for today", "error", err)
		return nil, nil, nil, nil, nil, fmt.Errorf("error.todos_not_found")
	}

	return sortTodos(highPrio), sortTodos(dueToday), sortTodos(inProgress), sortTodos(overDue), sortTodos(comingUp), nil
}

func (s *AppService) GetTodayCompletionStats() (completed int, total int, formattedTimeSpent string) {
	var timeSpent int64
	// Get tasks completed today
	completedToday, err := s.todoRepo.GetAll(repository.CompletedTodayFilter())
	if err != nil {
		log.Error("Failed to fetch today's completed tasks", "error", err)
		return 0, 0, ""
	}

	// Get all tasks that would show in today's dashboard (not completed yet)
	currentTodayTasks, err := s.todoRepo.GetAll(repository.AllTodayFilter(), repository.NotArchivedFilter())
	if err != nil {
		log.Error("Failed to fetch today's tasks", "error", err)
		return len(completedToday), len(completedToday), ""
	}

	for _, todo := range completedToday {
		timeSpent = todo.GetTotalSeconds()
	}
	for _, todo := range currentTodayTasks {
		timeSpent = todo.GetTotalSeconds()
	}

	return len(completedToday), len(completedToday) + len(currentTodayTasks), utils.FormatTime(timeSpent)
}

// Tag methods
func (s *AppService) AddTagToTodo(todoID int64, tag string) error {
	err := s.todoRepo.AddTagToTodo(todoID, tag)
	if err != nil {
		log.Error("Failed to add tag to todo", "error", err, "todoID", todoID, "tag", tag)
		return fmt.Errorf("error.tag_add_failed")
	}

	return nil
}

func (s *AppService) RemoveTagFromTodo(todoID int64, tag string) error {
	err := s.todoRepo.RemoveTagFromTodo(todoID, tag)
	if err != nil {
		log.Error("Failed to remove tag from todo", "error", err, "todoID", todoID, "tag", tag)
		return fmt.Errorf("error.tag_remove_failed")
	}

	return nil
}

// Due date methods
func (s *AppService) SetDueDate(todoID int64, dueDate time.Time) error {
	todo, err := s.todoRepo.GetByID(todoID)
	if err != nil {
		log.Error("Failed to fetch todo for setting due date", "error", err, "id", todoID)
		return fmt.Errorf("error.todos_not_found")
	}

	todo.DueDate = &dueDate
	todo.UpdatedAt = time.Now()

	err = s.todoRepo.Update(todo)
	if err != nil {
		log.Error("Failed to set due date", "error", err, "todoID", todoID, "dueDate", dueDate)
		return fmt.Errorf("error.update_failed")
	}

	return nil
}

func (s *AppService) ClearDueDate(todoID int64) error {
	todo, err := s.todoRepo.GetByID(todoID)
	if err != nil {
		log.Error("Failed to fetch todo for clearing due date", "error", err, "id", todoID)
		return fmt.Errorf("error.todos_not_found")
	}

	todo.DueDate = nil
	todo.UpdatedAt = time.Now()

	err = s.todoRepo.Update(todo)
	if err != nil {
		log.Error("Failed to clear due date", "error", err, "todoID", todoID)
		return fmt.Errorf("error.update_failed")
	}

	return nil
}

// Priority methods
func (s *AppService) SetPriority(todoID int64, priority models.Priority) error {
	todo, err := s.todoRepo.GetByID(todoID)
	if err != nil {
		log.Error("Failed to fetch todo for setting priority", "error", err, "id", todoID)
		return fmt.Errorf("error.todos_not_found")
	}

	todo.Priority = priority
	todo.UpdatedAt = time.Now()

	err = s.todoRepo.Update(todo)
	if err != nil {
		log.Error("Failed to set priority", "error", err, "todoID", todoID, "priority", priority)
		return fmt.Errorf("error.update_failed")
	}

	return nil
}

// Time tracking methods
func (s *AppService) GetTotalTimeSpent(todoID int64) (time.Duration, error) {
	todo, err := s.todoRepo.GetByID(todoID)
	if err != nil {
		log.Error("Failed to fetch todo for time tracking", "error", err, "id", todoID)
		return 0, fmt.Errorf("error.todos_not_found")
	}

	// Get the base time spent
	totalSeconds := todo.TimeSpent

	// If currently in "Doing" status, add the current session time
	if todo.Status == models.Doing && todo.TimeStarted != nil {
		currentSessionSeconds := int64(time.Since(*todo.TimeStarted).Seconds())
		totalSeconds += currentSessionSeconds
	}

	return time.Duration(totalSeconds) * time.Second, nil
}

// PauseTimeTracking pauses time tracking without changing the status
func (s *AppService) PauseTimeTracking(todoID int64) error {
	todo, err := s.todoRepo.GetByID(todoID)
	if err != nil {
		log.Error("Failed to fetch todo for pausing time tracking", "error", err, "id", todoID)
		return fmt.Errorf("error.todos_not_found")
	}

	// Only process if the task is in Doing status and has a start time
	if todo.Status == models.Doing && todo.TimeStarted != nil {
		elapsed := time.Since(*todo.TimeStarted).Seconds()
		todo.TimeSpent += int64(elapsed)
		todo.TimeStarted = nil // Clear the start time but keep status

		err = s.todoRepo.Update(todo)
		if err != nil {
			log.Error("Failed to pause time tracking", "error", err, "todoID", todoID)
			return fmt.Errorf("error.update_failed")
		}
	}

	return nil
}

// ResumeTimeTracking resumes time tracking without changing the status
func (s *AppService) ResumeTimeTracking(todoID int64) error {
	todo, err := s.todoRepo.GetByID(todoID)
	if err != nil {
		log.Error("Failed to fetch todo for resuming time tracking", "error", err, "id", todoID)
		return fmt.Errorf("error.todos_not_found")
	}

	// Only resume if the task is in Doing status but doesn't have a start time
	if todo.Status == models.Doing && todo.TimeStarted == nil {
		now := time.Now()
		todo.TimeStarted = &now

		err = s.todoRepo.Update(todo)
		if err != nil {
			log.Error("Failed to resume time tracking", "error", err, "todoID", todoID)
			return fmt.Errorf("error.update_failed")
		}
	}

	return nil
}

// ResetTimeTracking resets all time tracking data for a todo
func (s *AppService) ResetTimeTracking(todoID int64) error {
	todo, err := s.todoRepo.GetByID(todoID)
	if err != nil {
		log.Error("Failed to fetch todo for resetting time tracking", "error", err, "id", todoID)
		return fmt.Errorf("error.todos_not_found")
	}

	todo.TimeSpent = 0
	todo.TimeStarted = nil

	err = s.todoRepo.Update(todo)
	if err != nil {
		log.Error("Failed to reset time tracking", "error", err, "todoID", todoID)
		return fmt.Errorf("error.update_failed")
	}

	return nil
}

func sortTodos(todos []*models.Todo) []*models.Todo {
	// Sort todos by priority (high to low) and then by updatedAt (newest first)
	sort.Slice(todos, func(i, j int) bool {
		// If priorities are different, sort by priority (high to low)
		if todos[i].Priority != todos[j].Priority {
			return todos[i].Priority > todos[j].Priority
		}

		// If priorities are the same, sort by updatedAt (newest first)
		return todos[i].UpdatedAt.After(todos[j].UpdatedAt)
	})

	return todos
}

func (s *AppService) AdvanceStatus(todoID int64) (models.Status, error) {
	todo, err := s.todoRepo.GetByID(todoID)
	if err != nil {
		log.Error("Failed to fetch todo for status change", "error", err, "id", todoID)
		return 0, fmt.Errorf("error.todos_not_found")
	}

	if todo.Status == models.Done {
		return 0, fmt.Errorf("error.update_from_done")
	}

	var newStatus models.Status
	switch todo.Status {
	case models.Open:
		newStatus = models.Doing
		err = s.MarkAsDoing(todoID)
	case models.Doing:
		newStatus = models.Done
		err = s.MarkAsDone(todoID)
	}

	if err != nil {
		log.Error("Failed to advance status", "error", err, "todoID", todoID, "fromStatus", todo.Status, "toStatus", newStatus)
		return 0, fmt.Errorf("error.update_failed")
	}

	return newStatus, nil
}

func (s *AppService) GetFilteredTodos(currentView ViewType, showArchived bool) ([]*models.Todo, error) {
	var todos []*models.Todo
	var err error

	switch currentView {
	case OpenPane:
		todos, err = s.GetOpenTodos()
	case DoingPane:
		todos, err = s.GetActiveTodos()
	case DonePane:
		todos, err = s.GetCompletedTodos()
	case AllPane:
		todos, err = s.GetAllTodos(showArchived)
	default:
		err = fmt.Errorf("error.unknown")
	}

	if err != nil {
		return nil, err
	}

	return todos, nil
}

// ===========================================================================
// Update Info Methods
// ===========================================================================
func (s *AppService) SetUpdateInfo(version, releaseUrl, releaseNotes string, forceUpdate, hasUpdate bool) {
	s.updateInfo = &UpdateInfo{
		Version:     version,
		URL:         releaseUrl,
		Notes:       releaseNotes,
		ForceUpdate: forceUpdate,
		HasUpdate:   hasUpdate,
		CheckedAt:   time.Now(),
	}
}

func (s *AppService) GetUpdateInfo() *UpdateInfo {
	return s.updateInfo
}

func (s *AppService) HasUpdate() bool {
	return s.updateInfo != nil && s.updateInfo.HasUpdate
}

func (s *AppService) NeedsForceUpdate() bool {
	return s.updateInfo != nil && s.updateInfo.ForceUpdate
}
