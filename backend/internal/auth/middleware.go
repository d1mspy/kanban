package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const userIDContextKey = "userID"

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "Authorization header missing",
			})
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "Invalid token format",
			})
			return
		}
		
		claims, err := ValidateJWT(parts[1])
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"detail": "Invalid or expire token",
			})
			return
		}

		ctx.Set(userIDContextKey, claims.UserID)

		ctx.Next()
	}
}

func GetUserID(ctx *gin.Context) (string, bool) {
	userID, ok := ctx.Get(userIDContextKey)
	if !ok {
		return "", false
	}
	return userID.(string), true
}