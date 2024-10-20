package service

import (
	"database/sql"
	"errors"
	"golibrary/internal/models"
)

type BookRepositoryer interface {
	CreateBook(book models.Book) error
	GetBooks() []models.Book
	GetBookByID(id int) (models.Book, error)
	RentBook(bookID int, userID int) error
	ReturnBook(bookID int) error
}

type UserServicer interface {
	CreateUser(user models.User) error
	GetUserByID(id int) (models.User, error)
	GetUsers() []models.User
}

type AuthorServicer interface {
	CreateAuthor(user models.Author) error
	GetAuthorByID(id int) (models.Author, error)
	GetAuthors() []models.Author
}

type BookService struct {
	bookRepository BookRepositoryer
	userService    UserServicer
	authorService  AuthorServicer
}

func NewBookService(bookRepository BookRepositoryer, userService UserServicer, authorService AuthorServicer) *BookService {
	return &BookService{
		bookRepository: bookRepository,
		userService:    userService,
		authorService:  authorService,
	}
}

func (s *BookService) CreateBook(book models.Book) error {
	_, err := s.authorService.GetAuthorByID(int(book.AuthorID.Int64))
	if err != nil {
		return err
	}

	return s.bookRepository.CreateBook(book)
}

func (s *BookService) GetBooks() []models.Book {
	return s.bookRepository.GetBooks()
}

func (s *BookService) GetBookByID(id int) (models.Book, error) {
	book, err := s.bookRepository.GetBookByID(id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Book{}, errors.New("книга не найдена")
		}

		return models.Book{}, err
	}

	return book, nil
}

func (s *BookService) RentBook(bookID int, userID int) error {
	book, err := s.GetBookByID(bookID)
	if err != nil {
		return err
	}

	_, err = s.userService.GetUserByID(userID)
	if err != nil {
		return err
	}

	if book.UserID.Int64 == int64(userID) {
		return errors.New("книга уже выдана этому пользователю")
	}

	if book.UserID.Int64 != 0 {
		return errors.New("книга уже выдана")
	}

	return s.bookRepository.RentBook(bookID, userID)
}

func (s *BookService) ReturnBook(bookID int) error {
	book, err := s.GetBookByID(bookID)
	if err != nil {
		return err
	}

	if book.UserID.Int64 == 0 {
		return errors.New("книга не выдавалась")
	}
	
	return s.bookRepository.ReturnBook(bookID)
}