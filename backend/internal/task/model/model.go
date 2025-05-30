package taskModel

import "time"

type Task struct {
	ID          string     `json:"id"`
	ColumnID    string     `json:"column_id"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Name        string     `json:"name"`
	Description string     `json:"description"`
	Position    int        `json:"position"`
	Done        bool       `json:"done"`
	Deadline    *time.Time `json:"deadline"`
}

type CreateTaskRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

type UpdateTaskRequest struct {
	ColumnID    *string     `json:"column_id"`
	Name        *string     `json:"name"`
	Description *string     `json:"description"`
	Position    *int        `json:"position"`
	Done        *bool       `json:"done"`
	Deadline    *time.Time  `json:"deadline"`
}