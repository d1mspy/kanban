package main

import (
	"kanban/internal/config"
	"kanban/internal/postgres"
	"kanban/internal/server"
)

func main() {
	db := postgres.NewPostgres()
	defer db.Close()

	s := server.New(config.Load().Host, db)

	s.Start()
}