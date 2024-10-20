package modules

import (
	"golibrary/internal/infrastructure/responder"
	aC "golibrary/internal/modules/author/controller"
	bC "golibrary/internal/modules/book/controller"
	uC "golibrary/internal/modules/user/controller"
	"net/http"

	"github.com/go-chi/chi"
)

type Controller struct {
	Book   *bC.BookController
	Author *aC.AuthorController
	User   *uC.UserController
}

func NewController(services *Service, responder responder.Responder) *Controller {
	return &Controller{
		Book:   bC.NewBookController(services.Book, responder),
		Author: aC.NewAuthorController(services.Author, responder),
		User:   uC.NewUserController(services.User, responder),
	}
}

func (c *Controller) InitRoutesUser() http.Handler {
	r := chi.NewRouter()

	r.Post("/", c.User.CreateUser)
	r.Get("/list", c.User.GetUsers)

	return r
}

func (c *Controller) InitRoutesBook() http.Handler {
	r := chi.NewRouter()

	r.Post("/", c.Book.CreateBook)
	r.Get("/list", c.Book.GetBooks)
	r.Post("/rent/{bookID}/{userID}", c.Book.RentBook)
	r.Post("/return/{bookID}", c.Book.ReturnBook)

	return r
}

func (c *Controller) InitRoutesAuthor() http.Handler {
	r := chi.NewRouter()

	r.Post("/", c.Author.CreateAuthor)
	r.Get("/list", c.Author.GetAuthors)
	r.Get("/top", c.Author.TopAuthors)

	return r
}
