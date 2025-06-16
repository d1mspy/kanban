package column

import (
	"database/sql"
	columnHandler "kanban/internal/column/handler"
	columnProxy "kanban/internal/column/proxy"
	columnRepo "kanban/internal/column/repo"
	columnService "kanban/internal/column/service"

	"github.com/gin-gonic/gin"
)

func Init(db *sql.DB, grp *gin.RouterGroup) {
	repo := columnRepo.NewRepository(db)
	service := columnService.NewService(repo)
	proxy := columnProxy.NewProxy(service)
	handler := columnHandler.NewHandler(proxy)

	grp.POST("/boards/:id/columns", handler.CreateColumnHandler())
	grp.GET("/boards/:id/columns", handler.GetAllColumnsHandler())
	grp.GET("/columns/:id", handler.GetColumnHandler())
	grp.PATCH("/columns/:id", handler.UpdateColumnHandler())
	grp.DELETE("/columns/:id", handler.DeleteColumnHandler())
}