package controller

import (
	"encoding/json"
	"fmt"
	"golibrary/internal/infrastructure/responder"
	"golibrary/internal/models"
	"net/http"
	"strconv"

	"github.com/go-chi/chi"
)

type BookServicer interface {
	CreateBook(user models.Book) error
	GetBooks() []models.Book
	GetBookByID(id int) (models.Book, error)
	RentBook(bookID int, userID int) error
	ReturnBook(bookID int) error
}

type BookController struct {
	bookService BookServicer
	responder   responder.Responder
}

func NewBookController(bookService BookServicer, respond responder.Responder) *BookController {
	return &BookController{
		bookService: bookService,
		responder:   respond,
	}
}

// @Summary Создание книги
// @Tags book
// @Accept json
// @Produce json
// @Param object body string true "Данные книги" SchemaExample({"title": "Shibumi","authorID": 1})
// @Success 200 {object} responder.Response
// @Router /library/book [post]
func (c *BookController) CreateBook(w http.ResponseWriter, r *http.Request) {
	var book models.Book

	err := json.NewDecoder(r.Body).Decode(&book)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	err = c.bookService.CreateBook(book)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	c.responder.Success(w, "книга создана")
}

// @Summary Получение списка книг
// @Tags book
// @Accept json
// @Produce json
// @Success 200 {object} []models.Book
// @Router /library/book/list [get]
func (c *BookController) GetBooks(w http.ResponseWriter, r *http.Request) {
	books := c.bookService.GetBooks()

	jsonResp, err := json.MarshalIndent(books, "", "    ")
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	fmt.Fprintln(w, string(jsonResp))
}

// @Summary Взятие кники
// @Tags book
// @Accept json
// @Produce json
// @Param bookID path string true "ID книги"
// @Param userID path string true "ID пользователя"
// @Success 200 {object} responder.Response
// @Router /library/book/rent/{bookID}/{userID} [post]
func (c *BookController) RentBook(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookID")
	userID := chi.URLParam(r, "userID")

	bookIDInt, err := strconv.Atoi(bookID)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	err = c.bookService.RentBook(bookIDInt, userIDInt)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	c.responder.Success(w, "книга выдана")
}

// @Summary Возврат книги
// @Tags book
// @Accept json
// @Produce json
// @Param bookID path string true "ID книги"
// @Success 200 {object} responder.Response
// @Router /library/book/return/{bookID} [post]
func (c *BookController) ReturnBook(w http.ResponseWriter, r *http.Request) {
	bookID := chi.URLParam(r, "bookID")

	bookIDInt, err := strconv.Atoi(bookID)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	err = c.bookService.ReturnBook(bookIDInt)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	c.responder.Success(w, "книга возвращена")
}
