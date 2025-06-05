package taskProxy

import (
	"errors"
	"fmt"
	taskModel "kanban/internal/task/model"
)

var errForbidden = errors.New("this is not your board")

type Service interface {
	CreateTask(columnID, name, description string) error
	GetAllTasks(columnID string) ([]taskModel.Task, error)
	GetTask(taskID string) (*taskModel.Task, error)
	UpdateTask(req taskModel.UpdateTaskRequest, taskID string) error
	DeleteTask(taskID string) error
	GetUserByColumn(columnID string) (*string, error)
	GetUserByTask(taskID string) (*string, error)
}

type Proxy struct {
	service Service
}

func NewProxy(service Service) *Proxy {
	return &Proxy{service: service}
}

func (p *Proxy) CreateTask(columnID, name, description, userID string) error {
	isOwner, err := p.checkColumnOwnership(columnID, userID)
	if err != nil {
		return fmt.Errorf("taskProxy.CreateTask: %w", err)
	}

	if isOwner {
		return p.service.CreateTask(columnID, name, description)
	} else {
		return errForbidden
	}
}

func (p *Proxy) GetAllTasks(columnID, userID string) ([]taskModel.Task, error) {
	isOwner, err := p.checkColumnOwnership(columnID, userID)
	if err != nil {
		return nil, fmt.Errorf("taskProxy.GetAllTasks: %w", err)
	}

	if isOwner {
		return p.service.GetAllTasks(columnID)
	} else {
		return nil, errForbidden
	}
}

func (p *Proxy) GetTask(taskID, userID string) (*taskModel.Task, error) {
	isOwner, err := p.checkTaskOwnership(taskID, userID)
	if err != nil {
		return nil, fmt.Errorf("taskProxy.GetTask: %w", err)
	}

	if isOwner {
		return p.service.GetTask(taskID)
	} else {
		return nil, errForbidden
	}
}

func (p *Proxy) UpdateTask(req taskModel.UpdateTaskRequest, taskID, userID string) error {
	isOwner, err := p.checkTaskOwnership(taskID, userID)
	if err != nil {
		return fmt.Errorf("taskProxy.UpdateTask: %w", err)
	}

	if isOwner {
		return p.service.UpdateTask(req, taskID)
	} else {
		return errForbidden
	}
}

func (p *Proxy) DeleteTask(taskID, userID string) error {
	isOwner, err := p.checkTaskOwnership(taskID, userID)
	if err != nil {
		return fmt.Errorf("taskProxy.DeleteTask: %w", err)
	}

	if isOwner {
		return p.service.DeleteTask(taskID)
	} else {
		return errForbidden
	}
}

func (p *Proxy) checkColumnOwnership(columnID, userID string) (bool, error) {
	realUserID, err := p.service.GetUserByColumn(columnID)
	if err != nil {
		return false, fmt.Errorf("taskProxy.checkColumnOwnership: %w", err)
	}

	return *realUserID == userID, nil
}

func (p *Proxy) checkTaskOwnership(taskID, userID string) (bool, error) {
	realUserID, err := p.service.GetUserByTask(taskID)
	if err != nil {
		return false, fmt.Errorf("taskProxy.checkTaskOwnership: %w", err)
	}

	return *realUserID == userID, nil
}

