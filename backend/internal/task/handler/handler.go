package taskHandler

import (
	"errors"
	authctx "kanban/internal/auth/context"
	taskModel "kanban/internal/task/model"
	taskProxy "kanban/internal/task/proxy"
	taskRepo "kanban/internal/task/repo"
	taskService "kanban/internal/task/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Proxy interface {
	CreateTask(columnID, userID string, req taskModel.CreateRequest) error
	GetAllTasks(columnID, userID string) ([]taskModel.Task, error)
	GetTask(taskID, userID string) (*taskModel.Task, error)
	UpdateTask(taskID, userID string, req taskModel.UpdateRequest) error
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
		var req taskModel.CreateRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"detail": "Invalid request body",
			})
			return
		}

		userID, ok := authctx.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		columnID := ctx.Param("id")

		err := h.proxy.CreateTask(columnID, userID, req); 
		if err != nil {
			log.Printf("Failed to create task: %v", err)
			h.handleError(ctx, err, "Failed to create task")
			return
		} 

		ctx.Status(http.StatusCreated)
	}
}

func (h *Handler) GetAllTasksHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		columnID := ctx.Param("id")

		userID, ok := authctx.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		tasks, err := h.proxy.GetAllTasks(columnID, userID)
		if err != nil {
			log.Printf("Failed to get all tasks: %v", err)
			h.handleError(ctx, err, "Failed to get all tasks")
			return
		}

		ctx.JSON(http.StatusOK, tasks)
	}
}

func (h *Handler) GetTaskHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		taskID := ctx.Param("id")

		userID, ok := authctx.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		task, err := h.proxy.GetTask(taskID, userID)
		if err != nil {
			log.Printf("Failed to get task: %v", err)
			h.handleError(ctx, err, "Failed to get task")
			return
		}

		ctx.JSON(http.StatusOK, task)
	}
}

func (h *Handler) UpdateTaskHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req taskModel.UpdateRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"detail": "Invalid request body",
			})
			return
		}

		taskID := ctx.Param("id")

		userID, ok := authctx.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		err := h.proxy.UpdateTask(taskID, userID, req)
		if err != nil {
			log.Printf("Failed to update task: %v", err)
			h.handleError(ctx, err, "Failed to update task")
			return
		}
		
		ctx.Status(http.StatusOK)
	}
}

func (h *Handler) DeleteTaskHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		taskID := ctx.Param("id")

		userID, ok := authctx.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		err := h.proxy.DeleteTask(taskID, userID);
		if err != nil {
			log.Printf("Failed to delete task: %v", err)
			h.handleError(ctx, err, "Failed to delete task")
			return
		}

		ctx.Status(http.StatusOK)
	}
}

func (h *Handler) handleError(ctx *gin.Context, err error, message string) {
	switch {
	case errors.Is(err, taskProxy.ErrForbidden):
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"detail": "Access denied",
		})
	case errors.Is(err, taskService.ErrBadUpdateRequest):
		ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
			"detail": "Invalid combination of fields",
		})
	case errors.Is(err, taskService.ErrTaskNotFound):
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"detail": "Task not found",
		})
	case errors.Is(err, taskRepo.ErrTaskLimitReached):
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"detail": "Task limit reached",
		})
	case errors.Is(err, taskRepo.ErrIncorrectPosition):
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"detail": "Task position is greater than possible or not positive",
		})
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"detail": message,
		})
	}
}