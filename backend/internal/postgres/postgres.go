package postgres

import (
	"database/sql"
	"kanban/internal/config"
	"log"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/lib/pq"
	_ "github.com/golang-migrate/migrate/v4/source/file"
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
		} else {
			log.Println("DB is up, running migrations...")
			runMigrations(db)
			break
		}
		time.Sleep(dbConnectRetryDelay)
	}
	if err != nil {
		log.Fatalf("DB setup failed after %d attempts: %v", maxDBConnectAttempts, err)
	}

	return db
}

func runMigrations(db *sql.DB) {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to init migration driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://db/migrations",
		config.Get().DBname,
		driver,
	)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Migration failed: %v", err)
	}
}