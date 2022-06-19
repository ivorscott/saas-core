// Package db maintains the database connection and extensions.
package db

import (
	"context"
	"fmt"
	"net/url"

	"github.com/devpies/saas-core/internal/user/config"
	"github.com/devpies/saas-core/pkg/web"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The database driver in use.
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// PostgresDatabase represents a database connection.
type PostgresDatabase struct {
	dB     *sqlx.DB
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
		dB:     db,
		URL:    u,
	}

	return r, db.Close, nil
}

// GetConnection returns a tenant aware connection.
func (r *PostgresDatabase) GetConnection(ctx context.Context) (*sqlx.Conn, func() error, error) {
	values, ok := web.FromContext(ctx)
	if !ok {
		r.logger.Error("invalid context values")
		return nil, nil, web.CtxErr()
	}

	conn, err := r.dB.Connx(ctx)
	if err != nil {
		r.logger.Error("connection failed", zap.Error(err))
		_ = conn.Close()
		return nil, nil, err
	}

	stmt := fmt.Sprintf("select set_config('app.current_tenant', '%s', false);", values.TenantID)
	_, err = conn.ExecContext(ctx, stmt)
	if err != nil {
		r.logger.Error("setting session variable failed", zap.Error(err))
		_ = conn.Close()
		return nil, nil, err
	}
	return conn, conn.Close, nil
}

// StatusCheck returns nil if it can successfully talk to the database. It returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, pg *PostgresDatabase) error {
	const q = `SELECT true`
	var tmp bool
	return pg.dB.QueryRowxContext(ctx, q).Scan(&tmp)
}
