package columnRepo

import (
	"database/sql"
	"errors"
	columnModel "kanban/internal/column/model"
	"kanban/internal/postgres"
	"kanban/internal/utils"
)

const maxColumns int = 42

var errForbidden = errors.New("this is not your board")
var errColumnNotFound = errors.New("column not found")
var errColumnLimitReached = errors.New("column limit reached")
var errIncorrectPosition = errors.New("column position is greater than possible or not positive")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(column columnModel.Column, userID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var exists bool
	err = tx.QueryRow(postgres.QueryCheckBoardOwnershipForColumn, column.BoardID, userID).Scan(&exists)
	if err != nil {
		return err
	}
	if !exists {
		return errForbidden
	}

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

func (r *Repository) GetAll(boardID, userID string) ([]columnModel.Column, error) {
	rows, err := r.db.Query(postgres.QueryGetAllColumns, boardID, userID)
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

func (r *Repository) Get(id, userID string) (*columnModel.Column, error) {
	var column columnModel.Column
	err := r.db.QueryRow(postgres.QueryGetColumn, id, userID).Scan(
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

func (r *Repository) Update(userID, columnID string, newName *string, newPos *int) error {
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

	var ok bool
	err = tx.QueryRow(postgres.QueryCheckBoardOwnershipForColumn, boardID, userID).Scan(&ok)
	if err != nil {
		return err
	}
	if !ok {
		return errForbidden
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

func (r *Repository) Delete(userID, columnID string) error {
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

	var ok bool
	err = tx.QueryRow(postgres.QueryCheckBoardOwnershipForColumn, boardID, userID).Scan(&ok)
	if err != nil {
		return err
	}
	if !ok {
		return errForbidden
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