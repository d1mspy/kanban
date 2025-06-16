package taskHandler

import (
	authMiddleware "kanban/internal/auth/middleware"
	taskModel "kanban/internal/task/model"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Proxy interface {
	CreateTask(columnID, name, description, userID string) error
	GetAllTasks(columnID, userID string) ([]taskModel.Task, error)
	GetTask(taskID, userID string) (*taskModel.Task, error)
	UpdateTask(req taskModel.UpdateTaskRequest, taskID, userID string) error
	DeleteTask(taskID, userID string) error
}

type Handler struct {
	proxy Proxy
}

func NewHandler(proxy Proxy) *Handler {
	return &Handler{proxy: proxy}
}

func (h *Handler) CreateTaskHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req taskModel.CreateTaskRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"detail": "Invalid request body",
			})
			return
		}

		userID, ok := authMiddleware.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		columnID := ctx.Param("id")

		err := h.proxy.CreateTask(columnID, req.Name, req.Description, userID); 
		if err != nil {
			log.Printf("Failed to create task: %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to create task",
			})
			return
		} 

		ctx.Status(http.StatusCreated)
	}
}

func (h *Handler) GetAllTasksHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		columnID := ctx.Param("id")

		userID, ok := authMiddleware.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		tasks, err := h.proxy.GetAllTasks(columnID, userID)
		if err != nil {
			log.Printf("Failed to get all tasks: %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to get all tasks",
			})
			return
		}

		ctx.JSON(http.StatusOK, tasks)
	}
}

func (h *Handler) GetTaskHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		taskID := ctx.Param("id")

		userID, ok := authMiddleware.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		task, err := h.proxy.GetTask(taskID, userID)
		if err != nil {
			log.Printf("Failed to get task: %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to get task",
			})
			return
		}

		ctx.JSON(http.StatusOK, task)
	}
}

func (h *Handler) UpdateTaskHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req taskModel.UpdateTaskRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"detail": "Invalid request body",
			})
			return
		}

		taskID := ctx.Param("id")

		userID, ok := authMiddleware.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		err := h.proxy.UpdateTask(req, taskID, userID)
		if err != nil {
			log.Printf("Failed to update task: %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to update task",
			})
			return
		}
		
		ctx.Status(http.StatusOK)
	}
}

func (h *Handler) DeleteTaskHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		taskID := ctx.Param("id")

		userID, ok := authMiddleware.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		err := h.proxy.DeleteTask(taskID, userID);
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to delete task",
			})
			return
		}

		ctx.Status(http.StatusOK)
	}
}