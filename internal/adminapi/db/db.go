// Package db maintains the database connection and extensions.
package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/url"

	"github.com/devpies/core/internal/adminapi/config"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The database driver in use.
	"github.com/pkg/errors"
)

// PostgresRepository represents a database repository.
type PostgresRepository struct {
	*sqlx.DB
	SQ  squirrel.StatementBuilderType
	URL url.URL
}

// NewPostgresRepository creates a new repository, connecting it to the postgres server.
func NewPostgresRepository(cfg config.Config) (*PostgresRepository, error) {
	// Define SSL mode.
	sslMode := "require"
	if cfg.Postgres.DisableTLS {
		sslMode = "disable"
	}

	// Query parameters.
	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")

	// Construct url.
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.Postgres.User, cfg.Postgres.Password),
		Host:     fmt.Sprintf("%s:%d", cfg.Postgres.Host, cfg.Postgres.Port),
		Path:     cfg.Postgres.DB,
		RawQuery: q.Encode(),
	}

	db, err := sqlx.Open("postgres", u.String())
	if err != nil {
		return nil, errors.Wrap(err, "connecting to database")
	}

	r := &PostgresRepository{
		DB:  db,
		SQ:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db),
		URL: u,
	}

	return r, nil
}

// RunInTransaction runs callback function in a transaction.
func (r *PostgresRepository) RunInTransaction(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := r.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	return txRun(tx, fn)
}

func txRun(tx *sqlx.Tx, fn func(*sqlx.Tx) error) error {
	defer func() {
		if err := recover(); err != nil {
			if err := tx.Rollback(); err != nil {
				log.Printf("tx.Rollback panicked: %s", err)
			}
			panic(err)
		}
	}()

	if err := fn(tx); err != nil {
		if err := tx.Rollback(); err != nil {
			log.Printf("tx.Rollback failed: %s", err)
		}
		return err
	}
	return tx.Commit()
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db *PostgresRepository) error {
	// Run a simple query to determine connectivity.
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowxContext(ctx, q).Scan(&tmp)
}
