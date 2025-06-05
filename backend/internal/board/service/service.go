package boardService

import (
	"fmt"
	boardModel "kanban/internal/board/model"
	boardRepo "kanban/internal/board/repo"
	"kanban/internal/utils"
)

type Repository interface {
	Create(board boardModel.Board) error
	GetAll(userID string) ([]boardModel.Board, error)
	Get(boardID string) (*boardModel.Board, error)
	Update(board boardModel.Board) error
	Delete(boardID string) error
}

type Service struct {
	repo *boardRepo.Repository
}

func NewService(repo *boardRepo.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateBoard(userID, name string) error {
	board := boardModel.Board{
		ID:     utils.NewUUID(),
		UserID: userID,
		Name:   name,
	}

	err := s.repo.Create(board)
	if err != nil {
		return fmt.Errorf("boardService.CreateBoard: %w", err)
	}

	return nil
}

func (s *Service) GetAllBoards(userID string) ([]boardModel.Board, error) {
	boards, err := s.repo.GetAll(userID)
	if err != nil {
		return nil, fmt.Errorf("boardService.GetAllBoards: %w", err)
	}

	return boards, nil
}

func (s *Service) GetBoard(boardID string) (*boardModel.Board, error) {
	board, err := s.repo.Get(boardID)
	if err != nil {
		return nil, fmt.Errorf("boardService.GetBoard: %w", err)
	}

	return board, nil
}

func (s *Service) UpdateBoard(boardID, name string) error {
	board := boardModel.Board{
		ID:     boardID,
		Name:   name,
	}

	if err := s.repo.Update(board); err != nil {
		return fmt.Errorf("boardService.UpdateBoard: %w", err)
	}
	
	return nil
}

func (s *Service) DeleteBoard(boardID string) error {
	if err := s.repo.Delete(boardID); err != nil {
		return fmt.Errorf("boardService.DeleteBoard: %w", err)
	}

	return nil
}

func (s *Service) GetUserByBoard(boardID string) (*string, error) {
	return s.repo.GetUserByBoard(boardID)
}