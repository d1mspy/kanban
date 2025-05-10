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
	host string
	db *sql.DB
}

func New(host string, db *sql.DB) *Server {
	s := &Server{
		host: host,
		db: db,	
	}

	return s
}

func (r *Server) newAPI() *gin.Engine {
	engine := gin.New()

	engine.GET("/health", func(ctx *gin.Context) {
		ctx.Status(http.StatusNoContent)
	})

	authGroup := engine.Group("/auth")
	authGroup.POST("/register", auth.RegisterHandler(r.db))
	authGroup.POST("/login", auth.LoginHandler(r.db))

	protectedGroup := engine.Group("/", auth.AuthMiddleware())

	protectedGroup.GET("/me", func(ctx *gin.Context) {
		userID, _ := auth.GetUserID(ctx)
		ctx.JSON(http.StatusOK, gin.H{"user_id": userID})
	})

	protectedGroup.POST("/boards", board.CreateBoardHandler(r.db))
	protectedGroup.GET("/boards", board.GetAllBoardsHandler(r.db))
	protectedGroup.GET("/boards/:id", board.GetBoardHandler(r.db))
	protectedGroup.PUT("/boards/:id", board.UpdateBoardHandler(r.db))
	protectedGroup.DELETE("/boards/:id", board.DeleteBoardHandler(r.db))

	protectedGroup.POST("/boards/:id/columns", column.CreateColumnHandler(r.db))
	protectedGroup.GET("/boards/:id/columns", column.GetAllColumnsHandler(r.db))
	protectedGroup.GET("/columns/:id", column.GetColumnHandler(r.db))
	protectedGroup.PATCH("/columns/:id", column.UpdateColumnHandler(r.db))
	protectedGroup.DELETE("/columns/:id", column.DeleteColumnHandler(r.db))

	protectedGroup.POST("/columns/:id/tasks", task.CreateTaskHandler(r.db))
	protectedGroup.GET("/columns/:id/tasks", task.GetAllTasksHandler(r.db))
	protectedGroup.GET("/tasks/:id", task.GetTaskHandler(r.db))
	protectedGroup.PATCH("/tasks/:id", task.UpdateTaskHandler(r.db))
	protectedGroup.DELETE("/tasks/:id", task.DeleteTaskHandler(r.db))

	return engine
}

func (r *Server) Start() {
	r.newAPI().Run(r.host)
}