package taskProxy

import (
	"errors"
	"fmt"
	taskModel "kanban/internal/task/model"
)

var ErrForbidden = errors.New("access denied")

type Service interface {
	CreateTask(columnID string, req taskModel.CreateRequest) error
	GetAllTasks(columnID string) ([]taskModel.Task, error)
	GetTask(taskID string) (*taskModel.Task, error)
	UpdateTask(taskID string, req taskModel.UpdateRequest) error
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

func (p *Proxy) CreateTask(columnID, userID string, req taskModel.CreateRequest) error {
	isOwner, err := p.checkColumnOwnership(columnID, userID)
	if err != nil {
		return fmt.Errorf("taskProxy.CreateTask: %w", err)
	}
 
	if isOwner {
		return p.service.CreateTask(columnID, req)
	} else {
		return fmt.Errorf("taskProxy.CreateTask: %w", ErrForbidden)
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
		return nil, fmt.Errorf("taskProxy.GetAllTasks: %w", ErrForbidden)
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
		return nil, fmt.Errorf("taskProxy.GetTask: %w", ErrForbidden)
	}
}

func (p *Proxy) UpdateTask(taskID, userID string, req taskModel.UpdateRequest) error {
	isOwner, err := p.checkTaskOwnership(taskID, userID)
	if err != nil {
		return fmt.Errorf("taskProxy.UpdateTask: %w", err)
	}

	if isOwner {
		return p.service.UpdateTask(taskID, req)
	} else {
		return fmt.Errorf("taskProxy.UpdateTask: %w", ErrForbidden)
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
		return fmt.Errorf("taskProxy.DeleteTask: %w", ErrForbidden)
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