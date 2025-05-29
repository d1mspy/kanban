package board

import (
	"database/sql"
	boardHandler "kanban/internal/board/handler"
	boardRepo "kanban/internal/board/repo"
	boardService "kanban/internal/board/service"

	"github.com/gin-gonic/gin"
)

func Init(db *sql.DB, r *gin.Engine, grp *gin.RouterGroup) {
	repo := boardRepo.NewRepository(db)
	serv := boardService.NewService(repo)
	handl := boardHandler.NewHandler(serv)

	grp.POST("/boards", handl.CreateBoardHandler())
	grp.GET("/boards", handl.GetAllBoardsHandler())
	grp.GET("/boards/:id", handl.GetBoardHandler())
	grp.PUT("/boards/:id", handl.UpdateBoardHandler())
	grp.DELETE("/boards/:id", handl.DeleteBoardHandler())

}
