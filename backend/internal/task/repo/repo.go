package taskRepo

import (
	"database/sql"
	"errors"
	"kanban/internal/postgres"
	taskModel "kanban/internal/task/model"
	"kanban/internal/utils"
)

const maxTasks int = 52

var ErrTaskLimitReached error = errors.New("task limit reached")
var ErrIncorrectPosition error = errors.New("task position is greater than possible or not positive")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(task taskModel.Task) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	var count int
	err = tx.QueryRow(postgres.QueryGetTasksCount, task.ColumnID).Scan(&count)
	if err != nil {
		return err
	}
	if count >= maxTasks {
		return ErrTaskLimitReached
	}

	err = tx.QueryRow(postgres.QueryGetMaxTaskPosition, task.ColumnID).Scan(&task.Position)
	if err != nil{
		return err
	}

	_, err = tx.Exec(
		postgres.QueryCreateTask,
		utils.NewUUID(),
		task.ColumnID,
		utils.GenerateTimestamp(),
		utils.GenerateTimestamp(),
		task.Name,
		task.Description,
		task.Position,
		false,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) GetAll(columnID string) ([]taskModel.Task, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	rows, err := tx.Query(postgres.QueryGetAllTasks, columnID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []taskModel.Task
	for rows.Next() {
		var task taskModel.Task
		if err = rows.Scan(
			&task.ID,
			&task.ColumnID,
			&task.CreatedAt,
			&task.UpdatedAt,
			&task.Name,
			&task.Description,
			&task.Position,
			&task.Done,
			&task.Deadline,
		); err != nil {
			return nil, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	if err = tx.Commit(); err != nil {
		return nil, err
	}

	return tasks, nil
}

func (r *Repository) Get(taskID string) (*taskModel.Task, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	var task taskModel.Task
	err = tx.QueryRow(postgres.QueryGetTask, taskID).Scan(
		&task.ID,
		&task.ColumnID,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.Name,
		&task.Description,
		&task.Position,
		&task.Done,
		&task.Deadline,
	)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func (r *Repository) UpdateContent(taskID string, req taskModel.UpdateRequest) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec(
		postgres.QueryUpdateTaskContent,
		req.Name,
		req.Description,
		req.Done,
		req.Deadline,
		utils.GenerateTimestamp(),
		taskID,
	)
	if err != nil {
		return err
	}
	
	return tx.Commit()
}

func (r *Repository) UpdateColumn(taskID string, req taskModel.UpdateRequest) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	columnID, oldPos, err := getColumnIDAndPosition(tx, taskID)
	if err != nil {
		return err
	}

	err = checkPosition(tx, columnID, *req.Position)
	if err != nil {
		return err
	}

	_, err = tx.Exec(postgres.QueryMoveTasksForInsert, *req.ColumnID, *req.Position)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		postgres.QueryUpdateTaskColumn, 
		req.ColumnID,
		req.Position,
		utils.GenerateTimestamp(),
		taskID,
	)
	if err != nil {
		return err
	}

	_, err = tx.Exec(postgres.QueryMoveTaskForDelete, columnID, oldPos)
	if err != nil {
		return err
	}


	return tx.Commit()
}

func (r *Repository) UpdatePosition(taskID string, req taskModel.UpdateRequest) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	columnID, oldPos, err := getColumnIDAndPosition(tx, taskID)
	if err != nil {
		return err
	}

	err = checkPosition(tx, columnID, *req.Position + 1)
	if err != nil {
		return err
	}

	if *req.Position > oldPos {
		_, err = tx.Exec(postgres.QueryMoveTasksUp, columnID, oldPos, *req.Position)
	} else if *req.Position < oldPos {
		_, err = tx.Exec(postgres.QueryMoveTasksDown, columnID, *req.Position, oldPos)
	}
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		postgres.QueryUpdateTaskPosition, 
		req.Position, 
		utils.GenerateTimestamp(), 
		taskID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) Delete(taskID string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	columnID, pos, err := getColumnIDAndPosition(tx, taskID)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(postgres.QueryMoveTaskForDelete, columnID, pos)
	if err != nil {
		return err
	}

	_, err = r.db.Exec(postgres.QueryDeleteTask, taskID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (r *Repository) GetUserByColumn(columnID string) (*string, error) {
	var userID string
	err := r.db.QueryRow(postgres.QueryGetUserByColumnID, columnID).Scan(&userID)
	return &userID, err
}

func (r *Repository) GetUserByTask(taskID string) (*string, error) {
	var userID string
	err := r.db.QueryRow(postgres.QueryGetUserByTaskID, taskID).Scan(&userID)
	return &userID, err
}

func getColumnIDAndPosition(tx *sql.Tx, taskID string) (string, int, error) {
	var columnID string
	var pos int
	err := tx.QueryRow(postgres.QueryGetColumnIDAndPosition, taskID).Scan(&columnID, &pos)
	if err != nil {
		return "", -1, err
	}
	return columnID, pos, err
}

func checkPosition(tx *sql.Tx, columnID string, pos int) error {
	var maxPos int
	err := tx.QueryRow(postgres.QueryGetMaxTaskPosition, columnID).Scan(&maxPos)
	if pos > maxPos || pos <= 0 {
		return ErrIncorrectPosition
	}
	return err
}
