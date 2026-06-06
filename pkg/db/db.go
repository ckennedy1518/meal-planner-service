package db

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"github.com/supabase-community/supabase-go"
)

// NewDirectDB creates and returns a new supabase.Client connected to the database
// or an error. The caller does NOT need to handle closing the client.
func NewDirectDB() (*supabase.Client, error) {
	godotenv.Load()

	supabaseURL := os.Getenv("DATABASE_DIRECT_URL")
	if supabaseURL == "" {
		return nil, errors.New("DATABASE_DIRECT_URL environment variable is not set")
	}

	supabaseKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
	if supabaseKey == "" {
		return nil, errors.New("SUPABASE_SERVICE_ROLE_KEY environment variable is not set")
	}

	client, err := supabase.NewClient(supabaseURL, supabaseKey, &supabase.ClientOptions{})
	if err != nil {
		log.Fatalf("Failed to connect to the database (db.go): %v", err)
		return nil, err
	}
	log.Printf("Supabase URL=%q keyPrefix=%q", supabaseURL, supabaseKey[:28])

	return client, nil
}

// NewPoolDB creates and returns a new pgxpool.Pool connected to the database
// or an error. The caller must handle closing the pool when done.
func NewPoolDB() (*pgxpool.Pool, error) {
	godotenv.Load()

	databaseURL := os.Getenv("DATABASE_POOL_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_POOL_URL environment variable is not set")
	}
	db, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database (db.go): %v", err)
		return nil, err
	}

	return db, db.Ping(context.Background())
}

// ClosePoolDB closes the given pgxpool.Pool if it is not nil.
// Should be called for every NewPoolDB call.
func ClosePoolDB(db *pgxpool.Pool) {
	if db != nil {
		db.Close()
	}
}
