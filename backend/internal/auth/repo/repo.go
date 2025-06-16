package authRepo

import (
	"database/sql"
	"errors"
	authModel "kanban/internal/auth/model"
	"kanban/internal/postgres"
	"kanban/internal/utils"
)

var ErrUserNotFound = errors.New("user not found")

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

	return err
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

	if err == sql.ErrNoRows {
		return nil, ErrUserNotFound
	}
	return &user, err
}