package modules

import (
	aS "golibrary/internal/modules/author/service"
	bS "golibrary/internal/modules/book/service"
	uS "golibrary/internal/modules/user/service"
)

type Service struct {
	Book   *bS.BookService
	Author *aS.AuthorService
	User   *uS.UserService
}

func NewService(repos *Repository) *Service {
	userService := uS.NewUserService(repos.User)
	authorService := aS.NewAuthorService(repos.Author)
	bookService := bS.NewBookService(repos.Book, userService, authorService)

	s := &Service{
		Book:   bookService,
		Author: authorService,
		User:   userService,
	}

	return s
}
