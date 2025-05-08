package postgres

import (
	"database/sql"
	"kanban/internal/config"
	"log"

	_ "github.com/lib/pq"
)

func NewPostgres() *sql.DB {
	db, err := sql.Open("postgres", config.Load().PostgresURI)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal("Failed to ping DB:", err)
	}

	if err = InitDatabase(db); err != nil {
		log.Fatal("Failed to init DB:", err)
	}

	return db
}

func InitDatabase(db *sql.DB) error {
	queries := []string{
		QueryCreateUserTable,
		QueryCreateBoardTable,
		QueryCreateColumnTable,
		QueryCreateTaskTable,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}

	return nil
}