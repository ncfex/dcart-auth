package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

type database struct {
	*sql.DB
}

func NewDatabase(dsn string) (*database, error) {
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error initializing database: %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	return &database{db}, nil
}
