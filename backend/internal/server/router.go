package server

import (
	"database/sql"
	"kanban/internal/auth"
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

	return engine
}

func (r *Server) Start() {
	r.newAPI().Run(r.host)
}