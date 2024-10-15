package handlers

import (
	"encoding/json"
	"net/http"
)

type ExampleHandler struct{}

func NewExampleHandler() *ExampleHandler {
	return &ExampleHandler{}
}

type ExampleDto struct {
	Message string `json:"message"`
	Status  int    `json:"status"`
}

func (h *ExampleHandler) Get(w http.ResponseWriter, r *http.Request) {
	dto := ExampleDto{
		Message: "Ok",
		Status:  200,
	}
	json.NewEncoder(w).Encode(dto)
}
