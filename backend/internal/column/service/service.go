package columnService

import (
	"fmt"
	columnModel "kanban/internal/column/model"
	columnRepo "kanban/internal/column/repo"
	"kanban/internal/utils"
)

type Repository interface {
	Create(column columnModel.Column, userID string) error
	GetAll(boardID, userID string) ([]columnModel.Column, error)
	Get(id, userID string) (*columnModel.Column, error)
	Update(userID, columnID string, newName *string, newPos *int) error
	Delete(userID, columnID string) error
}

type Service struct {
	repo *columnRepo.Repository
}

func NewService(repo *columnRepo.Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateColumn(boardID, name, userID string) error {
	column := columnModel.Column{
			ID: utils.NewUUID(),
			BoardID: boardID,
			Name: name,
	}

	err := s.repo.Create(column, userID)
	if err != nil {
		return fmt.Errorf("columnService.CreateColumn: %w", err)
	}

	return nil
}

func (s *Service) GetAllColumns(boardID, userID string) ([]columnModel.Column, error) {
	columns, err := s.repo.GetAll(boardID, userID)
	if err != nil {
		return nil, fmt.Errorf("columnService.GetAllColumns: %w", err)
	}

	return columns, nil
}

func (s *Service) GetColumn(id, userID string) (*columnModel.Column, error) {
	column, err := s.repo.Get(id, userID)
	if err != nil {
		return nil, fmt.Errorf("columnService.GetColumn: %w", err)
	}

	return column, nil
}

func (s *Service) UpdateColumn(userID, columnID string, newName *string, newPos *int) error {
	err := s.repo.Update(userID, columnID, newName, newPos)
	if err != nil {
		return fmt.Errorf("columnService.UpdateColumn: %w", err)
	}

	return nil
}

func (s *Service) DeleteColumn(userID, columnID string) error {
	err := s.repo.Delete(userID, columnID)
	if err != nil {
		return fmt.Errorf("columnService.DeleteColumn: %w", err)
	}

	return nil
}