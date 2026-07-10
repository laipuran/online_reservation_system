package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"ors-be/internal/model"
	"ors-be/internal/repository"
)

type reviewRepo struct {
	pool *pgxpool.Pool
}

func NewReviewRepo(pool *pgxpool.Pool) repository.ReviewRepository {
	return &reviewRepo{pool: pool}
}

func (r *reviewRepo) Create(ctx context.Context, review *model.Review) error {
	return r.CreateAndRefreshServiceRating(ctx, review)
}

func (r *reviewRepo) CreateAndRefreshServiceRating(ctx context.Context, review *model.Review) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var serviceID int64
	if err := tx.QueryRow(ctx, `SELECT id FROM services WHERE id = $1 FOR UPDATE`, review.ServiceID).Scan(&serviceID); err != nil {
		return err
	}

	query := `
		INSERT INTO reviews (reservation_id, user_id, service_id, rating, comment, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, created_at`
	if err := tx.QueryRow(ctx, query,
		review.ReservationID,
		review.UserID,
		review.ServiceID,
		review.Rating,
		nullableString(review.Comment),
	).Scan(&review.ID, &review.CreatedAt); err != nil {
		return err
	}

	updateQuery := `
		UPDATE services
		SET avg_rating = COALESCE((
				SELECT AVG(rating)::DOUBLE PRECISION
				FROM reviews
				WHERE service_id = $1
			), 0),
			review_count = (
				SELECT COUNT(*)
				FROM reviews
				WHERE service_id = $1
			),
			updated_at = NOW()
		WHERE id = $1`
	if _, err := tx.Exec(ctx, updateQuery, review.ServiceID); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (r *reviewRepo) GetByID(ctx context.Context, id int64) (*model.Review, error) {
	query := reviewSelectSQL() + ` WHERE id = $1`
	return scanReview(r.pool.QueryRow(ctx, query, id))
}

func (r *reviewRepo) GetByReservationID(ctx context.Context, reservationID int64) (*model.Review, error) {
	query := reviewSelectSQL() + ` WHERE reservation_id = $1`
	return scanReview(r.pool.QueryRow(ctx, query, reservationID))
}

func (r *reviewRepo) ListByServiceID(ctx context.Context, serviceID int64, limit, offset int) ([]*model.Review, error) {
	limit, offset = normalizeLimitOffset(limit, offset)
	query := reviewSelectSQL() + `
		WHERE service_id = $1
		ORDER BY created_at DESC, id DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, serviceID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanReviews(rows)
}

func (r *reviewRepo) ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]*model.Review, error) {
	limit, offset = normalizeLimitOffset(limit, offset)
	query := reviewSelectSQL() + `
		WHERE user_id = $1
		ORDER BY created_at DESC, id DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, userID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanReviews(rows)
}

func (r *reviewRepo) ListByProviderID(ctx context.Context, providerID int64, limit, offset int) ([]*model.Review, error) {
	limit, offset = normalizeLimitOffset(limit, offset)
	query := `
		SELECT rv.id, rv.reservation_id, rv.user_id, rv.service_id, rv.rating,
			COALESCE(rv.comment, ''), rv.created_at
		FROM reviews rv
		INNER JOIN services s ON s.id = rv.service_id
		WHERE s.provider_id = $1
		ORDER BY rv.created_at DESC, rv.id DESC
		LIMIT $2 OFFSET $3`

	rows, err := r.pool.Query(ctx, query, providerID, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanReviews(rows)
}

func reviewSelectSQL() string {
	return `
		SELECT id, reservation_id, user_id, service_id, rating,
			COALESCE(comment, ''), created_at
		FROM reviews`
}

func scanReview(row rowScanner) (*model.Review, error) {
	review := &model.Review{}
	err := row.Scan(
		&review.ID,
		&review.ReservationID,
		&review.UserID,
		&review.ServiceID,
		&review.Rating,
		&review.Comment,
		&review.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return review, nil
}

func scanReviews(rows pgx.Rows) ([]*model.Review, error) {
	reviews := make([]*model.Review, 0)
	for rows.Next() {
		review := &model.Review{}
		if err := rows.Scan(
			&review.ID,
			&review.ReservationID,
			&review.UserID,
			&review.ServiceID,
			&review.Rating,
			&review.Comment,
			&review.CreatedAt,
		); err != nil {
			return nil, err
		}
		reviews = append(reviews, review)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return reviews, nil
}
