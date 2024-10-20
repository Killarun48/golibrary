package models

import (
	"database/sql"
	"encoding/json"
)

type Book struct {
	ID       sql.NullInt64  `json:"id,omitempty" example:"21" swaggertype:"integer"`
	Title    sql.NullString `json:"title,omitempty" example:"Shibumi" swaggertype:"string"`
	UserID   sql.NullInt64  `json:"userID,omitempty" example:"12" swaggertype:"integer"`
	AuthorID sql.NullInt64  `json:"authorID,omitempty" example:"12" swaggertype:"integer"`
	Author   Author         `json:"author,omitempty"`
}

func (b *Book) MarshalJSON() ([]byte, error) {
	//type Alias Book
	response := &struct {
		ID       int64   `json:"id,omitempty" example:"0"`
		Title    string  `json:"title,omitempty" example:"johnWick"`
		UserID   int64   `json:"userID,omitempty" example:"0"`
		AuthorID int64   `json:"authorID,omitempty" example:"0"`
		Author   *Author `json:"author,omitempty"`
		//*Alias
	}{
		ID:       b.ID.Int64,
		Title:    b.Title.String,
		UserID:   b.UserID.Int64,
		AuthorID: b.AuthorID.Int64,
		Author:   nil,
		//Alias:    (*Alias)(b),
	}

	if b.Author.ID.Valid {
		response.Author = &b.Author
	}

	/* jsonResp, err := json.MarshalIndent(response, "", "\t")
	if err != nil {
		return nil, err
	}

	return jsonResp, nil */
	return json.Marshal(response)
}

func (b *Book) UnmarshalJSON(data []byte) error {
	aux := &struct {
		Title    string `json:"title,omitempty"`
		AuthorID int64  `json:"authorID,omitempty"`
	}{}

	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	b.Title = sql.NullString{String: aux.Title, Valid: true}
	b.AuthorID = sql.NullInt64{Int64: aux.AuthorID, Valid: true}

	return nil
}
