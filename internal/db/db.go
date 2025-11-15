package db

import (
	auth_db "app/internal/db/auth"
	blog_db "app/internal/db/blog"
	product_db "app/internal/db/product"
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	AuthDBPool     *pgxpool.Pool
	ProductDBPool  *pgxpool.Pool
	BlogDBPool     *pgxpool.Pool
	AuthQueries    *auth_db.Queries
	ProductQueries *product_db.Queries
	BlogQueries    *blog_db.Queries
)

func Init() {
	ctx := context.Background()

	authURL := os.Getenv("AUTH_DB_URL")
	productURL := os.Getenv("PRODUCT_DB_URL")
	blogURL := os.Getenv("BLOG_DB_URL")

	if authURL == "" || productURL == "" {
		log.Fatal("Database URLs not set in environment variables")
	}

	authPool, err := pgxpool.New(ctx, authURL)
	if err != nil {
		log.Fatalf("Unable to connect to Auth DB: %v", err)
	}
	if err := authPool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping Auth DB: %v", err)
	}
	log.Println("Connected to Auth DB")

	productPool, err := pgxpool.New(ctx, productURL)
	if err != nil {
		log.Fatalf("Unable to connect to Product DB: %v", err)
	}
	if err := productPool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping Product DB: %v", err)
	}
	log.Println("Connected to Product DB")

	blogPool, err := pgxpool.New(ctx, blogURL)
	if err != nil {
		log.Fatalf("Unable to connect to Product DB: %v", err)
	}
	if err := blogPool.Ping(ctx); err != nil {
		log.Fatalf("Unable to ping Product DB: %v", err)
	}
	log.Println("Connected to Blog DB")

	AuthDBPool = authPool
	ProductDBPool = productPool
	BlogDBPool = blogPool
	AuthQueries = auth_db.New(authPool)
	ProductQueries = product_db.New(productPool)
	BlogQueries = blog_db.New(blogPool)
}

func Close() {
	if AuthDBPool != nil {
		AuthDBPool.Close()
	}
	if ProductDBPool != nil {
		ProductDBPool.Close()
	}
	if BlogDBPool != nil {
		BlogDBPool.Close()
	}
}
