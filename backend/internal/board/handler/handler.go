package boardHandler

import (
	"kanban/internal/auth"
	boardModel "kanban/internal/board/model"
	boardProxy "kanban/internal/board/proxy"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Proxy interface {
	CreateBoard(userID, name string) error
	GetAllBoards(userID string) ([]boardModel.Board, error)
	GetBoard(boardID, userID string) (*boardModel.Board, error)
	UpdateBoard(boardID, name, userID string) error
	DeleteBoard(boardID, userID string) error
}

type Handler struct {
	proxy *boardProxy.Proxy
}

func NewHandler(proxy *boardProxy.Proxy) *Handler {
	return &Handler{proxy: proxy}
}

func (h *Handler) CreateBoardHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req boardModel.BoardRequest
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

		err := h.proxy.CreateBoard(userID, req.Name)
		if err != nil {
			log.Printf("Failed create board: %v\n", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to create board",
			})
			return
		}

		ctx.Status(http.StatusCreated)
	}
}

func (h *Handler) GetAllBoardsHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, ok := auth.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		boards, err := h.proxy.GetAllBoards(userID)
		if err != nil {
			log.Printf("Failed to get boards: %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to get boards",
			})
			return
		}

		ctx.JSON(http.StatusOK, boards)
	}
}

func (h *Handler) GetBoardHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		userID, ok := auth.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		board, err := h.proxy.GetBoard(id, userID)
		if err != nil {
			log.Printf("Failed to get board: %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to get board",
			})
			return
		}

		ctx.JSON(http.StatusOK, board)
	}
}

func (h *Handler) UpdateBoardHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req boardModel.BoardRequest
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

		if err := h.proxy.UpdateBoard(id, req.Name, userID); err != nil {
			log.Printf("Failed to update board: %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to update board",
			})
			return
		}

		ctx.Status(http.StatusOK)
	}
}

func (h *Handler) DeleteBoardHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		userID, ok := auth.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		if err := h.proxy.DeleteBoard(id, userID); err != nil {
			log.Printf("Failed to delete board: %v", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to delete board",
			})
			return
		}

		ctx.Status(http.StatusOK)
	}
}
