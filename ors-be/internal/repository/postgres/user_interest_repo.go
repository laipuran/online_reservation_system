package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"ors-be/internal/model"
	"ors-be/internal/repository"
)

type userInterestRepo struct {
	pool *pgxpool.Pool
}

func NewUserInterestRepo(pool *pgxpool.Pool) repository.UserInterestRepository {
	return &userInterestRepo{pool: pool}
}

func (r *userInterestRepo) ReplaceByUserID(ctx context.Context, userID int64, tagIDs []int64) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM user_interests WHERE user_id = $1`, userID); err != nil {
		return err
	}

	for _, tagID := range tagIDs {
		if _, err := tx.Exec(ctx,
			`INSERT INTO user_interests (user_id, tag_id) VALUES ($1, $2)`,
			userID, tagID,
		); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *userInterestRepo) ListByUserID(ctx context.Context, userID int64) ([]*model.Tag, error) {
	query := `
		SELECT t.id, t.name, t.created_at
		FROM user_interests ui
		JOIN tags t ON t.id = ui.tag_id
		WHERE ui.user_id = $1
		ORDER BY t.id ASC`

	rows, err := r.pool.Query(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := make([]*model.Tag, 0)
	for rows.Next() {
		tag := &model.Tag{}
		if err := rows.Scan(&tag.ID, &tag.Name, &tag.CreatedAt); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return tags, nil
}
