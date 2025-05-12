package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/charmbracelet/log"
)

// Migration represents a database migration
type Migration struct {
	ID       int
	Name     string
	RunSQL   func(tx *sql.Tx) error
	Rollback func(tx *sql.Tx) error // Optional rollback function
}

// MigrationManager handles database migrations
type MigrationManager struct {
	db *sql.DB
}

// NewMigrationManager creates a new migration manager
func NewMigrationManager(db *sql.DB) *MigrationManager {
	return &MigrationManager{db: db}
}

// Initialize creates the migrations table if it doesn't exist
func (m *MigrationManager) Initialize() error {
	_, err := m.db.Exec(`
		CREATE TABLE IF NOT EXISTS schema_migrations (
			id INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			applied_at TIMESTAMP NOT NULL
		)
	`)
	return err
}

// GetAppliedMigrations returns the IDs of all applied migrations
func (m *MigrationManager) GetAppliedMigrations() (map[int]bool, error) {
	rows, err := m.db.Query("SELECT id FROM schema_migrations")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	appliedMigrations := make(map[int]bool)
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		appliedMigrations[id] = true
	}

	return appliedMigrations, rows.Err()
}

// ApplyMigrations applies all pending migrations
func (m *MigrationManager) ApplyMigrations(migrations []Migration) error {
	appliedMigrations, err := m.GetAppliedMigrations()
	if err != nil {
		return fmt.Errorf("couldn't get applied migrations: %w", err)
	}

	for _, migration := range migrations {
		// Skip if already applied, but validate critical migrations
		if appliedMigrations[migration.ID] {
			log.Debug("Migration already applied", "id", migration.ID, "name", migration.Name)

			// For critical migrations, validate they actually applied correctly
			valid, err := m.ValidateMigration(migration.ID)
			if err != nil {
				log.Warn("Failed to validate migration", "id", migration.ID, "error", err)
				// Continue anyway, since this is just a validation check
				continue
			}

			if !valid {
				log.Warn("Migration marked as applied but validation failed - removing record",
					"id", migration.ID, "name", migration.Name)

				// Remove the invalid migration record
				_, err := m.db.Exec("DELETE FROM schema_migrations WHERE id = ?", migration.ID)
				if err != nil {
					return fmt.Errorf("couldn't remove invalid migration record: %w", err)
				}

				// Now continue to re-apply this migration
			} else {
				// Migration is valid, continue to next one
				continue
			}
		}

		log.Info("Applying migration", "id", migration.ID, "name", migration.Name)

		// Start a transaction for this migration
		tx, err := m.db.Begin()
		if err != nil {
			return fmt.Errorf("couldn't start transaction for migration %d: %w", migration.ID, err)
		}

		// Run the migration
		if err := migration.RunSQL(tx); err != nil {
			tx.Rollback()
			return fmt.Errorf("migration %d failed: %w", migration.ID, err)
		}

		// For critical migrations, validate before committing
		if migration.ID == 2 { // Time tracking migration
			var timeSpentExists, timeStartedExists int
			err := tx.QueryRow("SELECT COUNT(*) FROM pragma_table_info('todos') WHERE name = 'time_spent'").Scan(&timeSpentExists)
			if err != nil || timeSpentExists == 0 {
				tx.Rollback()
				return fmt.Errorf("migration failed: time_spent column not created")
			}

			err = tx.QueryRow("SELECT COUNT(*) FROM pragma_table_info('todos') WHERE name = 'time_started'").Scan(&timeStartedExists)
			if err != nil || timeStartedExists == 0 {
				tx.Rollback()
				return fmt.Errorf("migration failed: time_started column not created")
			}
		}

		// Record that this migration was applied
		_, err = tx.Exec(
			"INSERT INTO schema_migrations (id, name, applied_at) VALUES (?, ?, ?)",
			migration.ID,
			migration.Name,
			time.Now(),
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("couldn't record migration %d: %w", migration.ID, err)
		}

		// Commit the transaction
		if err := tx.Commit(); err != nil {
			return fmt.Errorf("couldn't commit migration %d: %w", migration.ID, err)
		}

		log.Info("Migration applied successfully", "id", migration.ID)
	}

	return nil
}

func (m *MigrationManager) ColumnExists(table, column string) (bool, error) {
	var exists bool
	query := `SELECT COUNT(*) > 0 FROM pragma_table_info(?) WHERE name = ?`
	err := m.db.QueryRow(query, table, column).Scan(&exists)
	return exists, err
}

func (m *MigrationManager) ValidateMigration(migrationID int) (bool, error) {
	switch migrationID {
	case 2: // Time tracking migration
		timeSpentExists, err := m.ColumnExists("todos", "time_spent")
		if err != nil {
			return false, err
		}
		timeStartedExists, err := m.ColumnExists("todos", "time_started")
		if err != nil {
			return false, err
		}
		return timeSpentExists && timeStartedExists, nil
	default:
		// For other migrations, just assume they're valid
		return true, nil
	}
}
