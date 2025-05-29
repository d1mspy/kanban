package column

import (
	"database/sql"
	columnHandler "kanban/internal/column/handler"
	columnRepo "kanban/internal/column/repo"
	columnService "kanban/internal/column/service"

	"github.com/gin-gonic/gin"
)

func Init(db *sql.DB, r *gin.Engine, grp *gin.RouterGroup) {
	repo := columnRepo.NewRepository(db)
	serv := columnService.NewService(repo)
	handl := columnHandler.NewHandler(serv)

	grp.POST("/boards/:id/columns", handl.CreateColumnHandler())
	grp.GET("/boards/:id/columns", handl.GetAllColumnsHandler())
	grp.GET("/columns/:id", handl.GetColumnHandler())
	grp.PATCH("/columns/:id", handl.UpdateColumnHandler())
	grp.DELETE("/columns/:id", handl.DeleteColumnHandler())
}