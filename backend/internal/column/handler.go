package column

import (
	"database/sql"
	"kanban/internal/auth"
	"kanban/internal/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type columnRequest struct {
	Name     *string `json:"name"`
	Position *int    `json:"position"`
}

func CreateColumnHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req columnRequest
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
		column := Column{
			ID: utils.NewUUID(),
			BoardID: boardID,
			Name: *req.Name,
		}

		if err := CreateColumn(db, column, userID); err != nil {
			switch err {
			case errForbidden:
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"detail": "This is not your column",
				})
				return
			case errColumnLimitReached:
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"detail": "Column limit reached",
				})
				return
			default:
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"detail": "Failed to update column",
				})
				return
			}
		}

		ctx.Status(http.StatusCreated)
	}
}

func GetAllColumnsHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		boardID := ctx.Param("id")

		userID, ok := auth.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		columns, err := GetAllColumns(db, boardID, userID)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to get all columns",
			})
			return
		}

		ctx.JSON(http.StatusOK, columns)
	}
}

func GetColumnHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		userID, ok := auth.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		column, err := GetColumn(db, id, userID)
		if err != nil {
			if err == errColumnNotFound {
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"detail": "Column not found",
				})
				return
			}
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to get column",
			})
			return
		}

		ctx.JSON(http.StatusOK, column)
	}
}

func UpdateColumnHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req columnRequest
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

		if err := UpdateColumn(db, userID, id, req.Name, req.Position); err != nil {
			switch err {
			case errColumnNotFound:
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"detail": "Column not found",
				})
				return
			case errForbidden:
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"detail": "This is not your column",
				})
				return
			case errIncorrectPosition:
				ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
					"detail": "Position is greater than possible or not positive",
				})
				return
			default:
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"detail": "Failed to update column",
				})
				return
			}

		}

		ctx.Status(http.StatusOK)
	}
}

func DeleteColumnHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id := ctx.Param("id")

		userID, ok := auth.GetUserID(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "No token",
			})
			return
		}

		if err := DeleteColumn(db, userID, id); err != nil {
			switch err {
			case errColumnNotFound:
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"detail": "Column not found",
				})
				return
			case errForbidden:
				ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{
					"detail": "This is not your column",
				})
				return
			default:
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"detail": "Failed to delete column",
				})
				return
			}
		}

		ctx.Status(http.StatusOK)
	}
}