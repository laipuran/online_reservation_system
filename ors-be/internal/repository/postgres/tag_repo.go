package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"ors-be/internal/model"
	"ors-be/internal/repository"
)

type tagRepo struct {
	pool *pgxpool.Pool
}

func NewTagRepo(pool *pgxpool.Pool) repository.TagRepository {
	return &tagRepo{pool: pool}
}

func (r *tagRepo) Create(ctx context.Context, tag *model.Tag) error {
	query := `
		INSERT INTO tags (name)
		VALUES ($1)
		RETURNING id, created_at`

	return r.pool.QueryRow(ctx, query, tag.Name).Scan(&tag.ID, &tag.CreatedAt)
}

func (r *tagRepo) GetByID(ctx context.Context, id int64) (*model.Tag, error) {
	query := `SELECT id, name, created_at FROM tags WHERE id = $1`
	return scanTag(r.pool.QueryRow(ctx, query, id))
}

func (r *tagRepo) GetByName(ctx context.Context, name string) (*model.Tag, error) {
	query := `SELECT id, name, created_at FROM tags WHERE name = $1`
	return scanTag(r.pool.QueryRow(ctx, query, name))
}

func (r *tagRepo) List(ctx context.Context) ([]*model.Tag, error) {
	query := `SELECT id, name, created_at FROM tags ORDER BY id ASC`

	rows, err := r.pool.Query(ctx, query)
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

func scanTag(row pgx.Row) (*model.Tag, error) {
	tag := &model.Tag{}
	err := row.Scan(&tag.ID, &tag.Name, &tag.CreatedAt)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return tag, nil
}
