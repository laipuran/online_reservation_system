package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"

	"ors-be/internal/model"
	"ors-be/internal/repository"
)

type reservationRepo struct {
	pool *pgxpool.Pool
}

func NewReservationRepo(pool *pgxpool.Pool) repository.ReservationRepository {
	return &reservationRepo{pool: pool}
}

func (r *reservationRepo) Create(ctx context.Context, reservation *model.Reservation) error {
	status := reservation.Status
	if status == "" {
		status = "pending"
	}

	query := `
		INSERT INTO reservations (
			user_id, service_id, start_time, end_time, status, note, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, NOW(), NOW())
		RETURNING id, status, created_at, updated_at`

	err := r.pool.QueryRow(ctx, query,
		reservation.UserID,
		reservation.ServiceID,
		reservation.StartTime,
		reservation.EndTime,
		status,
		nullableString(reservation.Note),
	).Scan(
		&reservation.ID,
		&reservation.Status,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
	)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.ConstraintName == "reservations_service_start_time_unique" {
			return repository.ErrReservationTimeConflict
		}
		return err
	}
	return nil
}

func (r *reservationRepo) GetByID(ctx context.Context, id int64) (*model.Reservation, error) {
	query := reservationSelectSQL() + ` WHERE id = $1`
	return scanReservation(r.pool.QueryRow(ctx, query, id))
}

func (r *reservationRepo) GetByIDForUser(ctx context.Context, id, userID int64) (*model.Reservation, error) {
	query := reservationSelectSQL() + ` WHERE id = $1 AND user_id = $2`
	return scanReservation(r.pool.QueryRow(ctx, query, id, userID))
}

func (r *reservationRepo) GetByIDForProvider(ctx context.Context, id, providerID int64) (*model.Reservation, error) {
	query := `
		SELECT r.id, r.user_id, r.service_id, r.start_time, r.end_time, r.status,
			COALESCE(r.note, ''), r.created_at, r.updated_at
		FROM reservations r
		INNER JOIN services s ON s.id = r.service_id
		WHERE r.id = $1 AND s.provider_id = $2`

	return scanReservation(r.pool.QueryRow(ctx, query, id, providerID))
}

func (r *reservationRepo) ListByUserID(ctx context.Context, userID int64, status string, limit, offset int) ([]*model.Reservation, error) {
	limit, offset = normalizeLimitOffset(limit, offset)
	query := reservationSelectSQL() + `
		WHERE user_id = $1 AND ($2 = '' OR status = $2)
		ORDER BY start_time DESC, id DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.pool.Query(ctx, query, userID, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanReservations(rows)
}

func (r *reservationRepo) ListByProviderID(ctx context.Context, providerID int64, status string, limit, offset int) ([]*model.Reservation, error) {
	limit, offset = normalizeLimitOffset(limit, offset)
	query := `
		SELECT r.id, r.user_id, r.service_id, r.start_time, r.end_time, r.status,
			COALESCE(r.note, ''), r.created_at, r.updated_at
		FROM reservations r
		INNER JOIN services s ON s.id = r.service_id
		WHERE s.provider_id = $1 AND ($2 = '' OR r.status = $2)
		ORDER BY r.start_time DESC, r.id DESC
		LIMIT $3 OFFSET $4`

	rows, err := r.pool.Query(ctx, query, providerID, status, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanReservations(rows)
}

func (r *reservationRepo) UpdateStatus(ctx context.Context, id int64, status string) (*model.Reservation, error) {
	query := `
		UPDATE reservations
		SET status = $1, updated_at = NOW()
		WHERE id = $2
		RETURNING id, user_id, service_id, start_time, end_time, status,
			COALESCE(note, ''), created_at, updated_at`

	return scanReservation(r.pool.QueryRow(ctx, query, status, id))
}

func (r *reservationRepo) CompleteDue(ctx context.Context, now time.Time) ([]*model.Reservation, error) {
	query := `
		UPDATE reservations
		SET status = 'completed', updated_at = NOW()
		WHERE status = 'confirmed' AND end_time <= $1
		RETURNING id, user_id, service_id, start_time, end_time, status,
			COALESCE(note, ''), created_at, updated_at`

	rows, err := r.pool.Query(ctx, query, now)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanReservations(rows)
}

func reservationSelectSQL() string {
	return `
		SELECT id, user_id, service_id, start_time, end_time, status,
			COALESCE(note, ''), created_at, updated_at
		FROM reservations`
}

func scanReservation(row rowScanner) (*model.Reservation, error) {
	reservation := &model.Reservation{}
	err := row.Scan(
		&reservation.ID,
		&reservation.UserID,
		&reservation.ServiceID,
		&reservation.StartTime,
		&reservation.EndTime,
		&reservation.Status,
		&reservation.Note,
		&reservation.CreatedAt,
		&reservation.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return reservation, nil
}

func scanReservations(rows pgx.Rows) ([]*model.Reservation, error) {
	reservations := make([]*model.Reservation, 0)
	for rows.Next() {
		reservation := &model.Reservation{}
		if err := rows.Scan(
			&reservation.ID,
			&reservation.UserID,
			&reservation.ServiceID,
			&reservation.StartTime,
			&reservation.EndTime,
			&reservation.Status,
			&reservation.Note,
			&reservation.CreatedAt,
			&reservation.UpdatedAt,
		); err != nil {
			return nil, err
		}
		reservations = append(reservations, reservation)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reservations, nil
}

func normalizeLimitOffset(limit, offset int) (int, int) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
