package service

import (
	"database/sql"
	"errors"
	"golibrary/internal/models"
)

type AuthorRepositoryer interface {
	CreateAuthor(user models.Author) error
	GetAuthorByID(id int) (models.Author, error)
	GetAuthors() []models.Author
	TopAuthors() []models.Author
}

type AuthorService struct {
	AuthorRepository AuthorRepositoryer
}

func NewAuthorService(authorRepository AuthorRepositoryer) *AuthorService {
	return &AuthorService{authorRepository}
}

func (s *AuthorService) CreateAuthor(author models.Author) error {
	return s.AuthorRepository.CreateAuthor(author)
}

func (s *AuthorService) GetAuthorByID(id int) (models.Author, error) {
	author, err := s.AuthorRepository.GetAuthorByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Author{}, errors.New("автор не найден")
		}

		return models.Author{}, err
	}
	
	return author, err
}

func (s *AuthorService) GetAuthors() []models.Author {
	return s.AuthorRepository.GetAuthors()
}

func (s *AuthorService) TopAuthors() []models.Author {
	return s.AuthorRepository.TopAuthors()
}
