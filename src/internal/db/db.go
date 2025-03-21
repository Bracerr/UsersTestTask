package db

import (
	"database/sql"
	"fmt"
	"users-api/src/internal/config"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewDB(cfg *config.Config) (*DB, error) {
	db, err := sql.Open("postgres", cfg.GetDBURL())
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the database: %w", err)
	}

	database := &DB{
		DB: db,
	}

	if err := RunMigrations(database); err != nil {
		return nil, fmt.Errorf("error running migrations: %w", err)
	}

	return database, nil
}
