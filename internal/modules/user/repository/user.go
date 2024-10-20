package repository

import (
	"database/sql"
	"golibrary/internal/models"
	"log"

	sq "github.com/Masterminds/squirrel"
)

const (
	usersTable = "users"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db}
}

func (r UserRepository) CreateUser(user models.User) error {
	_, err := sq.Insert(usersTable).
		Columns("username", "email").
		Values(user.Username, user.Email).
		RunWith(r.db).
		Exec()
	if err != nil {
		return err
	}

	return nil
}

func (r UserRepository) GetUserByID(id int) (models.User, error) {
	var user models.User

	err := sq.Select("id", "username", "email").
		From(usersTable).
		Where(sq.Eq{"id": id}).
		RunWith(r.db).
		QueryRow().
		Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		return models.User{}, err
	}

	return user, nil
}

func (r UserRepository) GetUsers() []models.User {
	var users []*models.User

	rows, err := sq.Select("users.id", "users.username", "users.email", "books.id", "books.title", "books.user_id", "authors.id", "authors.name", "authors.birth_date").
		From(usersTable).
		LeftJoin("books ON books.user_id = users.id").
		LeftJoin("authors ON authors.id = books.author_id").
		RunWith(r.db).
		Query()

	if err != nil {
		log.Println(err)
		return make([]models.User, 0)
	}

	for rows.Next() {
		var user models.User
		var book models.Book
		var author models.Author
		var userID sql.NullInt64

		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&book.ID,
			&book.Title,
			&userID,
			&author.ID,
			&author.Name,
			&author.BirthDate,
		)

		if err != nil {
			log.Println(err)
			continue
		}

		var existingUser *models.User
		for _, u := range users {
			if u.ID == userID {
				existingUser = u
				break
			}
		}

		book.Author = author

		if existingUser != nil {
			existingUser.RentedBooks = append(existingUser.RentedBooks, book)
		} else {
			if book.ID.Int64 == 0 {
				user.RentedBooks = []models.Book{}
				users = append(users, &user)
				continue
			}
			user.RentedBooks = []models.Book{book}
			users = append(users, &user)
		}
	}

	result := make([]models.User, len(users))
	for i, user := range users {
		result[i] = *user
	}

	return result
}
