package task

import (
	"database/sql"
	"kanban/internal/auth"
	"net/http"
	"reflect"
	"time"

	"github.com/gin-gonic/gin"
)

type createTaskRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type updateTaskRequest struct {
	ColumnID    *string     `json:"column_id"`
	Name        *string     `json:"name"`
	Description *string     `json:"description"`
	Position    *int        `json:"position"`
	Done        *bool       `json:"done"`
	Deadline    *time.Time  `json:"deadline"`
}

type updateCase string
const (
	caseColumn = "column"
	casePosition = "position"
	caseContent = "content"
	caseUndefined = "undefined"
)

func CreateTaskHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req createTaskRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"detail": "Invalid request body",
			})
			return
		}

		userID, ok := auth.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		columnID := ctx.Param("id")
		task := Task{
			ColumnID: columnID,
			Name: req.Name,
			Description: req.Description,
		}

		if err := CreateTask(db, task, userID); err != nil {
			switch err {
			case errForbidden:
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"detail": "This is not your board",
				})
				return
			case errTaskLimitReached:
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"detail": "Task limit reached",
				})
				return
			default:
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"detail": "Failed to create task",
				})
				return
			}
		} 

		ctx.Status(http.StatusCreated)
	}
}

func GetAllTasksHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		columnID := ctx.Param("id")

		userID, ok := auth.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		tasks, err := GetAllTasks(db, columnID, userID)
		if err != nil {
			switch err {
			case errForbidden:
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"detail": "This is not your board",
				})
				return
			default:
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"detail": "Failed to get all tasks",
				})
				return
			}
		}

		ctx.JSON(http.StatusOK, tasks)
	}
}

func GetTaskHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		taskID := ctx.Param("id")

		userID, ok := auth.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		task, err := GetTask(db, taskID, userID)
		if err != nil {
			switch err {
			case errTaskNotFound:
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"detail": "Task not found",
				})
				return
			case errForbidden:
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"detail": "This is not your board",
				})
				return
			default:
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"detail": "Failed to get task",
				})
				return
			}
		}

		ctx.JSON(http.StatusOK, task)
	}
}

func UpdateTaskHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req updateTaskRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"detail": "Invalid request body",
			})
			return
		}

		taskID := ctx.Param("id")

		userID, ok := auth.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		updCase := validateUpdateTaskRequest(req)

		var err error
		switch updCase {
		case caseColumn:
			err = UpdateTaskColumn(db, req, taskID, userID)
		case casePosition:
			err = UpdateTaskPosition(db, req, taskID, userID)
		case caseContent:
			err = UpdateTaskContent(db, req, taskID, userID)
		default:
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"detail": "Invalid combination of fields",
			})
			return
		}
		
		if err != nil {
			switch err {
			case errTaskNotFound:
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"detail": "Task not found",
				})
				return
			case errForbidden:
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"detail": "This is not your board",
				})
				return
			case errIncorrectPosition:
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"detail": "Position is greater than possible or not positive",
				})
				return
			default:
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"detail": "Failed to update task",
				})
				return
			}
		}
		
		ctx.Status(http.StatusOK)
	}
}

func DeleteTaskHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		taskID := ctx.Param("id")

		userID, ok := auth.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		if err := DeleteTask(db, taskID, userID); err != nil {
			switch err {
			case errTaskNotFound:
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"detail": "Task not found",
				})
				return
			case errForbidden:
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"detail": "This is not your board",
				})
				return
			default:
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"detail": "Failed to delete task",
				})
				return
			}
		}

		ctx.Status(http.StatusOK)
	}
}

func validateUpdateTaskRequest(req updateTaskRequest) updateCase {
	contentFields := []any{
		req.Name,
		req.Description,
		req.Done,
		req.Deadline,
	}

	contentCount := 0
	for _, field := range contentFields {
		if !isNil(field) {
			contentCount++
		}
	}

	columnSet := req.ColumnID != nil
	positionSet := req.Position != nil

	switch {
	case contentCount == 1 && !columnSet && !positionSet:
		return caseContent
	case contentCount == 0 && columnSet && positionSet:
		return caseColumn
	case contentCount == 0 && !columnSet && positionSet:
		return casePosition
	default:
		return caseUndefined
	}
}

func isNil(v any) bool {
	return v == nil  || reflect.ValueOf(v).IsNil()
}