package board

import (
	"database/sql"
	"kanban/internal/auth"
	boardHandler "kanban/internal/board/handler"
	boardRepo "kanban/internal/board/repo"
	boardService "kanban/internal/board/service"

	"github.com/gin-gonic/gin"
)

func Init(db *sql.DB, r *gin.Engine) {
	repo := boardRepo.NewRepository(db)
	serv := boardService.NewService(repo)
	handl := boardHandler.NewHandler(serv)

	rp := r.Group("/", auth.AuthMiddleware())

	rp.POST("/boards", handl.CreateBoardHandler())
	rp.GET("/boards", handl.GetAllBoardsHandler())
	rp.GET("/boards/:id", handl.GetBoardHandler())
	rp.PUT("/boards/:id", handl.UpdateBoardHandler())
	rp.DELETE("/boards/:id", handl.DeleteBoardHandler())

}
