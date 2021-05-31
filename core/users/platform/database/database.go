package database

import (
	"context"
	"fmt"
	"net/url"
	"os"

	"github.com/pkg/errors"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq" // The database driver in use.
)

// Config is the required properties to use the database.
type Config struct {
	User       string
	Password   string
	Host       string
	Name       string
	DisableTLS bool
}

type Repository struct {
	DB  *sqlx.DB
	SQ  squirrel.StatementBuilderType
	URL url.URL
}

// NewRepository creates a new Directory, connecting it to the postgres server
func NewRepository(cfg Config) (*Repository, func(), error) {
	// Define SSL mode.
	sslMode := "require"
	if cfg.DisableTLS {
		sslMode = "disable"
	}

	// Query parameters.
	q := make(url.Values)
	q.Set("sslmode", sslMode)
	q.Set("timezone", "utc")

	// Construct url.
	u := url.URL{
		Scheme:   "postgres",
		User:     url.UserPassword(cfg.User, cfg.Password),
		Host:     cfg.Host,
		Path:     cfg.Name,
		RawQuery: q.Encode(),
	}

	db, err := sqlx.Open("postgres", u.String())
	if err != nil {
		return nil, nil, errors.Wrap(err, "connecting to database")
	}

	r := &Repository{
		DB:  db,
		SQ:  squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db),
		URL: u,
	}

	return r, r.Close, nil
}

func (d *Repository) Close() {
	err := d.DB.Close()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// StatusCheck returns nil if it can successfully talk to the database. It
// returns a non-nil error otherwise.
func StatusCheck(ctx context.Context, db *sqlx.DB) error {

	// Run a simple query to determine connectivity. The db has a "Ping" method
	// but it can false-positive when it was previously able to talk to the
	// database but the database has since gone away. Running this query forces a
	// round trip to the database.
	const q = `SELECT true`
	var tmp bool
	return db.QueryRowContext(ctx, q).Scan(&tmp)
}
