package task

import (
	"database/sql"
	"errors"
	"kanban/internal/postgres"
	"kanban/internal/utils"
)

const maxTasks int = 52

var errForbidden error = errors.New("this is not your board")
var errTaskNotFound error = errors.New("task not found")
var errTaskLimitReached error = errors.New("task limit reached")
var errIncorrectPosition error = errors.New("task position is greater than possible or not positive")

func CreateTask(db *sql.DB, task Task, userID string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	ok, err := checkOwnershipByColumn(tx, task.ColumnID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return errForbidden
	}

	var count int
	err = tx.QueryRow(postgres.QueryGetTasksCount, task.ColumnID).Scan(&count)
	if err != nil {
		return err
	}
	if count >= maxTasks {
		return errTaskLimitReached
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

func GetAllTasks(db *sql.DB, columnID, userID string) ([]Task, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	ok, err := checkOwnershipByColumn(tx, columnID, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errForbidden
	}

	rows, err := tx.Query(postgres.QueryGetAllTasks, columnID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
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

func GetTask(db *sql.DB, taskID, userID string) (*Task, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	columnID, err := getColumnID(tx, taskID)
	if err != nil {
		return nil, err
	}

	ok, err := checkOwnershipByColumn(tx, columnID, userID)
	if err != nil {
		return nil, err
	}
	if !ok {
		return nil, errForbidden
	}

	var task Task
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
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errTaskNotFound
		}
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return &task, nil
}

func UpdateTaskContent(db *sql.DB, updatedTask updateTaskRequest, taskID, userID string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	columnID, err := getColumnID(tx, taskID)
	if err != nil {
		return err
	}

	ok, err := checkOwnershipByColumn(tx, columnID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return errForbidden
	}

	_, err = tx.Exec(
		postgres.QueryUpdateTaskContent,
		updatedTask.Name,
		updatedTask.Description,
		updatedTask.Done,
		updatedTask.Deadline,
		utils.GenerateTimestamp(),
		taskID,
	)
	if err != nil {
		return err
	}
	
	return tx.Commit()
}

func UpdateTaskColumn(db *sql.DB, updatedTask updateTaskRequest, taskID, userID string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	columnID, oldPos, err := getColumnIDAndPosition(tx, taskID)
	if err != nil {
		return err
	}

	ok, err := checkOwnershipByColumn(tx, columnID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return errForbidden
	}

	err = checkPosition(tx, columnID, *updatedTask.Position)
	if err != nil {
		return err
	}

	_, err = tx.Exec(postgres.QueryMoveTasksForInsert, *updatedTask.ColumnID, *updatedTask.Position)
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		postgres.QueryUpdateTaskColumn, 
		updatedTask.ColumnID,
		updatedTask.Position,
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

func UpdateTaskPosition(db *sql.DB, updatedTask updateTaskRequest, taskID, userID string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	columnID, oldPos, err := getColumnIDAndPosition(tx, taskID)
	if err != nil {
		return err
	}

	ok, err := checkOwnershipByColumn(tx, columnID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return errForbidden
	}

	err = checkPosition(tx, columnID, *updatedTask.Position + 1)
	if err != nil {
		return err
	}

	if *updatedTask.Position > oldPos {
		_, err = tx.Exec(postgres.QueryMoveTasksUp, columnID, oldPos, *updatedTask.Position)
	} else if *updatedTask.Position < oldPos {
		_, err = tx.Exec(postgres.QueryMoveTasksDown, columnID, *updatedTask.Position, oldPos)
	}
	if err != nil {
		return err
	}

	_, err = tx.Exec(
		postgres.QueryUpdateTaskPosition, 
		updatedTask.Position, 
		utils.GenerateTimestamp(), 
		taskID,
	)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func DeleteTask(db *sql.DB, taskID, userID string) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	columnID, pos, err := getColumnIDAndPosition(tx, taskID)
	if err != nil {
		return err
	}

	ok, err := checkOwnershipByColumn(tx, columnID, userID)
	if err != nil {
		return err
	}
	if !ok {
		return errForbidden
	}

	_, err = db.Exec(postgres.QueryMoveTaskForDelete, columnID, pos)
	if err != nil {
		return err
	}

	_, err = db.Exec(postgres.QueryDeleteTask, taskID)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func checkOwnershipByColumn(tx *sql.Tx, columnID, userID string) (bool, error) {
	var ok bool
	err := tx.QueryRow(postgres.QueryCheckBoardOwnershipForTask, columnID, userID).Scan(&ok)
	return ok, err
}

func getColumnID(tx *sql.Tx, taskID string) (string, error) {
	var columnID string
	err := tx.QueryRow(postgres.QueryGetColumnID, taskID).Scan(&columnID)
	if errors.Is(err, sql.ErrNoRows) {
		return "", errTaskNotFound
	}
	return columnID, err
}

func getColumnIDAndPosition(tx *sql.Tx, taskID string) (string, int, error) {
	var columnID string
	var pos int
	err := tx.QueryRow(postgres.QueryGetColumnIDAndPosition, taskID).Scan(&columnID, &pos)
	if errors.Is(err, sql.ErrNoRows) {
		return "", -1, errTaskNotFound
	}
	return columnID, pos, err
}

func checkPosition(tx *sql.Tx, columnID string, pos int) error {
	var maxPos int
	err := tx.QueryRow(postgres.QueryGetMaxTaskPosition, columnID).Scan(&maxPos)
	if pos > maxPos || pos <= 0 {
		return errIncorrectPosition
	}
	return err
}
