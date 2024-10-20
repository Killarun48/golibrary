package db

import (
	"database/sql"
	"fmt"
	"log"

	sq "github.com/Masterminds/squirrel"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	//_ "github.com/glebarez/sqlite"
)

type DataBaseSqlite struct {
	DB *sql.DB
}

func NewDataBaseSqlite(path string) (*DataBaseSqlite, error) {
	db, err := sql.Open("sqlite3", fmt.Sprintf("%s?_fk=1", path))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &DataBaseSqlite{
		DB: db,
	}, nil
}

func (d *DataBaseSqlite) Migrate() error {
	driver, err := sqlite3.WithInstance(d.DB, &sqlite3.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/infrastructure/db/migrations/sqlite3",
		"sqlite3",
		driver)
	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	return nil
}

func (d *DataBaseSqlite) FillFakeData() {
	rows, err := sq.Select("1").
		From("users").
		RunWith(d.DB).
		Query()
	if err != nil {
		log.Println(err)
		return
	}
	defer rows.Close()

	if rows.Next() {
		return
	}

	d.addUsers()
	d.addAuthors()
	d.addBooks()
}

func (d *DataBaseSqlite) addUsers() {
	type fakeUser struct {
		Username string `fake:"{name}"`
		Email    string `fake:"{email}"`
	}

	var fakeUsers = make([]fakeUser, 55)
	for i := range fakeUsers {
		var f fakeUser

		if i == 0 {
			f.Username = "John Wick"
			f.Email = "johnWick@continental.fake"
			fakeUsers[i] = f
			continue
		}

		gofakeit.Struct(&f)
		fakeUsers[i] = f
	}

	insertBuilder := sq.Insert("users").Columns("username", "email")
	for _, f := range fakeUsers {
		insertBuilder = insertBuilder.Values(f.Username, f.Email)
	}
	_, err := insertBuilder.
		RunWith(d.DB).
		Exec()
	if err != nil {
		log.Println(err)
		return
	}
}

func (d *DataBaseSqlite) addAuthors() {
	type fakeAuthor struct {
		Book      gofakeit.BookInfo `fake:"{book}"`
		BirthDate string            `fake:"{year}-{month}-{day}"`
	}

	var fakeAuthors = make([]fakeAuthor, 10)
	for i := range fakeAuthors {
		var f fakeAuthor

		if i == 0 {
			f.Book.Author = "Rodney William Whitaker"
			f.BirthDate = "1931-06-12"
			fakeAuthors[i] = f
			continue
		}

		gofakeit.Struct(&f)
		fakeAuthors[i] = f
	}

	insertBuilder := sq.Insert("authors").Columns("name", "birth_date")
	for _, f := range fakeAuthors {
		insertBuilder = insertBuilder.Values(f.Book.Author, f.BirthDate)
	}
	_, err := insertBuilder.
		RunWith(d.DB).
		Exec()
	if err != nil {
		log.Println(err)
		return
	}
}

func (d *DataBaseSqlite) addBooks() {
	type fakeBook struct {
		AuthorID int               `fake:"{number:2,10}"`
		Book     gofakeit.BookInfo `fake:"{book}"`
	}

	var fakeBooks = make([]fakeBook, 100)
	for i := range fakeBooks {
		var f fakeBook

		if i == 0 {
			f.Book.Title = "Shibumi"
			f.AuthorID = 1
			fakeBooks[i] = f
			continue
		}

		gofakeit.Struct(&f)
		fakeBooks[i] = f
	}

	insertBuilder := sq.Insert("books").Columns("author_id", "title", "user_id")
	for i, f := range fakeBooks {
		if i == 0 {
			insertBuilder = insertBuilder.Values(f.AuthorID, f.Book.Title, 1)
			continue
		}

		insertBuilder = insertBuilder.Values(f.AuthorID, f.Book.Title, nil)
	}
	_, err := insertBuilder.
		RunWith(d.DB).
		Exec()
	if err != nil {
		log.Println(err)
		return
	}
}
