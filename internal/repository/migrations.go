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

	// Sort migrations by ID (should be already sorted, but just to be safe)
	// We can add sorting logic here if needed

	for _, migration := range migrations {
		// Skip if already applied
		if appliedMigrations[migration.ID] {
			log.Debug("Migration already applied", "id", migration.ID, "name", migration.Name)
			continue
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
