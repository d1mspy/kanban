package authRepo

import (
	"database/sql"
	"fmt"
	authModel "kanban/internal/auth/model"
	"kanban/internal/postgres"
	"kanban/internal/utils"
)

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(user authModel.User) error {
	_, err := r.db.Exec(
		postgres.QueryCreateUser, 
		user.ID, user.Username, 
		user.Password, 
		utils.GenerateTimestamp(),
	)
	if err != nil {
		return fmt.Errorf("authRepo.Create: %w", err)
	}
	return nil
}

func (r *Repository) GetByUsername(username string) (*authModel.User, error) {
	var user authModel.User
	err := r.db.QueryRow(
		postgres.QueryGetUserByUsername, 
		username,
		).Scan(
			&user.ID, 
			&user.Username, 
			&user.Password,
		)
	if err != nil {
		return nil, fmt.Errorf("authRepo.GetByUsername: %w", err)
	}
	return &user, nil
}