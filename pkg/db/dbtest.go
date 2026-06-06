package db

import (
	"context"
	"encoding/json"
	"log"
	"time"
)

func NewDirectDBTest() bool {
	client, err := NewDirectDB()
	if err != nil {
		log.Fatal("Failed to connect to the database directly (dbtest.go): ", err)
		return false
	}

	userRow := map[string]interface{}{
		"id":         1,
		"username":   "example_user",
		"email":      "user@example.com",
		"phone":      "555-555-1234",
		"created_at": time.Now(),
	}

	ret, _, err := client.From("user").Insert(userRow, false, "", "representation", "exact").Execute()
	if err != nil {
		log.Fatal("Failed to insert user via direct connection (dbtest.go): ", err)
		return false
	}

	// Parse JSON response
	var users []map[string]interface{}
	err = json.Unmarshal(ret, &users)
	if err != nil {
		log.Fatal("Failed to parse response:", err)
		return false
	}

	if len(users) > 0 {
		id := users[0]["id"]
		log.Printf("Inserted user with ID: %v", id)
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
