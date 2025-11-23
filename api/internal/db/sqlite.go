package db

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

// OpenSQLite opens SQLite with WAL and busy timeout for better concurrency.
func OpenSQLite(path string) *sql.DB {
	db, err := sql.Open("sqlite3", path+"?_busy_timeout=5000&_journal_mode=WAL")
	if err != nil {
		log.Fatal("sqlite open failed:", err)
	}
	return db
}
