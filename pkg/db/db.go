package db

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

// NewDB creates and returns a new pgxpool.Pool connected to the database
// or an error. The caller must handle closing the pool when done.
func NewDB() (*pgxpool.Pool, error) {
	godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return nil, errors.New("DATABASE_URL environment variable is not set")
	}
	db, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to the database (db.go): %v", err)
		return nil, err
	}

	return db, db.Ping(context.Background())
}

// CloseDB closes the given pgxpool.Pool if it is not nil.
// Should be called for every NewDB call.
func CloseDB(db *pgxpool.Pool) {
	if db != nil {
		db.Close()
	}
}
