package db

import (
	sqlc "app/db/sqlc"
	"database/sql"
	"log"
)

var (
	DB      *sql.DB
	Queries *sqlc.Queries
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

	Queries = sqlc.New(DB)
	log.Println("Database initialized at", filepath)
}
