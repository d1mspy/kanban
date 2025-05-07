package main

import (
	"kanban/internal/config"
	"kanban/internal/db"
	"kanban/internal/server"
)

func main() {
	db := db.NewPostgres()
	defer db.Close()

	s := server.New(config.Load().Host)

	s.Start()
}