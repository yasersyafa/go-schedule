package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func Connect(databaseURL string) (*sqlx.DB, error) {
	conn, err := sqlx.Connect("postgres", databaseURL)

	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}

	return conn, nil
}