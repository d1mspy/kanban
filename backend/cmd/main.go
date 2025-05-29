package main

import (
	"kanban/internal/config"
	"kanban/internal/postgres"
	"kanban/internal/server"
	"time"
)

func main() {
	time.Sleep(1*time.Second)

	db := postgres.NewPostgres()
	defer db.Close()

	s := server.New(config.Load().Host, db)

	s.Start()
}