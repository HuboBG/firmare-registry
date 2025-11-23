package db

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// RunMigrations applies all up migrations in migrationsDir.
// Safe to call on every boot.
func RunMigrations(dbPath, migrationsDir string) {
	source := fmt.Sprintf("file://%s", filepath.Clean(migrationsDir))
	dbURL := fmt.Sprintf("sqlite3://%s", filepath.Clean(dbPath))

	m, err := migrate.New(source, dbURL)
	if err != nil {
		log.Fatal("migrate init failed:", err)
	}
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal("migrate up failed:", err)
	}
	log.Println("migrations OK")
}
