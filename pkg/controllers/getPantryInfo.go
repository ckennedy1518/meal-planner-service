package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"meal-planner-service/pkg/db"
	"meal-planner-service/pkg/service"
)

type IngredientRelation struct {
	Name     string `json:"name"`
	IsStaple bool   `json:"is_staple"`
}

type PantryItemResponse struct {
	Amount     int                `json:"amount"`
	Unit       string             `json:"unit"`
	Ingredient IngredientRelation `json:"ingredient"`
}

type PantryItemDTO struct {
	Name     string `json:"name"`
	Quantity int    `json:"quantity"`
	Unit     string `json:"unit"`
	IsStaple bool   `json:"isStaple"`
}

func (h *Handler) GetPantryInfo(w http.ResponseWriter, r *http.Request) {
	claims, ok := service.GetClaims(r)
	if !ok {
		log.Println("Claims: ", claims)
		http.Error(w, "no claims in context", http.StatusUnauthorized)
		return
	}
	userID, ok := claims["id"].(string)
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
	resp, _, err := client.From("pantry_item").
		Select("amount,unit,ingredient(name,is_staple)", "", false).
		Eq("user_id", userID).
		Execute()
	if err != nil {
		log.Printf("Failed to query pantry items: %v", err)
		http.Error(w, "failed to retrieve pantry info", http.StatusInternalServerError)
		return
	}

	var pantryItems []PantryItemResponse
	if err := json.Unmarshal(resp, &pantryItems); err != nil {
		log.Printf("Failed to decode pantry items response: %v", err)
		http.Error(w, "failed to decode pantry info", http.StatusInternalServerError)
		return
	}

	// map query result to return object
	dtoItems := make([]PantryItemDTO, 0, len(pantryItems))
	for _, item := range pantryItems {
		dtoItems = append(dtoItems, PantryItemDTO{
			Name:     item.Ingredient.Name,
			Quantity: item.Amount,
			Unit:     item.Unit,
			IsStaple: item.Ingredient.IsStaple,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message":      "Pantry items retrieved",
		"user_email":   userEmail,
		"pantry_items": dtoItems,
	})
}
