package task

import (
	"database/sql"
	taskHandler "kanban/internal/task/handler"
	taskProxy "kanban/internal/task/proxy"
	taskRepo "kanban/internal/task/repo"
	taskService "kanban/internal/task/service"

	"github.com/gin-gonic/gin"
)

func Init(db *sql.DB, grp *gin.RouterGroup) {
	repo := taskRepo.NewRepository(db)
	service := taskService.NewService(repo)
	proxy := taskProxy.NewProxy(service)
	handler := taskHandler.NewHandler(proxy)

	grp.POST("/columns/:id/tasks", handler.CreateTaskHandler())
	grp.GET("/columns/:id/tasks", handler.GetAllTasksHandler())
	grp.GET("/tasks/:id", handler.GetTaskHandler())
	grp.PATCH("/tasks/:id", handler.UpdateTaskHandler())
	grp.DELETE("/tasks/:id", handler.DeleteTaskHandler())
}