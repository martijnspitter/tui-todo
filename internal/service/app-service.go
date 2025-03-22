package service

import (
	"fmt"
	"time"

	"github.com/charmbracelet/log"
	"github.com/martijnspitter/tui-todo/internal/models"
	"github.com/martijnspitter/tui-todo/internal/repository"
)

type AppService struct {
	todoRepo repository.TodoRepository
}

func NewAppService(todoRepo repository.TodoRepository) *AppService {
	return &AppService{
		todoRepo: todoRepo,
	}
}

func (s *AppService) CreateTodo(title, description string, priority models.Priority) (*models.Todo, error) {
	todo := &models.Todo{
		Title:       title,
		Description: description,
		Status:      models.Open,
		Priority:    priority,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err := s.todoRepo.Create(todo)
	if err != nil {
		log.Error("Failed to create todo", "error", err, "title", title)
		return nil, fmt.Errorf("couldn't create todo: %w", err)
	}

	return todo, nil
}

func (s *AppService) GetAllTodos() ([]*models.Todo, error) {
	todos, err := s.todoRepo.GetAll()
	if err != nil {
		log.Error("Failed to fetch todos", "error", err)
		return nil, fmt.Errorf("couldn't fetch todos: %w", err)
	}

	return todos, nil
}

func (s *AppService) GetTodo(id int64) (*models.Todo, error) {
	todo, err := s.todoRepo.GetByID(id)
	if err != nil {
		log.Error("Failed to fetch todo", "error", err, "id", id)
		return nil, fmt.Errorf("couldn't fetch todo #%d: %w", id, err)
	}

	return todo, nil
}

func (s *AppService) UpdateTodo(todo *models.Todo) error {
	todo.UpdatedAt = time.Now()
	err := s.todoRepo.Update(todo)
	if err != nil {
		log.Error("Failed to update todo", "error", err, "id", todo.ID)
		return fmt.Errorf("couldn't update todo #%d: %w", todo.ID, err)
	}

	return nil
}

func (s *AppService) DeleteTodo(id int64) error {
	err := s.todoRepo.Delete(id)
	if err != nil {
		log.Error("Failed to delete todo", "error", err, "id", id)
		return fmt.Errorf("couldn't delete todo #%d: %w", id, err)
	}

	return nil
}

func (s *AppService) MarkAsOpen(id int64) error {
	todo, err := s.todoRepo.GetByID(id)
	if err != nil {
		log.Error("Failed to fetch todo for status change", "error", err, "id", id)
		return fmt.Errorf("couldn't fetch todo #%d: %w", id, err)
	}

	todo.Status = models.Open
	todo.UpdatedAt = time.Now()

	err = s.todoRepo.Update(todo)
	if err != nil {
		log.Error("Failed to mark todo as open", "error", err, "id", id)
		return fmt.Errorf("couldn't mark todo #%d as open: %w", id, err)
	}

	return nil
}

func (s *AppService) MarkAsDoing(id int64) error {
	todo, err := s.todoRepo.GetByID(id)
	if err != nil {
		log.Error("Failed to fetch todo for status change", "error", err, "id", id)
		return fmt.Errorf("couldn't fetch todo #%d: %w", id, err)
	}

	todo.Status = models.Doing
	todo.UpdatedAt = time.Now()

	err = s.todoRepo.Update(todo)
	if err != nil {
		log.Error("Failed to mark todo as doing", "error", err, "id", id)
		return fmt.Errorf("couldn't mark todo #%d as doing: %w", id, err)
	}

	return nil
}

func (s *AppService) MarkAsDone(id int64) error {
	todo, err := s.todoRepo.GetByID(id)
	if err != nil {
		log.Error("Failed to fetch todo for status change", "error", err, "id", id)
		return fmt.Errorf("couldn't fetch todo #%d: %w", id, err)
	}

	todo.Status = models.Done
	todo.UpdatedAt = time.Now()

	err = s.todoRepo.Update(todo)
	if err != nil {
		log.Error("Failed to mark todo as done", "error", err, "id", id)
		return fmt.Errorf("couldn't mark todo #%d as done: %w", id, err)
	}

	return nil
}

func (s *AppService) ArchiveTodo(id int64) error {
	todo, err := s.todoRepo.GetByID(id)
	if err != nil {
		log.Error("Failed to fetch todo for archiving", "error", err, "id", id)
		return fmt.Errorf("couldn't fetch todo #%d: %w", id, err)
	}

	todo.Status = models.Archived
	todo.UpdatedAt = time.Now()

	err = s.todoRepo.Update(todo)
	if err != nil {
		log.Error("Failed to archive todo", "error", err, "id", id)
		return fmt.Errorf("couldn't archive todo #%d: %w", id, err)
	}

	return nil
}

// Filtered queries
func (s *AppService) GetOpenTodos() ([]*models.Todo, error) {
	todos, err := s.todoRepo.GetOpen()
	if err != nil {
		log.Error("Failed to fetch open todos", "error", err)
		return nil, fmt.Errorf("couldn't fetch open todos: %w", err)
	}

	return todos, nil
}

func (s *AppService) GetActiveTodos() ([]*models.Todo, error) {
	todos, err := s.todoRepo.GetActive()
	if err != nil {
		log.Error("Failed to fetch active todos", "error", err)
		return nil, fmt.Errorf("couldn't fetch active todos: %w", err)
	}

	return todos, nil
}

func (s *AppService) GetCompletedTodos() ([]*models.Todo, error) {
	todos, err := s.todoRepo.GetCompleted()
	if err != nil {
		log.Error("Failed to fetch completed todos", "error", err)
		return nil, fmt.Errorf("couldn't fetch completed todos: %w", err)
	}

	return todos, nil
}

func (s *AppService) GetArchivedTodos() ([]*models.Todo, error) {
	todos, err := s.todoRepo.GetArchived()
	if err != nil {
		log.Error("Failed to fetch archived todos", "error", err)
		return nil, fmt.Errorf("couldn't fetch archived todos: %w", err)
	}

	return todos, nil
}

func (s *AppService) SearchTodos(query string) ([]*models.Todo, error) {
	todos, err := s.todoRepo.Search(query)
	if err != nil {
		log.Error("Failed to search todos", "error", err, "query", query)
		return nil, fmt.Errorf("couldn't search for '%s': %w", query, err)
	}

	return todos, nil
}

// Tag methods
func (s *AppService) AddTagToTodo(todoID int64, tag string) error {
	err := s.todoRepo.AddTagToTodo(todoID, tag)
	if err != nil {
		log.Error("Failed to add tag to todo", "error", err, "todoID", todoID, "tag", tag)
		return fmt.Errorf("couldn't add tag '%s' to todo #%d: %w", tag, todoID, err)
	}

	return nil
}

func (s *AppService) RemoveTagFromTodo(todoID int64, tag string) error {
	err := s.todoRepo.RemoveTagFromTodo(todoID, tag)
	if err != nil {
		log.Error("Failed to remove tag from todo", "error", err, "todoID", todoID, "tag", tag)
		return fmt.Errorf("couldn't remove tag '%s' from todo #%d: %w", tag, todoID, err)
	}

	return nil
}

func (s *AppService) GetTodosByTag(tag string) ([]*models.Todo, error) {
	todos, err := s.todoRepo.FindTodosByTag(tag)
	if err != nil {
		log.Error("Failed to get todos by tag", "error", err, "tag", tag)
		return nil, fmt.Errorf("couldn't get todos with tag '%s': %w", tag, err)
	}

	return todos, nil
}

// Due date methods
func (s *AppService) SetDueDate(todoID int64, dueDate time.Time) error {
	todo, err := s.todoRepo.GetByID(todoID)
	if err != nil {
		log.Error("Failed to fetch todo for setting due date", "error", err, "id", todoID)
		return fmt.Errorf("couldn't fetch todo #%d: %w", todoID, err)
	}

	todo.DueDate = &dueDate
	todo.UpdatedAt = time.Now()

	err = s.todoRepo.Update(todo)
	if err != nil {
		log.Error("Failed to set due date", "error", err, "todoID", todoID, "dueDate", dueDate)
		return fmt.Errorf("couldn't set due date for todo #%d: %w", todoID, err)
	}

	return nil
}

func (s *AppService) ClearDueDate(todoID int64) error {
	todo, err := s.todoRepo.GetByID(todoID)
	if err != nil {
		log.Error("Failed to fetch todo for clearing due date", "error", err, "id", todoID)
		return fmt.Errorf("couldn't fetch todo #%d: %w", todoID, err)
	}

	todo.DueDate = nil
	todo.UpdatedAt = time.Now()

	err = s.todoRepo.Update(todo)
	if err != nil {
		log.Error("Failed to clear due date", "error", err, "todoID", todoID)
		return fmt.Errorf("couldn't clear due date for todo #%d: %w", todoID, err)
	}

	return nil
}

// Priority methods
func (s *AppService) SetPriority(todoID int64, priority models.Priority) error {
	todo, err := s.todoRepo.GetByID(todoID)
	if err != nil {
		log.Error("Failed to fetch todo for setting priority", "error", err, "id", todoID)
		return fmt.Errorf("couldn't fetch todo #%d: %w", todoID, err)
	}

	todo.Priority = priority
	todo.UpdatedAt = time.Now()

	err = s.todoRepo.Update(todo)
	if err != nil {
		log.Error("Failed to set priority", "error", err, "todoID", todoID, "priority", priority)
		return fmt.Errorf("couldn't set priority for todo #%d: %w", todoID, err)
	}

	return nil
}
