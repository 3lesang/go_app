package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var (
	SQliteQueries *Queries
)

func Init() {
	db, err := sql.Open("sqlite3", "./db.sqlite")
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("Cannot connect to database: %v", err)
	}
	SQliteQueries = New(db)
	log.Println("Connected to SQlite!")
}
