package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"meal-planner-service/pkg/service"
)

type Handler struct {
	AuthService *service.AuthService
}

// NewHandler Constructor to create a Handler object
func NewHandler(auth *service.AuthService) *Handler {
	return &Handler{AuthService: auth}
}

func (h *Handler) Profile(w http.ResponseWriter, r *http.Request) {
	claims, ok := service.GetClaims(r)
	if !ok {
		log.Println("Claims: ", claims)
		http.Error(w, "no claims in context", http.StatusUnauthorized)
		return
	}

	userID := claims["id"]

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message": "Profile retrieved",
		"user_id": userID,
	})
}
