package server

import (
	"database/sql"
	"kanban/internal/auth"
	authMiddleware "kanban/internal/auth/middleware"
	"kanban/internal/board"
	"kanban/internal/column"
	"kanban/internal/task"

	"github.com/gin-gonic/gin"
)

type Server struct {
	host   string
	engine *gin.Engine
}

func New(host string) *Server {
	s := &Server{
		host:   host,
		engine: gin.New(),
	}

	return s
}

func (r *Server) NewAPI(db *sql.DB) {
	authGroup := r.engine.Group("/auth")
	protectedGroup := r.engine.Group("/", authMiddleware.Middleware())

	auth.Init(db, authGroup)

	board.Init(db, protectedGroup)
	column.Init(db, protectedGroup)
	task.Init(db, protectedGroup)
}

func (r *Server) Start() {
	r.engine.Run(r.host)
}
