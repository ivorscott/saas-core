// Package res contains additional resources for services.
package res

import (
	"embed"
	"errors"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres" // required for golang-migrate
	_ "github.com/golang-migrate/migrate/v4/source/file"       // required for golang-migrate
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/lib/pq" // The database driver in use.
)

//go:embed migrations/*.sql
var content embed.FS

// MigrateUp applies the latest database migration.
func MigrateUp(databaseURL string) error {
	d, err := iofs.New(content, "migrations")
	if err != nil {
		return err
	}

	m, err := migrate.NewWithSourceInstance("iofs", d, databaseURL)
	if err != nil {
		println(databaseURL)
		return err
	}
	if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}
