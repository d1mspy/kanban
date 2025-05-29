package columnHandler

import (
	"kanban/internal/auth"
	columnModel "kanban/internal/column/model"
	columnService "kanban/internal/column/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Service interface {
	CreateColumn(boardID, name, userID string) error
	GetAllColumns(boardID, userID string) ([]columnModel.Column, error)
	GetColumn(id, userID string) (*columnModel.Column, error)
	UpdateColumn(userID, columnID string, newName *string, newPos *int) error
	DeleteColumn(userID, columnID string) error
}

type Handler struct {
	service *columnService.Service
}

func NewHandler(serv *columnService.Service) *Handler {
	return &Handler{service: serv}
}

func (h *Handler) CreateColumnHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req columnModel.ColumnRequest
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
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

		boardID := ctx.Param("id")

		err := h.service.CreateColumn(boardID, *req.Name, userID);
		if err != nil {
			log.Printf("Failed to create column: %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to create column",
			})
			return
		}

		ctx.Status(http.StatusCreated)
	}
}

func (h *Handler) GetAllColumnsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		boardID := ctx.Param("id")

		userID, ok := auth.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		columns, err := h.service.GetAllColumns(boardID, userID)
		if err != nil {
			log.Printf("Failed to get columns: %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to get all columns",
			})
			return
		}

		ctx.JSON(http.StatusOK, columns)
	}
}

func (h *Handler) GetColumnHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		userID, ok := auth.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		column, err := h.service.GetColumn(id, userID)
		if err != nil {
			log.Printf("Failed to get column: %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to get column",
			})
			return
		}

		ctx.JSON(http.StatusOK, column)
	}
}

func (h *Handler) UpdateColumnHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req columnModel.ColumnRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"detail": "Invalid request body",
			})
			return
		}

		id := ctx.Param("id")

		userID, ok := auth.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		err := h.service.UpdateColumn(userID, id, req.Name, req.Position);
		if err != nil {
			log.Printf("Failed to update column: %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to update column",
			})
			return
		}

		ctx.Status(http.StatusOK)
	}
}

func (h *Handler) DeleteColumnHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		userID, ok := auth.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		err := h.service.DeleteColumn(userID, id);
		if err != nil {
			log.Printf("Failed to delete column: %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to delete column",
			})
			return
		}

		ctx.Status(http.StatusOK)
	}
}