package authModel

import "github.com/golang-jwt/jwt/v5"

type User struct {
	ID       string `json:"id"`
	Email    string `json:"email"`
	Username string	`json:"username"`
	Password string `json:"password"`
}

type RegisterRequest struct {
	Email    string `json:"email"    binding:"required,email"`
	Username string	`json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type LoginRequest struct {
	Email    string `json:"email"    binding:"required,email"`
    Password string `json:"password" binding:"required"`
}

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}