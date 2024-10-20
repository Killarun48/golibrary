package controller

import (
	"encoding/json"
	"fmt"
	"golibrary/internal/infrastructure/responder"
	"golibrary/internal/models"
	"net/http"
)

type AuthorServicer interface {
	CreateAuthor(user models.Author) error
	GetAuthors() []models.Author
	TopAuthors() []models.Author
}

type AuthorController struct {
	authorService AuthorServicer
	responder     responder.Responder
}

func NewAuthorController(authorService AuthorServicer, respond responder.Responder) *AuthorController {
	return &AuthorController{
		authorService: authorService,
		responder:     respond,
	}
}

// @Summary Создание автора
// @Tags author
// @Accept json
// @Produce json
// @Param object body string true "Данные автора" SchemaExample({"name": "Rodney William Whitaker","birthDate": "1931-06-12"})
// @Success 200 {object} responder.Response
// @Failure 400 {object} responder.Response
// @Router /library/author [post]
func (c *AuthorController) CreateAuthor(w http.ResponseWriter, r *http.Request) {
	var author models.Author

	err := json.NewDecoder(r.Body).Decode(&author)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	err = c.authorService.CreateAuthor(author)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	c.responder.Success(w, "автор создан")
}

// @Summary Получение списка авторов
// @Tags author
// @Accept json
// @Produce json
// @Success 200 {object} []models.Author
// @Router /library/author/list [get]
func (c *AuthorController) GetAuthors(w http.ResponseWriter, r *http.Request) {
	authors := c.authorService.GetAuthors()
	
	respJson, err := json.MarshalIndent(authors, "", "    ")
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	fmt.Fprintln(w, string(respJson))
}

// @Summary Получение топа авторов
// @Tags author
// @Accept json
// @Produce json
// @Success 200 {object} []models.Author
// @Router /library/author/top [get]
func (c *AuthorController) TopAuthors(w http.ResponseWriter, r *http.Request) {
	authors := c.authorService.TopAuthors()

	respJson, err := json.MarshalIndent(authors, "", "    ")
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	fmt.Fprintln(w, string(respJson))
}