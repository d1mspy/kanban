package taskService

import (
	"database/sql"
	"errors"
	"fmt"
	taskModel "kanban/internal/task/model"
	"reflect"
)

var ErrTaskNotFound error = errors.New("task not found")
var ErrBadUpdateRequest error = errors.New("invalid combination of fields")

type updateCase string
const (
	caseContent   updateCase = "content"
	caseColumn	  updateCase = "column"
	casePosition  updateCase = "position"
	caseUndefined updateCase = "undefined"
)

type Repository interface {
	Create(task taskModel.Task) error
	GetAll(columnID string) ([]taskModel.Task, error)
	Get(taskID string) (*taskModel.Task, error)
	UpdateContent(taskID string, req taskModel.UpdateRequest) error
	UpdateColumn(taskID string, req taskModel.UpdateRequest) error
	UpdatePosition(taskID string, req taskModel.UpdateRequest) error
	Delete(taskID string) error
	GetUserByColumn(columnID string) (*string, error)
	GetUserByTask(taskID string) (*string, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTask(columnID string, req taskModel.CreateRequest) error {
	task := taskModel.Task{
		ColumnID: columnID,
		Name: req.Name,
	}

	err := s.repo.Create(task)
	if err != nil {
		return fmt.Errorf("taskService.CreateTask: %w", err)
	}

	return nil
}

func (s *Service) GetAllTasks(columnID string) ([]taskModel.Task, error) {
	tasks, err := s.repo.GetAll(columnID)
	if err != nil {
		return nil, fmt.Errorf("taskService.GetAllTasks: %w", err)
	}

	return tasks, nil
}

func (s *Service) GetTask(taskID string) (*taskModel.Task, error) {
	task, err := s.repo.Get(taskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("taskService.GetTask: %w", ErrTaskNotFound)
		}
		return nil, fmt.Errorf("taskService.GetTask: %w", err)
	}

	return task, nil
}

func (s *Service) UpdateTask(taskID string, req taskModel.UpdateRequest) error {
	updCase := validateUpdateTaskRequest(req)

	var err error
	switch updCase {
	case caseContent:
		err = s.repo.UpdateContent(taskID, req)
	case caseColumn:
		err = s.repo.UpdateColumn(taskID, req)
	case casePosition:
		err = s.repo.UpdatePosition(taskID, req)
	default:
		return fmt.Errorf("taskService.UpdateTask: %w", ErrBadUpdateRequest)
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("taskService.UpdateTask: %w", ErrTaskNotFound)
		}
		return fmt.Errorf("taskService.UpdateTask: %w", err)
	}

	return nil
}

func (s *Service) DeleteTask(taskID string) error {
	err := s.repo.Delete(taskID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("taskService.DeleteTask: %w", ErrTaskNotFound)
		}
		return fmt.Errorf("taskService.DeleteTask: %w", err)
	}

	return nil
}


func (s *Service) GetUserByColumn(columnID string) (*string, error) {
	userID, err := s.repo.GetUserByColumn(columnID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("taskService.GetUserByColumn: %w", ErrTaskNotFound)
		}
		return nil, fmt.Errorf("taskService.GetUserByColumn: %w", err)
	}

	return userID, nil
}

func (s *Service) GetUserByTask(taskID string) (*string, error) {
	userID, err := s.repo.GetUserByTask(taskID)
		if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("taskService.GetUserByTask: %w", ErrTaskNotFound)
		}
		return nil, fmt.Errorf("taskService.GetUserByTask: %w", err)
	}

	return userID, nil
}

func validateUpdateTaskRequest(req taskModel.UpdateRequest) updateCase {
	contentFields := []any{
		req.Name,
		req.Description,
		req.Done,
		req.Deadline,
	}

	contentCount := 0
	for _, field := range contentFields {
		if !isNil(field) {
			contentCount++
		}
	}

	columnSet := req.ColumnID != nil
	positionSet := req.Position != nil

	switch {
	case contentCount == 1 && !columnSet && !positionSet:
		return caseContent
	case contentCount == 0 && columnSet && positionSet:
		return caseColumn
	case contentCount == 0 && !columnSet && positionSet:
		return casePosition
	default:
		return caseUndefined
	}
}

func isNil(v any) bool {
	return v == nil  || reflect.ValueOf(v).IsNil()
}