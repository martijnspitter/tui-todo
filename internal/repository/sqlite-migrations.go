package repository

import (
	"database/sql"
	"fmt"
)

// GetAllMigrations returns all migrations in order
func GetAllMigrations() []Migration {
	return []Migration{
		{
			ID:   1,
			Name: "Initial schema",
			RunSQL: func(tx *sql.Tx) error {
				// This would be our baseline schema that's already created in initSchema
				// We don't need to run anything here since we assume the schema is already initialized
				// Just a placeholder for tracking purposes
				return nil
			},
		},
		{
			ID:   2,
			Name: "Add time tracking fields",
			RunSQL: func(tx *sql.Tx) error {
				// First check if columns already exist to avoid errors
				var timeSpentExists, timeStartedExists int

				err := tx.QueryRow(`
					SELECT COUNT(*) FROM pragma_table_info('todos')
					WHERE name = 'time_spent'
				`).Scan(&timeSpentExists)
				if err != nil {
					return fmt.Errorf("failed to check for time_spent column: %w", err)
				}

				err = tx.QueryRow(`
					SELECT COUNT(*) FROM pragma_table_info('todos')
					WHERE name = 'time_started'
				`).Scan(&timeStartedExists)
				if err != nil {
					return fmt.Errorf("failed to check for time_started column: %w", err)
				}

				// Add time_spent column if it doesn't exist
				if timeSpentExists == 0 {
					_, err := tx.Exec(`ALTER TABLE todos ADD COLUMN time_spent INTEGER DEFAULT 0 NOT NULL`)
					if err != nil {
						return fmt.Errorf("failed to add time_spent column: %w", err)
					}
				}

				// Add time_started column if it doesn't exist
				if timeStartedExists == 0 {
					_, err := tx.Exec(`ALTER TABLE todos ADD COLUMN time_started TIMESTAMP NULL`)
					if err != nil {
						return fmt.Errorf("failed to add time_started column: %w", err)
					}
				}

				return nil
			},
		},
	}
}
