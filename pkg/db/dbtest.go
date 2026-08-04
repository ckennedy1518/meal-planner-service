package db

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

func insertIngredient(client *supabase.Client) (bool, error) {
	ingredient := map[string]interface{}{
		"id":         1, // hard code for testing
		"name":       "Salt",
		"created_at": time.Now(),
		"is_staple":  true,
	}

	ret, _, err := client.From("ingredient").Insert(ingredient, false, "", "representation", "exact").Execute()
	if err != nil {
		log.Fatal("Failed to insert pantryItem via direct connection (dbtest.go): ", err)
		return false, err
	}

	// Parse JSON response
	var ingredients []map[string]interface{}
	err = json.Unmarshal(ret, &ingredients)
	if err != nil {
		log.Fatal("Failed to parse ingredients response:", err)
		return false, err
	}

	if len(ingredients) > 0 {
		id := ingredients[0]["id"]
		log.Printf("Inserted ingredient with ID: %v", id)
		return true, nil
	}

	log.Printf("Unable to insert ingredient")
	return false, nil
}

func insertPantryItem(client *supabase.Client, userID string) (bool, error) {
	pantryItem := map[string]interface{}{
		"id":            1,
		"ingredient_id": 1, // hard code for testing
		"amount":        1000,
		"unit":          "tsp",
		"user_id":       userID,
		"created_at":    time.Now(),
	}

	ret, _, err := client.From("pantry_item").Insert(pantryItem, false, "", "representation", "exact").Execute()
	if err != nil {
		log.Fatal("Failed to insert pantryItem via direct connection (dbtest.go): ", err)
		return false, err
	}

	// Parse JSON response
	var pantryItems []map[string]interface{}
	err = json.Unmarshal(ret, &pantryItems)
	if err != nil {
		log.Fatal("Failed to parse pantryItems response:", err)
		return false, err
	}

	if len(pantryItems) > 0 {
		id := pantryItems[0]["id"]
		log.Printf("Inserted pantryItem with ID: %v", id)
		return true, nil
	}

	log.Printf("Unable to insert pantryItem")
	return false, nil
}

func insertGroceryList(client *supabase.Client, userID string) (bool, error) {
	groceryList := map[string]interface{}{
		"id":                 1,
		"user_id":            userID,
		"date_planned_to_go": time.Now().Add(48 * time.Hour),
	}

	ret, _, err := client.From("grocery_list").Insert(groceryList, false, "", "representation", "exact").Execute()
	if err != nil {
		log.Fatal("Failed to insert groceryList via direct connection (dbtest.go): ", err)
		return false, err
	}

	// Parse JSON response
	var groceryLists []map[string]interface{}
	err = json.Unmarshal(ret, &groceryLists)
	if err != nil {
		log.Fatal("Failed to parse groceryLists response:", err)
		return false, err
	}

	if len(groceryLists) > 0 {
		id := groceryLists[0]["id"]
		log.Printf("Inserted groceryList with ID: %v", id)
		return true, nil
	}

	log.Printf("Unable to insert groceryList")
	return false, nil
}

func insertGroceryListItem(client *supabase.Client) (bool, error) {
	groceryListItem := map[string]interface{}{
		"id":              1,
		"grocery_list_id": 1, // hard code for testing
		"ingredient_id":   1, // hard code for testing
		"amount":          3,
		"unit":            "lb",
		// no amount purchased yet
	}

	ret, _, err := client.From("grocery_list_item").Insert(groceryListItem, false, "", "representation", "exact").Execute()
	if err != nil {
		log.Fatal("Failed to insert groceryListItem via direct connection (dbtest.go): ", err)
		return false, err
	}

	// Parse JSON response
	var groceryListItems []map[string]interface{}
	err = json.Unmarshal(ret, &groceryListItems)
	if err != nil {
		log.Fatal("Failed to parse groceryListItems response:", err)
		return false, err
	}

	if len(groceryListItems) > 0 {
		id := groceryListItems[0]["id"]
		log.Printf("Inserted groceryListItem with ID: %v", id)
		return true, nil
	}

	log.Printf("Unable to insert groceryListItem")
	return false, nil
}

func NewDirectDBTest() bool {
	godotenv.Load()
	myUserId := os.Getenv("MY_USER_ID")

	client, err := NewDirectDB()
	if err != nil {
		log.Fatal("Failed to connect to the database directly (dbtest.go): ", err)
		return false
	}

	// result, err := insertIngredient(client)
	// if err != nil || !result {
	// 	log.Fatal("Failed to insert ingredient: ", err)
	// 	return false
	// }

	// result, err = insertPantryItem(client, myUserId)
	// if err != nil || !result {
	// 	log.Fatal("Failed to insert pantry item: ", err)
	// 	return false
	// }

	result, err := insertGroceryList(client, myUserId)
	if err != nil || !result {
		log.Fatal("Failed to insert grocery list: ", err)
		return false
	}

	result, err = insertGroceryListItem(client)
	if err != nil || !result {
		log.Fatal("Failed to insert grocery list item: ", err)
		return false
	}

	return true
}

func NewPoolDBTest() bool {
	pool, err := NewPoolDB()
	if err != nil {
		log.Fatal("Failed to connect to the database with pool (dbtest.go): ", err)
		return false
	}
	defer ClosePoolDB(pool)

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		log.Fatalf("Failed to acquire a database connection (dbtest.go): %v", err)
		return false
	}
	defer conn.Release()

	// Insert a user
	var user string
	err = conn.QueryRow(
		context.Background(),
		`INSERT INTO "user" (id, username, email, phone, created_at)
	     VALUES ($1, $2, $3, $4, $5)
	     RETURNING id`,
		1, "example_user", "user@example.com", "555-555-1234", time.Now(),
	).Scan(&user)
	if err != nil {
		log.Fatalf("Failed to insert user (dbtest.go): %v", err)
	}
	log.Printf(`Inserted user with ID: %s`, user)

	return true
}
