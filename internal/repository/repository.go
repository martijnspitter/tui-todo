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
	Search(query string) ([]*models.Todo, error)

	// tags
	AddTagToTodo(id int64, tagname string) error
	RemoveTagFromTodo(id int64, tageName string) error
}

// Filter returns a WHERE clause fragment and associated arguments
type Filter func() (string, []any)

// Common filters you can use
func DueTodayOrPrioAboveHighFilter() Filter {
	return func() (string, []any) {
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		tomorrow := today.Add(24 * time.Hour)

		// Only include:
		// 1. Tasks due today (not completed), OR
		// 2. Tasks that are overdue (not completed), OR
		// 3. Tasks with Major or Critical priority (not completed)
		return "((" +
				// Due today
				"(due_date IS NOT NULL AND due_date >= ? AND due_date < ?)" +
				" OR " +
				// Overdue
				"(due_date IS NOT NULL AND due_date < ?)" +
				" OR " +
				// High priority (without future due date)
				"(priority >= ? AND (due_date IS NULL OR due_date < ?))" +
				") AND status != ?)",
			[]interface{}{
				today, tomorrow, // For due today
				today,        // For overdue
				models.Major, // For high priority
				tomorrow,     // Exclude future high priority tasks
				models.Done,  // Exclude completed tasks
			}
	}
}

func CompletedTodayFilter() Filter {
	return func() (string, []any) {
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())

		return "status = ? AND updated_at >= ?",
			[]interface{}{models.Done, today}
	}
}

func StatusFilter(status models.Status) Filter {
	return func() (string, []any) {
		return "status = ?", []any{status}
	}
}

func ArchivedFilter() Filter {
	return func() (string, []any) {
		return "archived = 1", []any{}
	}
}

func NotArchivedFilter() Filter {
	return func() (string, []any) {
		return "archived = 0", []any{}
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
