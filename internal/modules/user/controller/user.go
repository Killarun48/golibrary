package controller

import (
	_ "database/sql"
	"encoding/json"
	"fmt"
	"golibrary/internal/infrastructure/responder"
	"golibrary/internal/models"
	"net/http"
)

type UserServicer interface {
	CreateUser(user models.User) error
	GetUsers() []models.User
}

type UserController struct {
	userService UserServicer
	responder   responder.Responder
}

func NewUserController(userService UserServicer, respond responder.Responder) *UserController {
	return &UserController{
		userService: userService,
		responder:   respond,
	}
}

// @Summary Создание пользователя
// @Tags user
// @Accept json
// @Produce json
// @Param object body string true "Данные пользователя" SchemaExample({"username": "John Wick","email": "johnWick@continental.fake"})
// @Success 200 {object} responder.Response
// @Failure 400 {object} responder.Response
// @Router /library/user [post]
func (c *UserController) CreateUser(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	err = c.userService.CreateUser(user)
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	c.responder.Success(w, "пользователь создан")
}

// @Summary Получение списка пользователей
// @Tags user
// @Accept json
// @Produce json
// @Success 200 {object} []models.User
// @Router /library/user/list [get]
func (c *UserController) GetUsers(w http.ResponseWriter, r *http.Request) {
	users := c.userService.GetUsers()

	respJson, err := json.MarshalIndent(users, "", "    ")
	if err != nil {
		c.responder.ErrorBadRequest(w, err)
		return
	}

	fmt.Fprintln(w, string(respJson))
}
