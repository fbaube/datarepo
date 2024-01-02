package utils

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"log"
	// _ "modernc.org/sqlite"
)

// ErrNotFound is returned when no row is found for an SQL query.
var ErrNotFound = errors.New("DB record not found")

func Init() {
	dbString := "file:./generics.db"
	db, err := sql.Open("sqlite", dbString)
	if err != nil {
		log.Fatal(fmt.Errorf("failed to open DB: %w", err))
	}
	if err := db.Ping(); err != nil {
		log.Fatal(fmt.Errorf("failed to ping DB: %w", err))
	}
	R.DB = db
	R.Path = "./generics.db"
}

func Close() {
	_ = R.DB.Close()
}
