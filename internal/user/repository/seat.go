package repository

import (
	"context"

	"github.com/devpies/saas-core/internal/user/db"
	"github.com/devpies/saas-core/internal/user/fail"
	"github.com/devpies/saas-core/internal/user/model"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// SeatRepository manages seat data access.
type SeatRepository struct {
	logger *zap.Logger
	pg     *db.PostgresDatabase
}

// NewSeatRepository returns a new seat repository.
func NewSeatRepository(
	logger *zap.Logger,
	pg *db.PostgresDatabase,
) *SeatRepository {
	return &SeatRepository{
		logger: logger,
		pg:     pg,
	}
}

// InsertSeatsEntryTx inserts a new seats entry into the database.
func (sr *SeatRepository) InsertSeatsEntryTx(ctx context.Context, tx *sqlx.Tx, maxSeats model.MaximumSeatsType, tenantID string) error {
	var err error

	stmt := `
		insert into seats (tenant_id, max_seats, seats_used)
		values ($1, $2, $3)
	`
	if _, err = tx.ExecContext(
		ctx,
		stmt,
		tenantID,
		maxSeats,
		0,
	); err != nil {
		return err
	}

	return nil
}

// FindSeatsAvailable retrieves both the number of available seats and the maximum number of allowed seats.
func (sr *SeatRepository) FindSeatsAvailable(ctx context.Context) (model.Seats, error) {
	var (
		s   model.Seats
		err error
	)

	conn, Close, err := sr.pg.GetConnection(ctx)
	if err != nil {
		return s, fail.ErrConnectionFailed
	}
	defer Close()

	stmt := `select max_seats, seats_used from seats`

	if err = conn.GetContext(ctx, &s, stmt); err != nil {
		// Any error is not tolerated.
		return s, err
	}

	return s, nil
}

// IncrementSeatsUsedTx increments the seats used by a tenant.
func (sr *SeatRepository) IncrementSeatsUsedTx(ctx context.Context, tx *sqlx.Tx) error {
	stmt := `update seats set seats_used = seats_used + 1`

	if _, err := tx.ExecContext(ctx, stmt); err != nil {
		return err
	}

	return nil
}

// DecrementSeatsUsedTx decrements the seats used by a tenant.
func (sr *SeatRepository) DecrementSeatsUsedTx(ctx context.Context, tx *sqlx.Tx) error {
	stmt := `update seats set seats_used = seats_used - 1`

	if _, err := tx.ExecContext(ctx, stmt); err != nil {
		return err
	}

	return nil
}
