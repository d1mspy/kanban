package columnService

import (
	"database/sql"
	"errors"
	"fmt"
	columnModel "kanban/internal/column/model"
	"kanban/internal/utils"
)

var ErrColumnNotFound = errors.New("column not found")

type Repository interface {
	Create(column columnModel.Column) error
	GetAll(boardID string) ([]columnModel.Column, error)
	Get(columnID string) (*columnModel.Column, error)
	Update(columnID string, newName *string, newPos *int) error
	Delete(columnID string) error
	GetUserByBoard(boardID string) (*string, error)
	GetUserByColumn(columnID string) (*string, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateColumn(boardID, name string) error {
	column := columnModel.Column{
			ID: utils.NewUUID(),
			BoardID: boardID,
			Name: name,
	}

	err := s.repo.Create(column)
	if err != nil {
		return fmt.Errorf("columnService.CreateColumn: %w", err)
	}

	return nil
}

func (s *Service) GetAllColumns(boardID string) ([]columnModel.Column, error) {
	columns, err := s.repo.GetAll(boardID)
	if err != nil {
		return nil, fmt.Errorf("columnService.GetAllColumns: %w", err)
	}

	return columns, nil
}

func (s *Service) GetColumn(boardID string) (*columnModel.Column, error) {
	column, err := s.repo.Get(boardID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("columnService.GetColumn: %w", ErrColumnNotFound)
		}
		return nil, fmt.Errorf("columnService.GetColumn: %w", err)
	}

	return column, nil
}

func (s *Service) UpdateColumn(columnID string, newName *string, newPos *int) error {
	err := s.repo.Update(columnID, newName, newPos)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("columnService.UpdateColumn: %w", ErrColumnNotFound)
		}
		return fmt.Errorf("columnService.UpdateColumn: %w", err)
	}

	return nil
}

func (s *Service) DeleteColumn(columnID string) error {
	err := s.repo.Delete(columnID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("columnService.DeleteColumn: %w", ErrColumnNotFound)
		}
		return fmt.Errorf("columnService.DeleteColumn: %w", err)
	}

	return nil
}

func (s *Service) GetUserByBoard(boardID string) (*string, error) {
	userID, err := s.repo.GetUserByBoard(boardID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("columnService.GetUserByBoard: %w", ErrColumnNotFound)
		}
		return nil, fmt.Errorf("columnService.GetUserByBoard: %w", err)
	}

	return userID, nil
}

func (s *Service) GetUserByColumn(columnID string) (*string, error) {
	userID, err := s.repo.GetUserByColumn(columnID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("columnService.GetUserByColumn: %w", ErrColumnNotFound)
		}
		return nil, fmt.Errorf("columnService.GetUserByColumn: %w", err)
	}

	return userID, nil
}