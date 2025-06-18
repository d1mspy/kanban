package boardProxy

import (
	"errors"
	"fmt"
	boardModel "kanban/internal/board/model"
)

var ErrForbidden = errors.New("access denied")

type Service interface {
	CreateBoard(userID, name string) error
	GetAllBoards(userID string) ([]boardModel.Board, error)
	GetBoard(boardID string) (*boardModel.Board, error)
	UpdateBoard(boardID, name string) error
	DeleteBoard(boardID string) error
	GetUserByBoard(boardID string) (*string, error)
}

type Proxy struct {
	service Service
}

func NewProxy(service Service) *Proxy {
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
		return nil, fmt.Errorf("boardProxy.GetBoard: %w", ErrForbidden)
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
		return fmt.Errorf("boardProxy.UpdateBoard: %w", ErrForbidden)
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
		return fmt.Errorf("boardProxy.DeleteBoard: %w", ErrForbidden)
	}
}

func (p *Proxy) checkBoardOwnership(boardID, userID string) (bool, error) {
	realUserID, err := p.service.GetUserByBoard(boardID)
	if err != nil {
		return false, fmt.Errorf("boardProxy.checkBoardOwnership: %w", err)
	}

	return *realUserID == userID, nil
}