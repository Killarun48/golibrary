package main

import (
	"fmt"
	"golibrary/internal/infrastructure/responder"
	"golibrary/internal/models"
	"golibrary/internal/modules/author/controller"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

type mockAuthorService struct {
	createAuthor func(user models.Author) error
	getAuthors   func() []models.Author
	topAuthors   func() []models.Author
}

func (m mockAuthorService) CreateAuthor(user models.Author) error {
	return m.createAuthor(user)
}

func (m mockAuthorService) GetAuthors() []models.Author {
	return m.getAuthors()
}

func (m mockAuthorService) TopAuthors() []models.Author {
	return m.topAuthors()
}

func TestAuthor(t *testing.T) {
	tests := []struct {
		name          string
		body          string
		authorService func(user models.Author) error
		wantStatus    int
		wantBody      string
	}{
		{
			name: "valid body",
			body: `{"name": "Rodney William Whitaker","birthDate": "1931-06-12"}`,
			authorService: func(user models.Author) error {
				return nil
			},
			wantStatus: 200,
			wantBody:   "{\"code\":200,\"success\":true,\"message\":\"автор создан\"}\n",
		},
		{
			name: "error create author",
			body: `{"name": "Rodney William Whitaker","birthDate": "1931-06-12"}`,
			authorService: func(user models.Author) error {
				return fmt.Errorf("error")
			},
			wantStatus: 400,
			wantBody:   "{\"code\":400,\"success\":false,\"message\":\"error\"}\n",
		},
		{
			name: "invalid body",
			body: `{"name" "Rodney William Whitaker"}`,
			authorService: func(user models.Author) error {
				return nil
			},
			wantStatus: 400,
			wantBody:   "{\"code\":400,\"success\":false,\"message\":\"invalid character '\\\"' after object key\"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			r := httptest.NewRequest("POST", "/library/author", strings.NewReader(tt.body))
			r.Header.Set("Content-Type", "application/json")

			authorService := mockAuthorService{
				createAuthor: tt.authorService,
			}
			responder := responder.NewResponder()
			ac := controller.NewAuthorController(authorService, responder)
			ac.CreateAuthor(w, r)

			assert.Equal(t, tt.wantStatus, w.Code)
			assert.Equal(t, tt.wantBody, w.Body.String())
		})
	}
}
func TestServer(t *testing.T) {
	pathDB = "./users_test.db"
	s := NewServer(":8088")

	go s.Serve()

	s.Stop()
	time.Sleep(1 * time.Second)

	os.Remove("./users_test.db")
}
