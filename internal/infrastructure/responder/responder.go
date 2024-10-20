package responder

import (
	"encoding/json"
	"log"
	"net/http"
)

type Response struct {
	Code    int    `json:"code" example:"200"`
	Success bool   `json:"success" example:"true"`
	Message string `json:"message,omitempty" example:"any message"`
}

type Responder interface {
	ErrorBadRequest(w http.ResponseWriter, err error)
	ErrorNotFound(w http.ResponseWriter, err error)
	Success(w http.ResponseWriter, message string)
}

type Respond struct{}

func NewResponder() Responder {
	return &Respond{}
}

func (r *Respond) ErrorBadRequest(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(Response{
		Code:    http.StatusBadRequest,
		Success: false,
		Message: err.Error(),
	}); err != nil {
		log.Printf("response writer error on write: %v", err.Error())
	}
}

func (r *Respond) ErrorNotFound(w http.ResponseWriter, err error) {
	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(Response{
		Code:    http.StatusNotFound,
		Success: false,
		Message: err.Error(),
	}); err != nil {
		log.Printf("response writer error on write: %v", err.Error())
	}
}

func (r *Respond) Success(w http.ResponseWriter, message string) {
	//w.Header().Set("Content-Type", "application/json;charset=utf-8")
	if err := json.NewEncoder(w).Encode(Response{
		Code:    http.StatusOK,
		Success: true,
		Message: message,
	}); err != nil {
		log.Printf("response writer error on write: %v", err.Error())
	}
}
