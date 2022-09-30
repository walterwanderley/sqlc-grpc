package main

import (
	"database/sql"
	"embed"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	"github.com/golang-migrate/migrate/v4/source/iofs"
)

//go:embed sql/migrations
var migrations embed.FS

const schemaVersion = 1

func ensureSchema(db *sql.DB) error {
	source, err := iofs.New(migrations, "sql/migrations")
	if err != nil {
		return err
	}
	target, err := sqlite.WithInstance(db, new(sqlite.Config))
	if err != nil {
		return err
	}
	m, err := migrate.NewWithInstance("iofs", source, "sqlite", target)
	if err != nil {
		return err
	}
	err = m.Migrate(schemaVersion)
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return source.Close()
}
