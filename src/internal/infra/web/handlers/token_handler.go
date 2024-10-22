package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

type TokenHandler struct {
	Client *redis.Client
}

func NewTokenHandler(client *redis.Client) *TokenHandler {
	return &TokenHandler{
		Client: client,
	}
}

type TokenRequest struct {
	Token       string `json:"token"`
	MaxRequests int    `json:"max_requests"`
}

type TokenResponse struct {
	Message string `json:"message"`
}

func (h *TokenHandler) Create(w http.ResponseWriter, r *http.Request) {
	var dto TokenRequest

	err := json.NewDecoder(r.Body).Decode(&dto)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(TokenResponse{
			Message: "Unable to read the body",
		})
	}

	if dto.Token == "" || dto.MaxRequests <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(TokenResponse{
			Message: "Invalid body",
		})
	}

	key := fmt.Sprintf("token_max_req:%s", dto.Token)
	h.Client.Set(context.Background(), key, dto.MaxRequests, 0)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(TokenResponse{
		Message: "Token registered",
	})
}
