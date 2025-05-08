package auth

import (
	"database/sql"
	"kanban/internal/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

type authRequest struct {
	Username string `json:"username"`
    Password string `json:"password"`
}

func RegisterHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req authRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"detail": "Invalid request body",
			})
			return
		}

		hash, err := HashPassword(req.Password)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to hash the password",
			})
			return
		}

		user := User{
			ID:       utils.NewUUID(),
			Username: req.Username,
			Password: hash,
		}

		if err := CreateUser(db, user); err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to create user",
			})
			return
		}

		token, err := GenerateJWT(user.ID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to generate token",
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"token": token})
	}
}

func LoginHandler(db *sql.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req authRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"detail": "Invalid JSON body",
			})
			return
		}

		user, err := GetUserByUsername(db, req.Username)
		if err != nil {
			if err == errUserNotFound {
				ctx.AbortWithStatusJSON(http.StatusTeapot, gin.H{
					"detail": "Check your name, bro. Do you exist?",
				})
				return
			} else {
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"detail": "Failed to find the user",
				})
				return
			}
		}

		if err = CheckPasswordHash(req.Password, user.Password); err != nil {
			ctx.AbortWithStatusJSON(http.StatusTeapot, gin.H{
				"detail": "Check your password, bro",
			})
			return
		}

		token, err := GenerateJWT(user.ID)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to generate token",
			})
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"token": token})
	}
}