package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"meal-planner-service/pkg/db"
	"meal-planner-service/pkg/service"
)

func (h *Handler) GetIngredients(w http.ResponseWriter, r *http.Request) {
	claims, ok := service.GetClaims(r)
	if !ok {
		log.Println("Claims: ", claims)
		http.Error(w, "no claims in context", http.StatusUnauthorized)
		return
	}
	_, ok = claims["id"].(string)
	if !ok {
		log.Println("Claims did not contain a string id")
		http.Error(w, "invalid user id in claims", http.StatusUnauthorized)
		return
	}
	userEmail, ok := claims["email"].(string)
	if !ok {
		log.Println("Claims did not contain a string email")
	}

	client, err := db.NewDirectDB()
	if err != nil {
		log.Printf("Failed to connect to the database directly: %v", err)
		http.Error(w, "failed to connect to database", http.StatusInternalServerError)
		return
	}

	// execute query
	resp, _, err := client.From("ingredient").
		Select("name,is_staple", "", false).
		Execute()
	if err != nil {
		log.Printf("Failed to query ingredients: %v", err)
		http.Error(w, "failed to retrieve ingredients", http.StatusInternalServerError)
		return
	}

	var ingredients []IngredientRelation
	if err := json.Unmarshal(resp, &ingredients); err != nil {
		log.Printf("Failed to decode ingredients response: %v", err)
		http.Error(w, "failed to decode ingredients info", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message":     "Pantry items retrieved",
		"user_email":  userEmail,
		"ingredients": ingredients,
	})
}
