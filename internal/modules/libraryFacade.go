package modules

import (
	"database/sql"
	"golibrary/internal/infrastructure/db"
	"golibrary/internal/infrastructure/responder"
	"log"

	"github.com/go-chi/chi"
)

type LibraryFacade struct {
	controller *Controller
	service    *Service
	repository *Repository
}

func NewLibraryFacade(pathDB string) *LibraryFacade {
	db := initDB(pathDB)

	repository := NewRepository(db)
	service := NewService(repository)
	controller := NewController(service, responder.NewResponder())

	return &LibraryFacade{
		controller: controller,
		service:    service,
		repository: repository,
	}
}

func (f *LibraryFacade) GetRoutes() *chi.Mux {
	author := f.controller.InitRoutesAuthor()
	book := f.controller.InitRoutesBook()
	user := f.controller.InitRoutesUser()

	r := chi.NewRouter()

	r.Route("/library", func(r chi.Router) {
		r.Mount("/author", author)
		r.Mount("/book", book)
		r.Mount("/user", user)
	})

	return r
}

func initDB(pathDB string) *sql.DB {
	db, err := db.NewDataBaseSqlite(pathDB)
	if err != nil {
		log.Fatal(err)
		return nil
	}

	err = db.Migrate()
	if err != nil {
		log.Fatal(err)
		return nil
	}

	db.FillFakeData()

	return db.DB
}
