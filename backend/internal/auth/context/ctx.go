package authctx

import "github.com/gin-gonic/gin"

const userIDContextKey string = "userID"

func SetUserID(ctx *gin.Context, userID string) {
	ctx.Set(userIDContextKey, userID)
}


func GetUserID(ctx *gin.Context) (string, bool) {
	userID, ok := ctx.Get(userIDContextKey)
	if !ok {
		return "", false
	}
	return userID.(string), true
}