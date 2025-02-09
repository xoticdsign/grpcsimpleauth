package main

import (
	"errors"
	"fmt"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"sso/sso/internal/config"
)

func main() {
	cfg := config.MustLoad()

	if cfg.StoragePath == "" {
		panic("storage path must be set")
	}

	if cfg.MigrationsPath == "" {
		panic("migration path must be set")
	}

	m, err := migrate.New("file://"+cfg.MigrationsPath, fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", cfg.StoragePath, cfg.MigrationsTable))
	if err != nil {
		panic(err)
	}

	err = m.Up()
	if err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			log.Println("nothing to migrate")
		} else {
			panic(err)
		}
	}

	log.Println("migration successful")
}
