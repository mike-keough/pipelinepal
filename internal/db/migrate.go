package db

import (
	"database/sql"
	"embed"
)

//go:embed ../../migrations/*.sql
var migrationsFS embed.FS

func Migrate(db *sql.DB) error {
	// For MVP, just run the init file.
	// Later: track applied migrations in a schema_migrations table.
	b, err := migrationsFS.ReadFile("../../migrations/001_init.sql")
	if err != nil {
		return err
	}
	_, err = db.Exec(string(b))
	return err
}
