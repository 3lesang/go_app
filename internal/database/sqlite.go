package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	DB            *sql.DB
	SQliteQueries *Queries
)

func InitDB(filepath string) {
	var err error
	DB, err = sql.Open("sqlite3", filepath)
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}

	if err = DB.Ping(); err != nil {
		log.Fatalf("failed to connect db: %v", err)
	}

	SQliteQueries = New(DB)
	log.Println("Database initialized at", filepath)
}
