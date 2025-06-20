package boardHandler

import (
	"errors"
	authctx "kanban/internal/auth/context"
	boardModel "kanban/internal/board/model"
	boardProxy "kanban/internal/board/proxy"
	boardService "kanban/internal/board/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Proxy interface {
	CreateBoard(userID string, req boardModel.Request) error
	GetAllBoards(userID string) ([]boardModel.Board, error)
	GetBoard(boardID, userID string) (*boardModel.Board, error)
	UpdateBoard(boardID, userID string, req boardModel.Request) error
	DeleteBoard(boardID, userID string) error
}

type Handler struct {
	proxy Proxy
}

func NewHandler(proxy Proxy) *Handler {
	return &Handler{proxy: proxy}
}

func (h *Handler) CreateBoardHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req boardModel.Request
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

		err := h.proxy.CreateBoard(userID, req)
		if err != nil {
			log.Printf("Failed create board: %v\n", err)
			h.handleError(ctx, err, "Failed to create board")
			return
		}

		ctx.Status(http.StatusCreated)
	}
}

func (h *Handler) GetAllBoardsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, ok := authctx.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		boards, err := h.proxy.GetAllBoards(userID)
		if err != nil {
			log.Printf("Failed to get boards: %v", err)
			h.handleError(ctx, err, "Failed to get boards")
			return
		}

		ctx.JSON(http.StatusOK, boards)
	}
}

func (h *Handler) GetBoardHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		userID, ok := authctx.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		board, err := h.proxy.GetBoard(id, userID)
		if err != nil {
			log.Printf("Failed to get board: %v", err)
			h.handleError(ctx, err, "Failed to get board")
			return
		}

		ctx.JSON(http.StatusOK, board)
	}
}

func (h *Handler) UpdateBoardHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req boardModel.Request
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"detail": "Invalid request body",
			})
			return
		}

		boardID := ctx.Param("id")

		userID, ok := authctx.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		if err := h.proxy.UpdateBoard(boardID, userID, req); err != nil {
			log.Printf("Failed to update board: %v", err)
			h.handleError(ctx, err, "Failed to update board")
			return
		}

		ctx.Status(http.StatusOK)
	}
}

func (h *Handler) DeleteBoardHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		userID, ok := authctx.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		if err := h.proxy.DeleteBoard(id, userID); err != nil {
			log.Printf("Failed to delete board: %v", err)
			h.handleError(ctx, err, "Failed to delete board")
			return
		}

		ctx.Status(http.StatusOK)
	}
}

func (h *Handler) handleError(ctx *gin.Context, err error, message string) {
	switch {
	case errors.Is(err, boardProxy.ErrForbidden):
		ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
			"detail": "Access denied",
		})
	case errors.Is(err, boardService.ErrBoardNotFound):
		ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
			"detail": "Board not found",
		})
	default:
		ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"detail": message,
		})
	}
}