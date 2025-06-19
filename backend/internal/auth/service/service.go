package authService

import (
	"database/sql"
	"errors"
	"fmt"
	authModel "kanban/internal/auth/model"
	"kanban/internal/utils"

	"golang.org/x/crypto/bcrypt"
)

var ErrUserNotFound = errors.New("user not found")
var ErrIncorrectPassword = errors.New("incorrect password")

type Repository interface {
	Create(user authModel.User) error
	GetByUsername(username string) (*authModel.User, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateUser(req authModel.Request) (*string, error){
	hash, err := hashPassword(req.Password)
	if err != nil {
		return nil, fmt.Errorf("authService.CreateUser: %w", err)
	}

	user := authModel.User{
		ID:       utils.NewUUID(),
		Username: req.Username,
		Password: hash,
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, fmt.Errorf("authService.CreateUser: %w", err)
	}

	token, err := GenerateJWT(user.ID)
	if err != nil {
		return nil, fmt.Errorf("authService.CreateUser: %w", err)
	}

	return &token, nil
}

func (s *Service) LoginUser(req authModel.Request) (*string, error) {
	user, err := s.repo.GetByUsername(req.Username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("authService.LoginUser: %w", ErrUserNotFound)
		}
		return nil, fmt.Errorf("authService.LoginUser: %w", err)
	}  

	err = checkPasswordHash(req.Password, user.Password)
	if err != nil {
		return nil, fmt.Errorf("authService.LoginUser: %w", ErrIncorrectPassword)
	}

	token, err := GenerateJWT(user.ID)
	if err != nil {
		return nil, fmt.Errorf("authService.LoginUser: %w", err)
	}

	return &token, nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}
