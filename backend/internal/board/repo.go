package board

import (
	"database/sql"
	"errors"
	"kanban/internal/postgres"
	"kanban/internal/utils"
)

var errBoardNotFound error = errors.New("Board not found")

func CreateBoard(db *sql.DB, board Board) error {
	_, err := db.Exec(
		postgres.QueryCreateBoard, 
		board.ID,
		board.UserID,
		utils.GenerateTimestamp(),
		utils.GenerateTimestamp(),
		board.Name,
	)

	return err
}

func GetBoard(db *sql.DB, boardID, userID string) (*Board, error) {
	var board Board
	err := db.QueryRow(postgres.QueryGetBoard, boardID, userID).Scan(
		&board.ID, 
		&board.UserID,
		&board.CreatedAt,
		&board.UpdatedAt,
		&board.Name,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errBoardNotFound
		}
		return nil, err
	}

	return &board, nil
}

func UpdateBoard(db *sql.DB, board Board) error {
	_, err := db.Exec(postgres.QueryUpdateBoard, utils.GenerateTimestamp(), board.Name, board.ID, board.UserID)
	return err
}

func DeleteBoard(db *sql.DB, boardID, userID string) error {
	_, err := db.Exec(postgres.QueryDeleteBoard, boardID, userID)
	return err
}

func GetAllBoards(db *sql.DB, userID string) ([]Board, error) {
	rows, err := db.Query(postgres.QueryGetAllBoards, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var boards []Board
	for rows.Next() {
		var board Board
		if err := rows.Scan(
			&board.ID,
			&board.UserID,
			&board.CreatedAt,
			&board.UpdatedAt,
			&board.Name,
		); err != nil {
			return nil, err
		}
		boards = append(boards, board)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return boards, nil
}