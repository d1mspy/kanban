package taskService

import (
	"fmt"
	taskModel "kanban/internal/task/model"
	"reflect"
)

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
	UpdateContent(updatedTask taskModel.UpdateTaskRequest, taskID string) error
	UpdateColumn(updatedTask taskModel.UpdateTaskRequest, taskID string) error
	UpdatePosition(updatedTask taskModel.UpdateTaskRequest, taskID string) error
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

func (s *Service) CreateTask(columnID, name, description string) error {
	task := taskModel.Task{
		ColumnID: columnID,
		Name: name,
		Description: description,
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
		return nil, fmt.Errorf("taskService.GetTask: %w", err)
	}

	return task, nil
}

func (s *Service) UpdateTask(req taskModel.UpdateTaskRequest, taskID string) error {
	updCase := validateUpdateTaskRequest(req)

	var err error
	switch updCase {
	case caseContent:
		err = s.repo.UpdateContent(req, taskID)
	case caseColumn:
		err = s.repo.UpdateColumn(req, taskID)
	case casePosition:
		err = s.repo.UpdatePosition(req, taskID)
	default:
		return fmt.Errorf("taskService.UpdateTask: %w", err)
	}

	if err != nil {
		return fmt.Errorf("taskService.UpdateTask: %w", err)
	}

	return nil
}

func (s *Service) DeleteTask(taskID string) error {
	err := s.repo.Delete(taskID)
	if err != nil {
		return fmt.Errorf("taskService.DeleteTask: %w", err)
	}

	return nil
}


func (s *Service) GetUserByColumn(columnID string) (*string, error) {
	return s.repo.GetUserByColumn(columnID)
}

func (s *Service) GetUserByTask(taskID string) (*string, error) {
	return s.repo.GetUserByTask(taskID)
}

func validateUpdateTaskRequest(req taskModel.UpdateTaskRequest) updateCase {
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