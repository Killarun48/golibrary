package models

import (
	"database/sql"
	"encoding/json"
)

type User struct {
	ID          sql.NullInt64  `json:"id,omitempty" example:"0" swaggertype:"integer"`
	Username    sql.NullString `json:"username,omitempty" example:"johnWick" swaggertype:"string"`
	Email       sql.NullString `json:"email,omitempty" example:"johnWick@continental.fake" swaggertype:"string"`
	RentedBooks []Book         `json:"rentedBooks,omitempty"`
}

func (u *User) MarshalJSON() ([]byte, error) {
	type Alias User
	return json.Marshal(&struct {
		ID          int64  `json:"id,omitempty"`
		Username    string `json:"username,omitempty"`
		Email       string `json:"email,omitempty"`
		RentedBooks []Book `json:"rentedBooks,omitempty"`
		*Alias
	}{
		ID:          u.ID.Int64,
		Username:    u.Username.String,
		Email:       u.Email.String,
		RentedBooks: u.RentedBooks,
		Alias:       (*Alias)(u),
	})
}

func (u *User) UnmarshalJSON(data []byte) error {
	//type Alias User
	aux := &struct {
		//*Alias
		Username string `json:"username,omitempty"`
		Email    string `json:"email,omitempty"`
	}{
		//Alias: (*Alias)(u),
	}
	
	err := json.Unmarshal(data, &aux)
	if err != nil {
		return err
	}

	u.Username = sql.NullString{String: aux.Username, Valid: true}
	u.Email = sql.NullString{String: aux.Email, Valid: true}

	return nil
}

/* func (u *User) MarshalJSON() ([]byte, error) {
	var data map[string]interface{}
	data = make(map[string]interface{})

	if u.ID.Valid {
		data["id"] = u.ID.Int64
	}
	if u.Username.Valid {
		data["username"] = u.Username.String
	}
	if u.Email.Valid {
		data["email"] = u.Email.String
	}
	data["rentedBooks"] = u.RentedBooks
	return json.Marshal(data)
} */
