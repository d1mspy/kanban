package main

import (
	"kanban/internal/config"
	"kanban/internal/postgres"
	"kanban/internal/server"
)

func main() {
	config.Load()

	db := postgres.NewPostgres()
	defer db.Close()

	s := server.New(config.Get().Host)
	s.NewAPI(db)
	s.Start()
}
