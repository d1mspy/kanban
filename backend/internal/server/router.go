package server

import (
	"database/sql"
	"kanban/internal/auth"
	"kanban/internal/board"
	"kanban/internal/column"
	"kanban/internal/task"
	"net/http"

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
	r.engine.GET("/health", func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})

	authGroup := r.engine.Group("/auth")
	authGroup.POST("/register", auth.RegisterHandler(db))
	authGroup.POST("/login", auth.LoginHandler(db))

	protectedGroup := r.engine.Group("/", auth.AuthMiddleware())

	protectedGroup.GET("/me", func(ctx *gin.Context) {
		userID, _ := auth.GetUserID(ctx)
		ctx.JSON(http.StatusOK, gin.H{"user_id": userID})
	})

	board.Init(db, r.engine, protectedGroup)
	column.Init(db, r.engine, protectedGroup)
	task.Init(db, r.engine, protectedGroup)
}

func (r *Server) Start() {
	r.engine.Run(r.host)
}
