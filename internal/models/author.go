package models

import (
	"database/sql"
	"encoding/json"
)

type Author struct {
	ID                 sql.NullInt64  `json:"id,omitempty" example:"0" swaggertype:"integer"`
	Name               sql.NullString `json:"name,omitempty" example:"Rodney William Whitaker" swaggertype:"string"`
	BirthDate          sql.NullString `json:"birthDate,omitempty" example:"1931-06-12" swaggertype:"string"`
	CountOfRentedBooks sql.NullInt64  `json:"countOfRentedBooks,omitempty" example:"0" swaggertype:"integer"`
	Books              []Book         `json:"books,omitempty"`
}

func (a *Author) MarshalJSON() ([]byte, error) {
	//type Alias Author
	response := &struct {
		ID                 int64  `json:"id,omitempty" example:"0"`
		Name               string `json:"name,omitempty" example:"johnWick"`
		BirthDate          string `json:"birthDate,omitempty" example:"johnWick"`
		CountOfRentedBooks int64  `json:"countOfRentedBooks,omitempty" example:"0"`
		Books              []Book `json:"books,omitempty"`
		//*Alias
	}{
		ID:                 a.ID.Int64,
		Name:               a.Name.String,
		BirthDate:          a.BirthDate.String,
		CountOfRentedBooks: a.CountOfRentedBooks.Int64,
		Books:              a.Books,
		//Books:     nil,
		//Alias:     (*Alias)(a),
	}

	/* if len(a.Books) > 0 {
		response.Books = a.Books
	} */

	return json.Marshal(response)
}

func (a *Author) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Name      string `json:"name,omitempty"`
		BirthDate string `json:"birthDate,omitempty"`
	}{}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	a.Name = sql.NullString{String: aux.Name, Valid: true}
	a.BirthDate = sql.NullString{String: aux.BirthDate, Valid: true}

	return nil
}
