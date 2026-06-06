package db

import (
	"context"
	"log"
	"time"
)

func NewDBTest() bool {
	pool, err := NewDB()
	if err != nil {
		log.Fatal("Failed to connect to the database (dbtest.go): ", err)
		return false
	}
	defer CloseDB(pool)

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
