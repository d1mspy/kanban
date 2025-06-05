package columnService

import (
	"fmt"
	columnModel "kanban/internal/column/model"
	columnRepo "kanban/internal/column/repo"
	"kanban/internal/utils"
)

type Repository interface {
	Create(column columnModel.Column) error
	GetAll(boardID string) ([]columnModel.Column, error)
	Get(columnID string) (*columnModel.Column, error)
	Update(columnID string, newName *string, newPos *int) error
	Delete(columnID string) error
}

type Service struct {
	repo *columnRepo.Repository
}

func NewService(repo *columnRepo.Repository) *Service {
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
		return nil, fmt.Errorf("columnService.GetColumn: %w", err)
	}

	return column, nil
}

func (s *Service) UpdateColumn(columnID string, newName *string, newPos *int) error {
	err := s.repo.Update(columnID, newName, newPos)
	if err != nil {
		return fmt.Errorf("columnService.UpdateColumn: %w", err)
	}

	return nil
}

func (s *Service) DeleteColumn(columnID string) error {
	err := s.repo.Delete(columnID)
	if err != nil {
		return fmt.Errorf("columnService.DeleteColumn: %w", err)
	}

	return nil
}

func (s *Service) GetUserByBoard(boardID string) (*string, error) {
	return s.repo.GetUserByBoard(boardID)
}

func (s *Service) GetUserByColumn(columnID string) (*string, error) {
	return s.repo.GetUserByColumn(columnID)
}