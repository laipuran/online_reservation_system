package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"ors-be/internal/model"
	"ors-be/internal/repository"
)

type categoryRepo struct {
	pool *pgxpool.Pool
}

func NewCategoryRepo(pool *pgxpool.Pool) repository.CategoryRepository {
	return &categoryRepo{pool: pool}
}

func (r *categoryRepo) Create(ctx context.Context, category *model.Category) error {
	query := `
		INSERT INTO categories (name, description, parent_id)
		VALUES ($1, $2, $3)
		RETURNING id, created_at`

	parentID := nullableParentID(category.ParentID)
	return r.pool.QueryRow(ctx, query,
		category.Name, category.Description, parentID,
	).Scan(&category.ID, &category.CreatedAt)
}

func (r *categoryRepo) GetByID(ctx context.Context, id int64) (*model.Category, error) {
	query := `
		SELECT id, name, COALESCE(description, ''), parent_id, created_at
		FROM categories WHERE id = $1`

	category := &model.Category{}
	var parentID sql.NullInt64
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&category.ID, &category.Name, &category.Description,
		&parentID, &category.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	category.ParentID = parentIDPtr(parentID)
	return category, nil
}

func (r *categoryRepo) List(ctx context.Context) ([]*model.Category, error) {
	query := `
		SELECT id, name, COALESCE(description, ''), parent_id, created_at
		FROM categories
		ORDER BY id ASC`

	rows, err := r.pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	categories := make([]*model.Category, 0)
	for rows.Next() {
		category := &model.Category{}
		var parentID sql.NullInt64
		if err := rows.Scan(
			&category.ID, &category.Name, &category.Description,
			&parentID, &category.CreatedAt,
		); err != nil {
			return nil, err
		}
		category.ParentID = parentIDPtr(parentID)
		categories = append(categories, category)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return categories, nil
}

func nullableParentID(parentID *int64) interface{} {
	if parentID == nil {
		return nil
	}
	return *parentID
}

func parentIDPtr(parentID sql.NullInt64) *int64 {
	if !parentID.Valid {
		return nil
	}
	id := parentID.Int64
	return &id
}
