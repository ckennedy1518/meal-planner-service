package db

import (
	"context"
	"encoding/json"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func NewDirectDBTest() bool {
	godotenv.Load()
	myUserId := os.Getenv("MY_USER_ID")

	client, err := NewDirectDB()
	if err != nil {
		log.Fatal("Failed to connect to the database directly (dbtest.go): ", err)
		return false
	}

	ingredient := map[string]interface{}{
		"id":         1,
		"name":       "Salt",
		"created_at": time.Now(),
	}

	ret, _, err := client.From("ingredient").Insert(ingredient, false, "", "representation", "exact").Execute()
	if err != nil {
		log.Fatal("Failed to insert pantryItem via direct connection (dbtest.go): ", err)
		return false
	}

	// Parse JSON response
	var ingredients []map[string]interface{}
	err = json.Unmarshal(ret, &ingredients)
	if err != nil {
		log.Fatal("Failed to parse ingredients response:", err)
		return false
	}

	if len(ingredients) > 0 {
		id := ingredients[0]["id"]
		log.Printf("Inserted ingredient with ID: %v", id)
	}

	pantryItem := map[string]interface{}{
		"id":            1,
		"ingredient_id": 1,
		"amount":        1000,
		"unit":          "tsp",
		"is_staple":     true,
		"user_id":       myUserId,
		"created_at":    time.Now(),
	}

	ret, _, err = client.From("pantry_item").Insert(pantryItem, false, "", "representation", "exact").Execute()
	if err != nil {
		log.Fatal("Failed to insert pantryItem via direct connection (dbtest.go): ", err)
		return false
	}

	// Parse JSON response
	var pantryItems []map[string]interface{}
	err = json.Unmarshal(ret, &pantryItems)
	if err != nil {
		log.Fatal("Failed to parse pantryItems response:", err)
		return false
	}

	if len(pantryItems) > 0 {
		id := pantryItems[0]["id"]
		log.Printf("Inserted pantryItem with ID: %v", id)
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
