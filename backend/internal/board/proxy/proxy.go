package boardProxy

import (
	"fmt"
	boardModel "kanban/internal/board/model"
	boardService "kanban/internal/board/service"
)

type Service interface {
	CreateBoard(userID, name string) error
	GetAllBoards(userID string) ([]boardModel.Board, error)
	GetBoard(boardID string) (boardModel.Board, error)
	UpdateBoard(boardID, name string) error
	DeleteBoard(boardID string) error
	GetUserByBoard(boardID string) (*string, error)
}

type Proxy struct {
	service *boardService.Service
}

func NewProxy(service *boardService.Service) *Proxy {
	return &Proxy{service: service}
}

func (p *Proxy) CreateBoard(userID, name string) error {
	return p.service.CreateBoard(userID, name)
}

func (p *Proxy) GetAllBoards(userID string) ([]boardModel.Board, error) {
	return p.service.GetAllBoards(userID)
}

func (p *Proxy) GetBoard(boardID, userID string) (*boardModel.Board, error) {
	isOwner, err := p.checkBoardOwnership(boardID, userID)
	if err != nil {
		return nil, fmt.Errorf("boardProxy.GetBoard: %w", err)
	}

	if isOwner {
		return p.service.GetBoard(boardID)
	} else {
		return nil, fmt.Errorf("boardProxy.GetBoard: this is not your board")
	}
}

func (p *Proxy) UpdateBoard(boardID, name, userID string) error {
	isOwner, err := p.checkBoardOwnership(boardID, userID)
	if err != nil {
		return err
	}

	if isOwner {
		return p.service.UpdateBoard(boardID, name)
	} else {
		return fmt.Errorf("this is not your board")
	}
}

func (p *Proxy) DeleteBoard(boardID, userID string) error {
	isOwner, err := p.checkBoardOwnership(boardID, userID)
	if err != nil {
		return err
	}

	if isOwner {
		return p.service.DeleteBoard(boardID)
	} else {
		return fmt.Errorf("this is not your board")
	}
}

func (p *Proxy) checkBoardOwnership(boardID, userID string) (bool, error) {
	realUserID, err := p.service.GetUserByBoard(boardID)
	if err != nil {
		return false, fmt.Errorf("boardProxy.checkBoardOwnership: %w", err)
	}

	return *realUserID == userID, nil
}