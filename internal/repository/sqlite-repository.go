package repository

import (
	"database/sql"
	"fmt"
	"sort"
	"strings"
	"time"

	"slices"

	"github.com/martijnspitter/tui-todo/internal/models"
	osoperations "github.com/martijnspitter/tui-todo/internal/os-operations"
	_ "modernc.org/sqlite"
)

type SQLiteTodoRepository struct {
	db *sql.DB
}

func NewSQLiteTodoRepository(version string) (*SQLiteTodoRepository, error) {
	path := osoperations.GetFilePath("todo.sql", version)

	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(time.Hour)

	// Initialize database schema if needed
	if err := initSchema(db); err != nil {
		return nil, err
	}

	migrationManager := NewMigrationManager(db)
	if err := migrationManager.Initialize(); err != nil {
		return nil, fmt.Errorf("couldn't initialize migration manager: %w", err)
	}

	if err := migrationManager.ApplyMigrations(GetAllMigrations()); err != nil {
		return nil, fmt.Errorf("couldn't apply migrations: %w", err)
	}

	return &SQLiteTodoRepository{db: db}, nil
}

func (r *SQLiteTodoRepository) Close() error {
	return r.db.Close()
}

func (r *SQLiteTodoRepository) Create(todo *models.Todo) error {
	// Implementation with SQL
	stmt, err := r.db.Prepare(`
        INSERT INTO todos (title, description, status, created_at, updated_at, priority, due_date, archived)
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	result, err := stmt.Exec(
		todo.Title,
		todo.Description,
		todo.Status,
		todo.CreatedAt,
		todo.UpdatedAt,
		todo.Priority,
		todo.DueDate,
		todo.Archived,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	todo.ID = id
	return nil
}

func (r *SQLiteTodoRepository) GetByID(id int64) (*models.Todo, error) {
	// Query to get todo with its tags in a single operation
	rows, err := r.db.Query(`
        SELECT t.id, t.title, t.description, t.status, t.created_at, t.updated_at,
               t.due_date, t.priority, t.archived, tag.name as tag_name
        FROM todos t
        LEFT JOIN todo_tags tt ON t.id = tt.todo_id
        LEFT JOIN tags tag ON tt.tag_id = tag.id
        WHERE t.id = ?
    `, id)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Variables for todo data
	var todo *models.Todo
	var foundTodo bool = false

	// Process result rows
	for rows.Next() {
		var todoID int64
		var title, description string
		var status models.Status
		var createdAt, updatedAt time.Time
		var dueDate sql.NullTime
		var priority models.Priority
		var tagName sql.NullString
		var archived bool

		// Scan row data
		if err := rows.Scan(
			&todoID,
			&title,
			&description,
			&status,
			&createdAt,
			&updatedAt,
			&dueDate,
			&priority,
			&archived,
			&tagName,
		); err != nil {
			return nil, err
		}

		// If this is our first row, initialize the todo
		if !foundTodo {
			todo = &models.Todo{
				ID:          todoID,
				Title:       title,
				Description: description,
				Status:      status,
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
				Priority:    priority,
				Tags:        []string{},
				Archived:    archived,
			}

			if dueDate.Valid {
				todo.DueDate = &dueDate.Time
			}

			foundTodo = true
		}

		// Add tag if present and not already in the list
		if tagName.Valid && tagName.String != "" {
			// Check if tag already exists in the slice
			tagExists := slices.Contains(todo.Tags, tagName.String)

			if !tagExists {
				todo.Tags = append(todo.Tags, tagName.String)
			}
		}
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// If no todo was found
	if !foundTodo {
		return nil, fmt.Errorf("todo with id %d not found", id)
	}

	return todo, nil
}

func (r *SQLiteTodoRepository) GetAll(filters ...Filter) ([]*models.Todo, error) {
	// Base query with joins to fetch todos and their tags
	query := `
     SELECT t.id, t.title, t.description, t.status, t.created_at, t.updated_at,
            t.due_date, t.priority, t.archived, tag.name as tag_name
     FROM todos t
     LEFT JOIN todo_tags tt ON t.id = tt.todo_id
     LEFT JOIN tags tag ON tt.tag_id = tag.id
 `

	// Apply any filters
	args := []any{}
	if len(filters) > 0 {
		query += " WHERE "
		filterClauses := []string{}

		for _, filter := range filters {
			clause, filterArgs := filter()
			if clause != "" {
				filterClauses = append(filterClauses, clause)
				args = append(args, filterArgs...)
			}
		}

		query += strings.Join(filterClauses, " AND ")
	}

	// Add ordering
	query += " ORDER BY t.created_at DESC"

	// Execute the query
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Map to store todos by ID to avoid duplicates
	todosMap := make(map[int64]*models.Todo)

	// Iterate through the result set
	for rows.Next() {
		var todoID int64
		var title, description string
		var status models.Status
		var createdAt, updatedAt time.Time
		var dueDate sql.NullTime
		var priority models.Priority
		var tagName sql.NullString
		var archived bool

		// Scan the row
		if err := rows.Scan(
			&todoID,
			&title,
			&description,
			&status,
			&createdAt,
			&updatedAt,
			&dueDate,
			&priority,
			&archived,
			&tagName,
		); err != nil {
			return nil, err
		}

		// Get or create todo in the map
		todo, exists := todosMap[todoID]
		if !exists {
			todo = &models.Todo{
				ID:          todoID,
				Title:       title,
				Description: description,
				Status:      status,
				CreatedAt:   createdAt,
				UpdatedAt:   updatedAt,
				Priority:    priority,
				Archived:    archived,
				Tags:        []string{},
			}

			if dueDate.Valid {
				todo.DueDate = &dueDate.Time
			}

			todosMap[todoID] = todo
		}

		// Add tag if present
		if tagName.Valid && tagName.String != "" {
			// Check if tag is already in the slice to avoid duplicates
			tagExists := slices.Contains(todo.Tags, tagName.String)

			if !tagExists {
				todo.Tags = append(todo.Tags, tagName.String)
			}
		}
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, err
	}

	// Convert map to slice
	todos := make([]*models.Todo, 0, len(todosMap))
	for _, todo := range todosMap {
		todos = append(todos, todo)
	}

	// Sort by created_at to maintain order
	sort.Slice(todos, func(i, j int) bool {
		return todos[i].CreatedAt.After(todos[j].CreatedAt)
	})

	return todos, nil
}

func (r *SQLiteTodoRepository) Update(todo *models.Todo) error {
	stmt, err := r.db.Prepare(`
        UPDATE todos
        SET title = ?, description = ?, status = ?, updated_at = ?, due_date = ?, priority = ?, archived = ?
        WHERE id = ?
    `)
	if err != nil {
		return err
	}
	defer stmt.Close()

	var dueDate any
	if todo.DueDate != nil {
		dueDate = *todo.DueDate
	} else {
		dueDate = nil
	}

	_, err = stmt.Exec(
		todo.Title,
		todo.Description,
		todo.Status,
		time.Now(), // Update the updated_at time
		dueDate,
		todo.Priority,
		todo.Archived,
		todo.ID,
	)
	return err
}

func (r *SQLiteTodoRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM todos WHERE id = ?", id)
	return err
}

func (r *SQLiteTodoRepository) GetOpen() ([]*models.Todo, error) {
	return r.GetAll(StatusFilter(models.Open), NotArchivedFilter())
}

func (r *SQLiteTodoRepository) GetActive() ([]*models.Todo, error) {
	return r.GetAll(StatusFilter(models.Doing), NotArchivedFilter())
}

// Get completed todos
func (r *SQLiteTodoRepository) GetCompleted() ([]*models.Todo, error) {
	return r.GetAll(StatusFilter(models.Done), NotArchivedFilter())
}

// Get archived todos
func (r *SQLiteTodoRepository) GetArchived() ([]*models.Todo, error) {
	return r.GetAll(ArchivedFilter())
}

// Search todos
func (r *SQLiteTodoRepository) Search(query string) ([]*models.Todo, error) {
	return r.GetAll(SearchFilter(query))
}

func (r *SQLiteTodoRepository) AddTagToTodo(todoID int64, tagName string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Get or create tag
	var tagID int64
	err = tx.QueryRow("SELECT id FROM tags WHERE name = ?", tagName).Scan(&tagID)
	if err == sql.ErrNoRows {
		// Tag doesn't exist, create it
		result, err := tx.Exec("INSERT INTO tags (name) VALUES (?)", tagName)
		if err != nil {
			return err
		}
		tagID, err = result.LastInsertId()
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	// Add relationship (ignore if already exists)
	_, err = tx.Exec(
		"INSERT OR IGNORE INTO todo_tags (todo_id, tag_id) VALUES (?, ?)",
		todoID, tagID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *SQLiteTodoRepository) RemoveTagFromTodo(todoID int64, tagName string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var tagID int64
	err = tx.QueryRow("SELECT id FROM tags WHERE name = ?", tagName).Scan(&tagID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		"DELETE FROM todo_tags WHERE todo_id = ? AND tag_id = ?",
		todoID, tagID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *SQLiteTodoRepository) GetTodoTags(todoID int64) ([]string, error) {
	rows, err := r.db.Query(`
        SELECT t.name
        FROM tags t
        JOIN todo_tags tt ON t.id = tt.tag_id
        WHERE tt.todo_id = ?
        ORDER BY t.name
    `, todoID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tags []string
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tags, nil
}

// FindTodosByTag returns todos with the specified tag
func (r *SQLiteTodoRepository) FindTodosByTag(tagName string) ([]*models.Todo, error) {
	rows, err := r.db.Query(`
        SELECT t.id, t.title, t.description, t.status, t.created_at, t.updated_at, t.due_date, t.priority, t.archived
        FROM todos t
        JOIN todo_tags tt ON t.id = tt.todo_id
        JOIN tags tag ON tt.tag_id = tag.id
        WHERE tag.name = ?
    `, tagName)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var todos []*models.Todo
	for rows.Next() {
		todo := &models.Todo{}
		var dueDate sql.NullTime

		if err := rows.Scan(
			&todo.ID,
			&todo.Title,
			&todo.Description,
			&todo.Status,
			&todo.CreatedAt,
			&todo.UpdatedAt,
			&dueDate,
			&todo.Priority,
			&todo.Archived,
		); err != nil {
			return nil, err
		}

		if dueDate.Valid {
			todo.DueDate = &dueDate.Time
		}

		// Get tags for this todo
		tags, err := r.GetTodoTags(todo.ID)
		if err != nil {
			return nil, err
		}
		todo.Tags = tags

		todos = append(todos, todo)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return todos, nil
}

func initSchema(db *sql.DB) error {
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS todos (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            title TEXT NOT NULL,
            description TEXT,
            status INTEGER NOT NULL,
            created_at TIMESTAMP NOT NULL,
            updated_at TIMESTAMP NOT NULL,
            due_date TIMESTAMP,
            priority INTEGER DEFAULT 0,
            archived BOOLEAN DEFAULT 0
        )
    `)
	if err != nil {
		return err
	}

	// Create tags table
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS tags (
            id INTEGER PRIMARY KEY AUTOINCREMENT,
            name TEXT NOT NULL UNIQUE
        )
    `)
	if err != nil {
		return err
	}

	// Create junction table
	_, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS todo_tags (
            todo_id INTEGER NOT NULL,
            tag_id INTEGER NOT NULL,
            PRIMARY KEY (todo_id, tag_id),
            FOREIGN KEY (todo_id) REFERENCES todos(id) ON DELETE CASCADE,
            FOREIGN KEY (tag_id) REFERENCES tags(id) ON DELETE CASCADE
        )
    `)
	if err != nil {
		return err
	}

	// Create indexes
	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_todo_tags_todo_id ON todo_tags(todo_id)`)
	if err != nil {
		return err
	}

	_, err = db.Exec(`CREATE INDEX IF NOT EXISTS idx_todo_tags_tag_id ON todo_tags(tag_id)`)
	if err != nil {
		return err
	}

	return nil
}
