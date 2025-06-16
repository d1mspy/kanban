package auth

import (
	"database/sql"
	authHandler "kanban/internal/auth/handler"
	authRepo "kanban/internal/auth/repo"
	authService "kanban/internal/auth/service"

	"github.com/gin-gonic/gin"
)

func Init(db *sql.DB, grp *gin.RouterGroup) {
	repo := authRepo.NewRepository(db)
	service := authService.NewService(repo)
	handler := authHandler.NewHandler(service)

	grp.POST("/register", handler.RegisterHandler())
	grp.POST("/login", handler.LoginHandler())

}