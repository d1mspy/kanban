package columnHandler

import (
	"errors"
	authctx "kanban/internal/auth/context"
	columnModel "kanban/internal/column/model"
	columnProxy "kanban/internal/column/proxy"
	columnRepo "kanban/internal/column/repo"
	columnService "kanban/internal/column/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Proxy interface {
	CreateColumn(boardID, userID string, req columnModel.CreateRequest) error
	GetAllColumns(boardID, userID string) ([]columnModel.Column, error)
	GetColumn(columnID, userID string) (*columnModel.Column, error)
	UpdateColumn(columnID, userID string, req columnModel.UpdateRequest) error
	DeleteColumn(columnID, userID string) error
}

type Handler struct {
	proxy Proxy
}

func NewHandler(serv Proxy) *Handler {
	return &Handler{proxy: serv}
}

func (h *Handler) CreateColumnHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req columnModel.CreateRequest
		if err := ctx.ShouldBindBodyWithJSON(&req); err != nil {
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

		boardID := ctx.Param("id")

		err := h.proxy.CreateColumn(boardID, userID, req);
		if err != nil {
			log.Printf("Failed to create column: %v", err)
			h.handleError(ctx, err, "Failed to create column")
			return
		}

		ctx.Status(http.StatusCreated)
	}
}

func (h *Handler) GetAllColumnsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		boardID := ctx.Param("id")

		userID, ok := authctx.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		columns, err := h.proxy.GetAllColumns(boardID, userID)
		if err != nil {
			log.Printf("Failed to get columns: %v", err)
			h.handleError(ctx, err, "Failed to get columns")
			return
		}

		ctx.JSON(http.StatusOK, columns)
	}
}

func (h *Handler) GetColumnHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		userID, ok := authctx.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		column, err := h.proxy.GetColumn(id, userID)
		if err != nil {
			log.Printf("Failed to get column: %v", err)
			h.handleError(ctx, err, "Failed to get column")
			return
		}

		ctx.JSON(http.StatusOK, column)
	}
}

func (h *Handler) UpdateColumnHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req columnModel.UpdateRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"detail": "Invalid request body",
			})
			return
		}

		id := ctx.Param("id")

		userID, ok := authctx.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		err := h.proxy.UpdateColumn(id, userID, req);
		if err != nil {
			log.Printf("Failed to update column: %v", err)
			h.handleError(ctx, err, "Failed to update column")
			return
		}

		ctx.Status(http.StatusOK)
	}
}

func (h *Handler) DeleteColumnHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		userID, ok := authctx.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		err := h.proxy.DeleteColumn(id, userID);
		if err != nil {
			log.Printf("Failed to delete column: %v", err)
			h.handleError(ctx, err, "Failed to delete column")
			return
		}

		ctx.Status(http.StatusOK)
	}
}

func (h *Handler) handleError(ctx *gin.Context, err error, message string) {
	switch {
	case errors.Is(err, columnProxy.ErrForbidden):
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"detail": "Access denied",
		})
	case errors.Is(err, columnService.ErrColumnNotFound):
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"detail": "Column not found",
		})
	case errors.Is(err, columnRepo.ErrColumnLimitReached):
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"detail": "Column limit reached",
		})
	case errors.Is(err, columnRepo.ErrIncorrectPosition):
		ctx.AbortWithStatusJSON(http.StatusUnprocessableEntity, gin.H{
			"detail": "Column position is greater than possible or not positive",
		})
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"detail": message,
		})
	}
}