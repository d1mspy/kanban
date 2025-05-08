package board

import (
	"database/sql"
	"kanban/internal/auth"
	"kanban/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type boardRequest struct {
	Name string `json:"name"`
}

func CreateBoardHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req boardRequest
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

		board := Board{
			ID: utils.NewUUID(),
			UserID: userID,
			Name: req.Name,
		}

		if err := CreateBoard(db, board); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to create board",
			})
		}

		ctx.Status(http.StatusCreated)
	}
}

func GetBoardHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		userID, ok := auth.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		board, err := GetBoard(db, id, userID)
		if err != nil {
			if err == errBoardNotFound {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"detail": "Board not found",
				})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to get board",
			})
			return
		}

		ctx.JSON(http.StatusOK, board)
	}
}

func UpdateBoardHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req boardRequest
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

		board := Board{
			ID: id,
			Name: req.Name,
			UserID: userID,
		}

		if err := UpdateBoard(db, board); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to update board",
			})
			return
		}

		ctx.Status(http.StatusOK)
	}
}

func DeleteBoardHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		userID, ok := auth.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		if err := DeleteBoard(db, id, userID); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to delete board",
			})
			return
		}

		ctx.Status(http.StatusOK)
	}
}

func GetAllBoardsHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, ok := auth.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		boards, err := GetAllBoards(db, userID); 
		
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to load all board data",
			})
			return
		}

		ctx.JSON(http.StatusOK, boards)
	}
}