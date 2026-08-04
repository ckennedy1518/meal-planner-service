package controllers

import (
	"encoding/json"
	"log"
	"net/http"

	"meal-planner-service/pkg/db"
	"meal-planner-service/pkg/service"
)

type GroceryListItemResponse struct {
	Amount          int                `json:"amount"`
	AmountPurchased int                `json:"amount_purchased"`
	Unit            string             `json:"unit"`
	Ingredient      IngredientRelation `json:"ingredient"`
}

type GroceryListResponse struct {
	DatePlannedToGo  string                    `json:"date_planned_to_go"`
	GroceryListItems []GroceryListItemResponse `json:"grocery_list_item"`
}

type GroceryListIngredientDTO struct {
	Name              string `json:"name"`
	Quantity          int    `json:"quantity"`
	QuantityPurchased int    `json:"quantity_purchased"`
	Unit              string `json:"unit"`
	IsStaple          bool   `json:"isStaple"`
}

type GroceryListDTO struct {
	Date        string                     `json:"date"`
	Ingredients []GroceryListIngredientDTO `json:"ingredients"`
}

func BuildGroceryListDTOs(groceryLists []GroceryListResponse) []GroceryListDTO {
	dtoLists := make([]GroceryListDTO, 0, len(groceryLists))
	for _, list := range groceryLists {
		ingredients := make([]GroceryListIngredientDTO, 0, len(list.GroceryListItems))
		for _, item := range list.GroceryListItems {
			ingredients = append(ingredients, GroceryListIngredientDTO{
				Name:              item.Ingredient.Name,
				Quantity:          item.Amount,
				QuantityPurchased: item.AmountPurchased,
				Unit:              item.Unit,
				IsStaple:          item.Ingredient.IsStaple,
			})
		}

		dtoLists = append(dtoLists, GroceryListDTO{
			Date:        list.DatePlannedToGo,
			Ingredients: ingredients,
		})
	}
	return dtoLists
}

func (h *Handler) GetGroceryLists(w http.ResponseWriter, r *http.Request) {
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
	resp, _, err := client.From("grocery_list").
		Select("date_planned_to_go,grocery_list_item(amount,amount_purchased,unit,ingredient(name,is_staple))", "", false).
		Eq("user_id", userID).
		Execute()
	if err != nil {
		log.Printf("Failed to query grocery lists: %v", err)
		http.Error(w, "failed to retrieve grocery lists", http.StatusInternalServerError)
		return
	}

	var groceryLists []GroceryListResponse
	if err := json.Unmarshal(resp, &groceryLists); err != nil {
		log.Printf("Failed to decode grocery lists response: %v", err)
		http.Error(w, "failed to decode grocery lists info", http.StatusInternalServerError)
		return
	}

	// map query result to return object
	dtoLists := BuildGroceryListDTOs(groceryLists)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"message":       "Grocery lists retrieved",
		"user_email":    userEmail,
		"grocery_lists": dtoLists,
	})
}
