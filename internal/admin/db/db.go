// Package db maintains the database connection and extensions.
package db

import (
	"context"
	"database/sql"
	"fmt"
	"net/url"

	"github.com/devpies/saas-core/internal/admin/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The database driver in use.
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// PostgresDatabase represents a database connection.
type PostgresDatabase struct {
	DB     *sqlx.DB
	logger *zap.Logger
	URL    url.URL
}

// NewPostgresDatabase creates a new postgres database.
func NewPostgresDatabase(logger *zap.Logger, cfg config.Config) (*PostgresDatabase, func() error, error) {
	sslMode := "require"
	if cfg.DB.DisableTLS {
		sslMode = "disable"
	}

	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")

	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.DB.User, cfg.DB.Password),
		Host:     fmt.Sprintf("%s:%d", cfg.DB.Host, cfg.DB.Port),
		Path:     cfg.DB.Name,
		RawQuery: q.Encode(),
	}

	db, err := sqlx.Open("postgres", u.String())
	if err != nil {
		return nil, nil, errors.Wrap(err, "connecting to database")
	}

	r := &PostgresDatabase{
		logger: logger,
		DB:     db,
		URL:    u,
	}

	return r, db.Close, nil
}

// RunInTransaction runs callback function in a transaction.
func (pg *PostgresDatabase) RunInTransaction(ctx context.Context, fn func(*sqlx.Tx) error) error {
	tx, err := pg.DB.BeginTxx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return err
	}
	return pg.txRun(tx, fn)
}

func (pg *PostgresDatabase) txRun(tx *sqlx.Tx, fn func(*sqlx.Tx) error) error {
	defer func() {
		if err := recover(); err != nil {
			if rbErr := tx.Rollback(); rbErr != nil {
				pg.logger.Info("tx.Rollback panicked", zap.Error(rbErr))
			}
			panic(err)
		}
	}()

	if err := fn(tx); err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			pg.logger.Info("tx.Rollback failed", zap.Error(rbErr))
		}
		return err
	}
	return tx.Commit()
}

// StatusCheck returns nil if it can successfully talk to the database. It returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, pg *PostgresDatabase) error {
	const q = `SELECT true`
	var tmp bool
	return pg.DB.QueryRowxContext(ctx, q).Scan(&tmp)
}
