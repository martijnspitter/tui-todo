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
	GetBlocked() ([]*models.Todo, error)
	Search(query string) ([]*models.Todo, error)

	// tags
	AddTagToTodo(id int64, tagname string) error
	RemoveTagFromTodo(id int64, tageName string) error
	GetAllTags() ([]*models.Tag, error)
	DeleteTag(id int64) error
}

// Filter returns a WHERE clause fragment and associated arguments
type Filter func() (string, []any)

func PrioAboveHighFilter() Filter {
	return func() (string, []any) {
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		return "(priority >= ? AND (due_date IS NULL OR due_date < ?) AND status != ?)", []interface{}{
			models.Major,
			today,
			models.Done,
		}
	}
}

func OverDueFilter() Filter {
	return func() (string, []any) {
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		return "(due_date IS NOT NULL AND due_date < ? AND status != ?)", []interface{}{
			today,
			models.Done,
		}
	}
}

func DueTodayFilter() Filter {
	return func() (string, []any) {
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		tomorrow := today.Add(24 * time.Hour)

		return "(due_date IS NOT NULL AND due_date >= ? AND due_date < ? AND status != ?)",
			[]interface{}{
				today,
				tomorrow,
				models.Done, // Exclude completed tasks
			}
	}
}

func ComingUpFilter() Filter {
	return func() (string, []any) {
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		tomorrow := today.Add(24 * time.Hour)    // Start from tomorrow
		inThreeDays := today.Add(72 * time.Hour) // Show next 3 days from today

		return "(due_date IS NOT NULL AND due_date >= ? AND due_date < ? AND status != ?)",
			[]interface{}{
				tomorrow,
				inThreeDays,
				models.Done, // Exclude completed tasks
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

// AllTodayFilter combines all filters for the Today dashboard:
// - High priority tasks
// - Overdue tasks
// - Tasks due today
// - Tasks in progress
// - Coming up tasks (next 3 days)
func AllTodayFilter() Filter {
	return func() (string, []any) {
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		tomorrow := today.Add(24 * time.Hour)
		inThreeDays := today.Add(72 * time.Hour)

		whereClause := `(
            -- High priority tasks
            (priority >= ? AND status != ?)
            OR
            -- Overdue tasks
            (due_date IS NOT NULL AND due_date < ? AND status != ?)
            OR
            -- Due today
            (due_date IS NOT NULL AND due_date >= ? AND due_date < ? AND status != ?)
            OR
            -- In progress tasks
            (status = ?)
            OR
            -- Coming up tasks
            (due_date IS NOT NULL AND due_date >= ? AND due_date < ? AND status != ?)
        )`

		args := []interface{}{
			// High priority args
			models.Major,
			models.Done,

			// Overdue args
			today,
			models.Done,

			// Due today args
			today,
			tomorrow,
			models.Done,

			// In progress arg
			models.Doing,

			// Coming up args
			tomorrow,
			inThreeDays,
			models.Done,
		}

		return whereClause, args
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
