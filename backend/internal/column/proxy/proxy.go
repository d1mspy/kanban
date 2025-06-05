package columnProxy

import (
	"errors"
	"fmt"
	columnModel "kanban/internal/column/model"
	columnService "kanban/internal/column/service"
)

var errForbidden = errors.New("this is not your board")

type Service interface {
	CreateColumn(boardID, name string) error
	GetAllColumns(boardID string)
	GetColumn(boardID string) (*columnModel.Column, error)
	UpdateColumn(columnID string, newName *string, newPos *int) error
	DeleteColumn(columnID string) error
	GetUserByBoard(boardID string) (*string, error)
	GetUserByColumn(columnID string) (*string, error)
}

type Proxy struct {
	service *columnService.Service
}

func NewProxy(service *columnService.Service) *Proxy {
	return &Proxy{service: service}
}

func (p *Proxy) CreateColumn(boardID, name, userID string) error {
	isOwner, err := p.checkBoardOwnership(boardID, userID)
	if err != nil {
		return fmt.Errorf("columnProxy.CreateColumn: %w", err)
	}

	if isOwner {
		return p.service.CreateColumn(boardID, name)
	} else {
		return errForbidden
	}
}

func (p *Proxy) GetAllColumns(boardID, userID string) ([]columnModel.Column, error) {
	isOwner, err := p.checkBoardOwnership(boardID, userID)
	if err != nil {
		return nil, fmt.Errorf("columnProxy.GetAllColumns: %w", err)
	}

	if isOwner {
		return p.service.GetAllColumns(boardID)
	} else {
		return nil, errForbidden
	}
}

func (p *Proxy) GetColumn(columnID, userID string) (*columnModel.Column, error) {
	isOwner, err := p.checkColumnOwnership(columnID, userID)
	if err != nil {
		return nil, fmt.Errorf("columnProxy.GetColumn: %w", err)
	}

	if isOwner {
		return p.service.GetColumn(columnID)
	} else {
		return nil, errForbidden
	}
}

func (p *Proxy) UpdateColumn(columnID, userID string, newName *string, newPos *int) error {
	isOwner, err := p.checkColumnOwnership(columnID, userID)
	if err != nil {
		return fmt.Errorf("columnProxy.UpdateColumn: %w", err)
	}

	if isOwner {
		return p.service.UpdateColumn(columnID, newName, newPos)
	} else {
		return errForbidden
	}
}

func (p *Proxy) DeleteColumn(columnID, userID string) error {
	isOwner, err := p.checkColumnOwnership(columnID, userID)
	if err != nil {
		return fmt.Errorf("columnProxy.UpdateColumn: %w", err)
	}

	if isOwner {
		return p.service.DeleteColumn(columnID)
	} else {
		return errForbidden
	}
}

func (p *Proxy) checkBoardOwnership(boardID, userID string) (bool, error) {
	realUserID, err := p.service.GetUserByBoard(boardID)
	if err != nil {
		return false, fmt.Errorf("columnProxy.checkBoardOwnership: %w", err)
	}

	return *realUserID == userID, nil
}

func (p *Proxy) checkColumnOwnership(columnID, userID string) (bool, error) {
	realUserID, err := p.service.GetUserByColumn(columnID)
	if err != nil {
		return false, fmt.Errorf("columnProxy.checkColumnOwnership: %w", err)
	}

	return *realUserID == userID, nil
}
