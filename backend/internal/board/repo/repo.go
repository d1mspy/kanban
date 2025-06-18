package boardRepo

import (
	"database/sql"
	"fmt"
	boardModel "kanban/internal/board/model"
	"kanban/internal/postgres"
	"kanban/internal/utils"
)

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
		return nil, fmt.Errorf("boardRepo.GetAll: %w", err)
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
			return nil, fmt.Errorf("boardRepo.GetAll: %w", err)
		}
		boards = append(boards, board)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("boardRepo.GetAll: %w", err)
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
		return nil, fmt.Errorf("boardRepo.Get: %w", err)
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
	
	if err != nil {
		return fmt.Errorf("boardRepo.Update: %w", err)
	}
	return nil
}

func (r *Repository) Delete(boardID string) error {
	_, err := r.db.Exec(postgres.QueryDeleteBoard, boardID)
	if err != nil {
		return fmt.Errorf("boardRepo.Delete: %w", err)
	}
	return nil
}

func (r *Repository) GetUserByBoard(boardID string) (*string, error) {
	var userID string
	err := r.db.QueryRow(postgres.QueryGetUserByBoardID, boardID).Scan(&userID)
	if err != nil {
		return nil, fmt.Errorf("boardRepo.GetUserByBoard: %w", err)
	}
	return &userID, nil
}