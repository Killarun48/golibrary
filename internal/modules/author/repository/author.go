package repository

import (
	"database/sql"
	"golibrary/internal/models"
	"log"

	sq "github.com/Masterminds/squirrel"
)

const (
	authorsTable = "authors"
	booksTable   = "books"
)

type AuthorRepository struct {
	db *sql.DB
}

func NewAuthorRepository(db *sql.DB) *AuthorRepository {
	return &AuthorRepository{db}
}

func (r AuthorRepository) CreateAuthor(user models.Author) error {
	_, err := sq.Insert(authorsTable).
		Columns("name", "birth_date").
		Values(user.Name, user.BirthDate).
		RunWith(r.db).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func (r AuthorRepository) GetAuthorByID(id int) (models.Author, error) {
	var author models.Author

	err := sq.Select("id", "name", "birth_date").
		From(authorsTable).
		Where(sq.Eq{"id": id}).
		RunWith(r.db).
		QueryRow().
		Scan(&author.ID, &author.Name, &author.BirthDate)
	if err != nil {
		return models.Author{}, err
	}

	return author, nil
}

func (r AuthorRepository) GetAuthors() []models.Author {
	var authors []*models.Author

	rows, err := sq.Select("authors.id", "authors.name", "authors.birth_date", "books.id", "books.title").
		From(authorsTable).
		LeftJoin("books ON books.author_id = authors.id").
		RunWith(r.db).
		Query()
	if err != nil {
		return make([]models.Author, 0)
	}

	for rows.Next() {
		var author models.Author
		var book models.Book

		err := rows.Scan(
			&author.ID,
			&author.Name,
			&author.BirthDate,
			&book.ID,
			&book.Title,
		)
		if err != nil {
			log.Println(err)
			continue
		}

		var existingAuthor *models.Author
		for _, a := range authors {
			if a.ID == author.ID {
				existingAuthor = a
				break
			}
		}

		if existingAuthor != nil {
			existingAuthor.Books = append(existingAuthor.Books, book)
		} else {
			if book.ID.Int64 == 0 {
				author.Books = []models.Book{}
				authors = append(authors, &author)
				continue
			}
			author.Books = []models.Book{book}
			authors = append(authors, &author)
		}
	}

	result := make([]models.Author, len(authors))
	for i, author := range authors {
		result[i] = *author
	}

	return result
}

func (r AuthorRepository) TopAuthors() []models.Author {
	var authors = make([]models.Author, 0)
	rows, err := sq.Select("books.author_id", "authors.name", "authors.birth_date", "count(books.author_id) as count_of_rent").
		From(booksTable).
		LeftJoin("authors ON authors.id = books.author_id").
		Where(sq.NotEq{"books.user_id": nil}).
		GroupBy("books.author_id", "authors.name", "authors.birth_date").
		OrderBy("count_of_rent DESC").
		RunWith(r.db).
		Query()
	if err != nil {
		return make([]models.Author, 0)
	}

	for rows.Next() {
		var author models.Author

		err := rows.Scan(&author.ID, &author.Name, &author.BirthDate, &author.CountOfRentedBooks)
		if err != nil {
			log.Println(err)
			continue
		}

		authors = append(authors, author)
	}
	
	return authors
}
