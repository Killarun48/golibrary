package modules

import (
	"database/sql"
	aR "golibrary/internal/modules/author/repository"
	bR "golibrary/internal/modules/book/repository"
	uR "golibrary/internal/modules/user/repository"
)

type Repository struct {
	User   *uR.UserRepository
	Author *aR.AuthorRepository
	Book   *bR.BookRepository
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		User:   uR.NewUserRepository(db),
		Author: aR.NewAuthorRepository(db),
		Book:   bR.NewUserRepository(db),
	}
}
