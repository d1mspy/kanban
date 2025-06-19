package main

import (
	"kanban/internal/config"
	"kanban/internal/postgres"
	"kanban/internal/server"
	"time"
)

func main() {
	config.Load()
	time.Sleep(1 * time.Second)

	db := postgres.NewPostgres()
	defer db.Close()

	s := server.New(config.Get().Host)
	s.NewAPI(db)
	s.Start()
}
