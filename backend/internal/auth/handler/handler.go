package authHandler

import (
	"errors"
	authModel "kanban/internal/auth/model"
	authService "kanban/internal/auth/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Service interface {
	CreateUser(req authModel.Request) (*string, error)
	LoginUser(req authModel.Request) (*string, error)
}

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req authModel.Request
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"detail": "Invalid request body",
			})
			return
		}
		
		token, err := h.service.CreateUser(req)
		if err != nil {
			log.Printf("Failed to create user: %v\n", err)
			ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
				"detail": "Failed to create user",
			})
			return
		}

		ctx.JSON(http.StatusCreated, gin.H{"token": token})
	}
}

func (h *Handler) LoginHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req authModel.Request
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{
				"detail": "Invalid JSON body",
			})
			return
		}

		token, err := h.service.LoginUser(req)
		if err != nil {
			log.Printf("Failed to login: %v", err)
			switch {
			case errors.Is(err, authService.ErrUserNotFound):
				ctx.AbortWithStatusJSON(http.StatusNotFound, gin.H{
					"detail": "User not found",
				})
				return
			case errors.Is(err, authService.ErrIncorrectPassword):
				ctx.AbortWithStatusJSON(http.StatusTeapot, gin.H{
					"detail": "Check your password, bro",
				})
				return
			default:
				ctx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
					"detail": "Failed to login",
				})
				return
			}
		}

		ctx.JSON(http.StatusOK, gin.H{"token": token})
	}
}