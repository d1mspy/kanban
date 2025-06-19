package postgres

import (
	"database/sql"
	"kanban/internal/config"
	"log"
	"time"

	_ "github.com/lib/pq"
)

var (
	maxDBConnectAttempts = 10
	dbConnectRetryDelay = 1 * time.Second
)

func NewPostgres() *sql.DB {
	db, err := sql.Open("postgres", config.Get().PostgresURI)
	if err != nil {
		log.Fatal("Failed to open DB:", err)
	}
	
	for i := 1; i <= maxDBConnectAttempts; i++ {
		if err = db.Ping(); err != nil {
			log.Printf("Attempt %d: failed to ping DB: %v", i, err)
		} else if err = InitDatabase(db); err != nil {
			log.Printf("Attempt %d: failed to init DB: %v", i, err)
		} else {
			log.Println("DB is ready")
			break
		}
		time.Sleep(dbConnectRetryDelay)
	}
	if err != nil {
		log.Fatalf("DB setup failed after %d attempts: %v", maxDBConnectAttempts, err)
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