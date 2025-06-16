package board

import (
	"database/sql"
	boardHandler "kanban/internal/board/handler"
	boardProxy "kanban/internal/board/proxy"
	boardRepo "kanban/internal/board/repo"
	boardService "kanban/internal/board/service"

	"github.com/gin-gonic/gin"
)

func Init(db *sql.DB, grp *gin.RouterGroup) {
	repo := boardRepo.NewRepository(db)
	service := boardService.NewService(repo)
	proxy := boardProxy.NewProxy(service)
	handler := boardHandler.NewHandler(proxy)

	grp.POST("/boards", handler.CreateBoardHandler())
	grp.GET("/boards", handler.GetAllBoardsHandler())
	grp.GET("/boards/:id", handler.GetBoardHandler())
	grp.PUT("/boards/:id", handler.UpdateBoardHandler())
	grp.DELETE("/boards/:id", handler.DeleteBoardHandler())

}
