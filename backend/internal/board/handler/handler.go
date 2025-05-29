package boardHandler

import (
	"kanban/internal/auth"
	boardModel "kanban/internal/board/model"
	boardService "kanban/internal/board/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Service interface {
	CreateBoard(userID, name string) error
	GetAllBoards(userID string) ([]boardModel.Board, error)
	GetBoard(id, userID string) (boardModel.Board, error)
	UpdateBoard(id, name, userID string) error
	DeleteBoard(boardID, userID string) error
}

type Handler struct {
	service *boardService.Service
}

func NewHandler(serv *boardService.Service) *Handler {
	return &Handler{service: serv}
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

		err := h.service.CreateBoard(userID, req.Name)
		if err != nil {
			log.Printf("Failed create board: %w\n", err)
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

		boards, err := h.service.GetAllBoards(userID)
		if err != nil {
			log.Printf("Failed to get boards: %w", err)
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

		board, err := h.service.GetBoard(id, userID)
		if err != nil {
			log.Printf("Failed to get board: %w", err)
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

		if err := h.service.UpdateBoard(id, req.Name, userID); err != nil {
			log.Printf("Failed to update board: %w", err)
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

		if err := h.service.DeleteBoard(id, userID); err != nil {
			log.Printf("Failed to delete board: %w", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to delete board",
			})
			return
		}

		ctx.Status(http.StatusOK)
	}
}
