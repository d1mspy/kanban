package board

import (
	"database/sql"
	boardHandler "kanban/internal/board/handler"
	boardProxy "kanban/internal/board/proxy"
	boardRepo "kanban/internal/board/repo"
	boardService "kanban/internal/board/service"

	"github.com/gin-gonic/gin"
)

func Init(db *sql.DB, r *gin.Engine, grp *gin.RouterGroup) {
	repo := boardRepo.NewRepository(db)
	serv := boardService.NewService(repo)
	proxy := boardProxy.NewProxy(serv, repo)
	handl := boardHandler.NewHandler(proxy)

	grp.POST("/boards", handl.CreateBoardHandler())
	grp.GET("/boards", handl.GetAllBoardsHandler())
	grp.GET("/boards/:id", handl.GetBoardHandler())
	grp.PUT("/boards/:id", handl.UpdateBoardHandler())
	grp.DELETE("/boards/:id", handl.DeleteBoardHandler())

}
