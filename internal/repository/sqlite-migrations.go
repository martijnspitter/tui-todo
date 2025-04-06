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
	}
}
