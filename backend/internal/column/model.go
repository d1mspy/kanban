package column

import "time"

type Column struct {
	ID        string    `json:"id"`
	BoardID   string    `json:"board_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	Position  int       `json:"position"`
}