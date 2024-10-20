package repository

import (
	"database/sql"
	"golibrary/internal/models"
	"log"

	sq "github.com/Masterminds/squirrel"
)

const (
	booksTable = "books"
)

type BookRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *BookRepository {
	return &BookRepository{db}
}

func (r *BookRepository) CreateBook(book models.Book) error {
	_, err := sq.Insert(booksTable).
		Columns("title", "author_id").
		Values(book.Title, book.AuthorID).
		RunWith(r.db).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func (r *BookRepository) GetBooks() []models.Book {
	var books = make([]models.Book, 0)

	rows, err := sq.Select("books.id", "books.title", "authors.id", "authors.name", "authors.birth_date").
		From(booksTable).
		LeftJoin("authors ON authors.id = books.author_id").
		RunWith(r.db).
		Query()
	if err != nil {
		return make([]models.Book, 0)
	}

	for rows.Next() {
		var book models.Book
		var author models.Author

		err := rows.Scan(
			&book.ID,
			&book.Title,
			&author.ID,
			&author.Name,
			&author.BirthDate,
		)

		if err != nil {
			log.Println(err)
			continue
		}

		book.Author = author

		books = append(books, book)
	}

	return books
}

func (r *BookRepository) GetBookByID(id int) (models.Book, error) {
	var book models.Book
	err := sq.Select("id", "title", "author_id", "user_id").
		From(booksTable).
		Where(sq.Eq{"id": id}).
		RunWith(r.db).
		QueryRow().
		Scan(&book.ID, &book.Title, &book.AuthorID, &book.UserID)
	if err != nil {
		return models.Book{}, err
	}

	return book, nil
}

func (r *BookRepository) RentBook(bookID int, userID int) error {
	_, err := sq.Update(booksTable).
		Set("user_id", userID).
		Where(sq.Eq{"id": bookID}).
		RunWith(r.db).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func (r *BookRepository) ReturnBook(bookID int) error {
	_, err := sq.Update(booksTable).
		Set("user_id", nil).
		Where(sq.Eq{"id": bookID}).
		RunWith(r.db).
		Exec()
	if err != nil {
		return err
	}

	return nil
}
