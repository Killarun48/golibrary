package service

import (
	"database/sql"
	"errors"
	"golibrary/internal/models"
)

type UserRepositoryer interface {
	CreateUser(user models.User) error
	GetUserByID(id int) (models.User, error)
	GetUsers() []models.User
}

type UserService struct {
	UserRepository UserRepositoryer
}

func NewUserService(userRepository UserRepositoryer) *UserService {
	return &UserService{userRepository}
}

func (s *UserService) CreateUser(user models.User) error {
	return s.UserRepository.CreateUser(user)
}

func (s *UserService) GetUserByID(id int) (models.User, error) {
	user, err := s.UserRepository.GetUserByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, errors.New("пользователь не найден")
		}

		return models.User{}, err
	}

	return user, nil
}

func (s *UserService) GetUsers() []models.User {
	return s.UserRepository.GetUsers()
}
