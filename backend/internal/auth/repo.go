package auth

import (
	"database/sql"
	"errors"
	"kanban/internal/postgres"
	"kanban/internal/utils"
)

var errUserNotFound = errors.New("User not found")

func CreateUser(db *sql.DB, user User) error {
	_, err := db.Exec(postgres.QueryCreateUser, user.ID, user.Username, user.Password, utils.GenerateTimestamp())

	return err
}

func GetUserByUsername(db *sql.DB, username string) (User, error) {
	var user User
	err := db.QueryRow(postgres.QueryGetUserByUsername, username).Scan(&user.ID, &user.Username, &user.Password)
	if err == sql.ErrNoRows {
		return user, errUserNotFound
	}
	return user, err
}