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

	board.Init(db, r.engine)

	//protectedGroup.POST("/boards", board.CreateBoardHandler(db))
	//protectedGroup.GET("/boards", board.GetAllBoardsHandler(db))
	//protectedGroup.GET("/boards/:id", board.GetBoardHandler(db))
	//protectedGroup.PUT("/boards/:id", board.UpdateBoardHandler(db))
	//protectedGroup.DELETE("/boards/:id", board.DeleteBoardHandler(db))

	protectedGroup.POST("/boards/:id/columns", column.CreateColumnHandler(db))
	protectedGroup.GET("/boards/:id/columns", column.GetAllColumnsHandler(db))
	protectedGroup.GET("/columns/:id", column.GetColumnHandler(db))
	protectedGroup.PATCH("/columns/:id", column.UpdateColumnHandler(db))
	protectedGroup.DELETE("/columns/:id", column.DeleteColumnHandler(db))

	protectedGroup.POST("/columns/:id/tasks", task.CreateTaskHandler(db))
	protectedGroup.GET("/columns/:id/tasks", task.GetAllTasksHandler(db))
	protectedGroup.GET("/tasks/:id", task.GetTaskHandler(db))
	protectedGroup.PATCH("/tasks/:id", task.UpdateTaskHandler(db))
	protectedGroup.DELETE("/tasks/:id", task.DeleteTaskHandler(db))
}

func (r *Server) Start() {
	r.engine.Run(r.host)
}
