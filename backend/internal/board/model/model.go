package boardModel

import "time"

type Board struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
}

type Request struct {
	Name string `json:"name" binding:"required"`
}
