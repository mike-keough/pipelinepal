package db

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	*sql.DB
	path string
}

func Open(path string) (*DB, error) {
	// Busy timeout helps with “database is locked” during fast UI operations
	dsn := fmt.Sprintf("file:%s?_foreign_keys=on&_busy_timeout=5000", path)
	sqldb, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, err
	}
	sqldb.SetConnMaxLifetime(30 * time.Minute)
	sqldb.SetMaxOpenConns(1) // SQLite is happiest with 1 writer
	sqldb.SetMaxIdleConns(1)

	if err := sqldb.Ping(); err != nil {
		_ = sqldb.Close()
		return nil, err
	}
	return &DB{DB: sqldb, path: path}, nil
}

func (d *DB) Close() error {
	if d.DB == nil {
		return nil
	}
	return d.DB.Close()
}

func (d *DB) Migrate(ctx context.Context) error {
	return runMigrations(ctx, d.DB)
}
