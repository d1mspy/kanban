package task

import (
	"database/sql"
	taskHandler "kanban/internal/task/handler"
	taskRepo "kanban/internal/task/repo"
	taskService "kanban/internal/task/service"

	"github.com/gin-gonic/gin"
)

func Init(db *sql.DB, r *gin.Engine, grp *gin.RouterGroup) {
	repo := taskRepo.NewRepository(db)
	serv := taskService.NewService(repo)
	handl := taskHandler.NewHandler(serv)

	grp.POST("/columns/:id/tasks", handl.CreateTaskHandler())
	grp.GET("/columns/:id/tasks", handl.GetAllTasksHandler())
	grp.GET("/tasks/:id", handl.GetTaskHandler())
	grp.PATCH("/tasks/:id", handl.UpdateTaskHandler())
	grp.DELETE("/tasks/:id", handl.DeleteTaskHandler())
}