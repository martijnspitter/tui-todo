package repository

import (
	"database/sql"
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
			Name: "Add archived field to todos",
			RunSQL: func(tx *sql.Tx) error {
				// Check if the column already exists to make this idempotent
				var columnCount int
				err := tx.QueryRow(`
					SELECT COUNT(*) FROM pragma_table_info('todos') WHERE name = 'archived'
				`).Scan(&columnCount)

				if err != nil {
					return err
				}

				if columnCount == 0 {
					// Add the archived column
					_, err = tx.Exec(`ALTER TABLE todos ADD COLUMN archived BOOLEAN DEFAULT 0`)
					if err != nil {
						return err
					}
				}

				// Migrate existing data - set archived=1 for todos with status=3 (archived)
				_, err = tx.Exec(`
					UPDATE todos SET archived = 1, status = 2 WHERE status = 3
				`)

				return err
			},
			Rollback: func(tx *sql.Tx) error {
				// Convert archived todos back to status 3
				_, err := tx.Exec(`
					UPDATE todos SET status = 3 WHERE archived = 1
				`)
				if err != nil {
					return err
				}

				// SQLite doesn't support dropping columns directly,
				// so a full rollback would require recreating the table.
				// For simplicity, we'll just leave the column.
				return nil
			},
		},
		// Add more migrations here in the future
	}
}
