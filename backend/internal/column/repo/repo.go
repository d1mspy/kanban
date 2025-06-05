package columnRepo

import (
	"database/sql"
	"errors"
	columnModel "kanban/internal/column/model"
	"kanban/internal/postgres"
	"kanban/internal/utils"
)

const maxColumns int = 42

var errColumnNotFound = errors.New("column not found")
var errColumnLimitReached = errors.New("column limit reached")
var errIncorrectPosition = errors.New("column position is greater than possible or not positive")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(column columnModel.Column) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var count int
	err = tx.QueryRow(postgres.QueryGetColumnsCount, column.BoardID).Scan(&count)
	if err != nil {
		return err
	}
	if count >= maxColumns {
		return errColumnLimitReached
	}

	err = tx.QueryRow(postgres.QueryGetMaxColumnPosition, column.BoardID).Scan(&column.Position)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		postgres.QueryCreateColumn,
		column.ID,
		column.BoardID,
		utils.GenerateTimestamp(),
		utils.GenerateTimestamp(),
		column.Name,
		column.Position,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) GetAll(boardID string) ([]columnModel.Column, error) {
	rows, err := r.db.Query(postgres.QueryGetAllColumns, boardID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []columnModel.Column
	for rows.Next() {
		var column columnModel.Column
		if err := rows.Scan(
			&column.ID,
			&column.BoardID,
			&column.CreatedAt,
			&column.UpdatedAt,
			&column.Name,
			&column.Position,
		); err != nil {
			return nil, err
		}
		columns = append(columns, column)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return columns, nil
}

func (r *Repository) Get(columnID string) (*columnModel.Column, error) {
	var column columnModel.Column
	err := r.db.QueryRow(postgres.QueryGetColumn, columnID).Scan(
		&column.ID,
		&column.BoardID,
		&column.CreatedAt,
		&column.UpdatedAt,
		&column.Name,
		&column.Position,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errColumnNotFound
		}
		return nil, err
	}

	return &column, nil
}

func (r *Repository) Update(columnID string, newName *string, newPos *int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var boardID string
	var oldPos int
	err = tx.QueryRow(postgres.QueryGetBoardIDAndColumnPos, columnID).Scan(&boardID, &oldPos)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errColumnNotFound
		}
		return err
	}

	var maxPos int
	if newPos != nil {
		err = tx.QueryRow(postgres.QueryGetMaxColumnPosition, boardID).Scan(&maxPos)
		if err != nil {
			return err
		}
		if *newPos >= maxPos || *newPos <= 0 {
			return errIncorrectPosition
		}
	}
	
	if newPos != nil && *newPos != oldPos {
		if *newPos > oldPos {
			_, err = tx.Exec(postgres.QueryMoveColumnsLeft, boardID, oldPos, *newPos)
		} else {
			_, err = tx.Exec(postgres.QueryMoveColumnsRight, boardID, *newPos, oldPos)
		}
		if err != nil {
			return err
		}
	}

	_, err = tx.Exec(postgres.QueryUpdateColumn, newName, newPos, utils.GenerateTimestamp(), columnID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) Delete(columnID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var boardID string
	var pos int
	err = tx.QueryRow(postgres.QueryGetBoardIDAndColumnPos, columnID).Scan(&boardID, &pos)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return errColumnNotFound
		}
		return err
	}

	_, err = tx.Exec(postgres.QueryDeleteColumn, columnID)
	if err != nil {
		return err
	}

	_, err = tx.Exec(postgres.QueryDecreaseColumnsPosition, boardID, pos)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) GetUserByBoard(boardID string) (*string, error) {
	var userID string
	err := r.db.QueryRow(postgres.QueryGetUserByBoardID, boardID).Scan(&userID)
	return &userID, err
}

func (r *Repository) GetUserByColumn(columnID string) (*string, error) {
	var userID string
	err := r.db.QueryRow(postgres.QueryGetUserByColumnID, columnID).Scan(&userID)
	return &userID, err
}
