package postgres

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"

	"ors-be/internal/model"
	"ors-be/internal/repository"
)

type serviceTagRepo struct {
	pool *pgxpool.Pool
}

func NewServiceTagRepo(pool *pgxpool.Pool) repository.ServiceTagRepository {
	return &serviceTagRepo{pool: pool}
}

func (r *serviceTagRepo) ReplaceByServiceID(ctx context.Context, serviceID int64, tagIDs []int64) error {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `DELETE FROM service_tags WHERE service_id = $1`, serviceID); err != nil {
		return err
	}

	for _, tagID := range tagIDs {
		if _, err := tx.Exec(ctx,
			`INSERT INTO service_tags (service_id, tag_id) VALUES ($1, $2)`,
			serviceID, tagID,
		); err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *serviceTagRepo) ListByServiceID(ctx context.Context, serviceID int64) ([]*model.Tag, error) {
	query := `
		SELECT t.id, t.name, t.created_at
		FROM service_tags st
		JOIN tags t ON t.id = st.tag_id
		WHERE st.service_id = $1
		ORDER BY t.id ASC`

	rows, err := r.pool.Query(ctx, query, serviceID)
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
