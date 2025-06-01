package repository

import (
	"database/sql"
	"fmt"
	"time"
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
		{
			ID:   3,
			Name: "Add tag description and timestamp fields",
			RunSQL: func(tx *sql.Tx) error {
				// Check if columns already exist to avoid errors
				var descriptionExists, createdAtExists, updatedAtExists int

				err := tx.QueryRow(`
                    SELECT COUNT(*) FROM pragma_table_info('tags')
                    WHERE name = 'description'
                `).Scan(&descriptionExists)
				if err != nil {
					return fmt.Errorf("failed to check for description column: %w", err)
				}

				err = tx.QueryRow(`
                    SELECT COUNT(*) FROM pragma_table_info('tags')
                    WHERE name = 'created_at'
                `).Scan(&createdAtExists)
				if err != nil {
					return fmt.Errorf("failed to check for created_at column: %w", err)
				}

				err = tx.QueryRow(`
                    SELECT COUNT(*) FROM pragma_table_info('tags')
                    WHERE name = 'updated_at'
                `).Scan(&updatedAtExists)
				if err != nil {
					return fmt.Errorf("failed to check for updated_at column: %w", err)
				}

				// Add description column if it doesn't exist
				if descriptionExists == 0 {
					_, err := tx.Exec(`ALTER TABLE tags ADD COLUMN description TEXT`)
					if err != nil {
						return fmt.Errorf("failed to add description column: %w", err)
					}
				}

				// For SQLite, we can't add columns with DEFAULT CURRENT_TIMESTAMP
				// So we add them without defaults and then update values

				// Add created_at column if it doesn't exist
				if createdAtExists == 0 {
					_, err := tx.Exec(`ALTER TABLE tags ADD COLUMN created_at TIMESTAMP`)
					if err != nil {
						return fmt.Errorf("failed to add created_at column: %w", err)
					}

					// Set current timestamp for all existing rows
					now := time.Now().Format("2006-01-02 15:04:05")
					_, err = tx.Exec(`UPDATE tags SET created_at = ? WHERE created_at IS NULL`, now)
					if err != nil {
						return fmt.Errorf("failed to set created_at values: %w", err)
					}
				}

				// Add updated_at column if it doesn't exist
				if updatedAtExists == 0 {
					_, err := tx.Exec(`ALTER TABLE tags ADD COLUMN updated_at TIMESTAMP`)
					if err != nil {
						return fmt.Errorf("failed to add updated_at column: %w", err)
					}

					// Set current timestamp for all existing rows
					now := time.Now().Format("2006-01-02 15:04:05")
					_, err = tx.Exec(`UPDATE tags SET updated_at = ? WHERE updated_at IS NULL`, now)
					if err != nil {
						return fmt.Errorf("failed to set updated_at values: %w", err)
					}
				}

				return nil
			},
		},
	}
}
