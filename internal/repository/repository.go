package repository

import (
	"time"

	"github.com/martijnspitter/tui-todo/internal/models"
)

// TodoRepository defines the contract for todo data operations
type TodoRepository interface {
	Create(todo *models.Todo) error
	GetByID(id int64) (*models.Todo, error)
	GetAll(filters ...Filter) ([]*models.Todo, error)
	Update(todo *models.Todo) error
	Delete(id int64) error

	// Additional methods specific to todos
	GetOpen() ([]*models.Todo, error)
	GetActive() ([]*models.Todo, error)
	GetCompleted() ([]*models.Todo, error)
	GetArchived() ([]*models.Todo, error)
	Search(query string) ([]*models.Todo, error)

	// tags
	AddTagToTodo(id int64, tagname string) error
	RemoveTagFromTodo(id int64, tageName string) error
	GetTodoTags(id int64) ([]string, error)
	FindTodosByTag(tagName string) ([]*models.Todo, error)
}

// Filter returns a WHERE clause fragment and associated arguments
type Filter func() (string, []any)

// Common filters you can use
func StatusFilter(status models.Status) Filter {
	return func() (string, []any) {
		return "status = ?", []any{status}
	}
}

func PriorityFilter(minPriority models.Priority) Filter {
	return func() (string, []any) {
		return "priority >= ?", []any{minPriority}
	}
}

func DueDateFilter(beforeDate time.Time) Filter {
	return func() (string, []any) {
		return "due_date IS NOT NULL AND due_date <= ?", []interface{}{beforeDate}
	}
}

func SearchFilter(query string) Filter {
	return func() (string, []any) {
		searchTerm := "%" + query + "%"
		return "(title LIKE ? OR description LIKE ?)", []interface{}{searchTerm, searchTerm}
	}
}

func TagFilter(tagName string) Filter {
	return func() (string, []any) {
		return `id IN (
            SELECT tt.todo_id
            FROM todo_tags tt
            JOIN tags t ON tt.tag_id = t.id
            WHERE t.name = ?
        )`, []any{tagName}
	}
}
