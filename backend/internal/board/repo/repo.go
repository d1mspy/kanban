package boardRepo

import (
	"database/sql"
	"errors"
	"fmt"
	boardModel "kanban/internal/board/model"
	"kanban/internal/postgres"
	"kanban/internal/utils"
)

var errBoardNotFound error = errors.New("board not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(board boardModel.Board) error {
	_, err := r.db.Exec(
		postgres.QueryCreateBoard,
		board.ID,
		board.UserID,
		utils.GenerateTimestamp(),
		utils.GenerateTimestamp(),
		board.Name,
	)
	if err != nil {
		return fmt.Errorf("boardRepo.Create: %w", err)
	}

	return nil
}

func (r *Repository) GetAll(userID string) ([]boardModel.Board, error) {
	rows, err := r.db.Query(postgres.QueryGetAllBoards, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var boards []boardModel.Board
	for rows.Next() {
		var board boardModel.Board
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

func (r *Repository) Get(boardID string) (*boardModel.Board, error) {
	var board boardModel.Board
	err := r.db.QueryRow(postgres.QueryGetBoard, boardID).Scan(
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

func (r *Repository) Update(board boardModel.Board) error {
	_, err := r.db.Exec(
		postgres.QueryUpdateBoard, 
		utils.GenerateTimestamp(), 
		board.Name, 
		board.ID, 
	)
	
	return err
}

func (r *Repository) Delete(boardID string) error {
	_, err := r.db.Exec(postgres.QueryDeleteBoard, boardID)
	return err
}

func (r *Repository) GetUserByBoard(boardID string) (*string, error) {
	var userID string
	err := r.db.QueryRow(postgres.QueryGetUserByBoardID, boardID).Scan(&userID)
	return &userID, err
}