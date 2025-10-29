package db

import (
	auth_db "app/internal/db/auth"
	product_db "app/internal/db/product"
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	AuthDBPool     *pgxpool.Pool
	ProductDBPool  *pgxpool.Pool
	AuthQueries    *auth_db.Queries
	ProductQueries *product_db.Queries
)

func Init() {
	ctx := context.Background()

	// Read from environment
	authURL := os.Getenv("AUTH_DB_URL")
	productURL := os.Getenv("PRODUCT_DB_URL")

	if authURL == "" || productURL == "" {
		log.Fatal("Database URLs not set in environment variables")
	}

	// === Connect Auth DB ===
	authPool, err := pgxpool.New(ctx, authURL)
	if err != nil {
		log.Fatalf("Unable to connect to Auth DB: %v", err)
	}
	if err := authPool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping Auth DB: %v", err)
	}
	log.Println("Connected to Auth DB")

	// === Connect Product DB ===
	productPool, err := pgxpool.New(ctx, productURL)
	if err != nil {
		log.Fatalf("Unable to connect to Product DB: %v", err)
	}
	if err := productPool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping Product DB: %v", err)
	}
	log.Println("Connected to Product DB")

	// Assign globals
	AuthDBPool = authPool
	ProductDBPool = productPool
	AuthQueries = auth_db.New(authPool)
	ProductQueries = product_db.New(productPool)
}

// Close closes all database pools
func Close() {
	if AuthDBPool != nil {
		AuthDBPool.Close()
	}
	if ProductDBPool != nil {
		ProductDBPool.Close()
	}
}
