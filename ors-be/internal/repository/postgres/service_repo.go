package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"ors-be/internal/model"
	"ors-be/internal/repository"
)

type serviceRepo struct {
	pool *pgxpool.Pool
}

func NewServiceRepo(pool *pgxpool.Pool) repository.ServiceRepository {
	return &serviceRepo{pool: pool}
}

func (r *serviceRepo) Create(ctx context.Context, service *model.Service) error {
	query := `
		INSERT INTO services (
			provider_id, category_id, title, description, price, duration_minutes,
			image_url, status, avg_rating, review_count, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, 0, 0, NOW(), NOW())
		RETURNING id, avg_rating, review_count, created_at, updated_at`

	return r.pool.QueryRow(ctx, query,
		service.ProviderID,
		service.CategoryID,
		service.Title,
		nullableString(service.Description),
		service.Price,
		service.DurationMinutes,
		nullableString(service.ImageURL),
		service.Status,
	).Scan(&service.ID, &service.AvgRating, &service.ReviewCount, &service.CreatedAt, &service.UpdatedAt)
}

func (r *serviceRepo) GetByID(ctx context.Context, id int64) (*model.Service, error) {
	query := `
		SELECT id, provider_id, category_id, title, COALESCE(description, ''), price,
			duration_minutes, COALESCE(image_url, ''), status, avg_rating, review_count,
			created_at, updated_at
		FROM services WHERE id = $1`

	service := &model.Service{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&service.ID,
		&service.ProviderID,
		&service.CategoryID,
		&service.Title,
		&service.Description,
		&service.Price,
		&service.DurationMinutes,
		&service.ImageURL,
		&service.Status,
		&service.AvgRating,
		&service.ReviewCount,
		&service.CreatedAt,
		&service.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return service, nil
}

func (r *serviceRepo) GetViewByID(ctx context.Context, id int64) (*model.ServiceView, error) {
	query := serviceViewSelect() + ` WHERE s.id = $1`
	return scanServiceView(r.pool.QueryRow(ctx, query, id))
}

func (r *serviceRepo) List(ctx context.Context, filter model.ServiceFilter) ([]*model.ServiceView, int, error) {
	where, args := buildServiceWhere(filter)

	countQuery := `SELECT COUNT(*) FROM services s` + where
	var total int
	if err := r.pool.QueryRow(ctx, countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	orderBy := serviceOrderBy(filter.SortBy, filter.SortOrder)
	limitArg := len(args) + 1
	offsetArg := len(args) + 2
	args = append(args, filter.PageSize, (filter.Page-1)*filter.PageSize)

	query := serviceViewSelect() + where + orderBy + fmt.Sprintf(" LIMIT $%d OFFSET $%d", limitArg, offsetArg)
	rows, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	items := make([]*model.ServiceView, 0)
	for rows.Next() {
		item, err := scanServiceView(rows)
		if err != nil {
			return nil, 0, err
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}
	return items, total, nil
}

func (r *serviceRepo) Update(ctx context.Context, service *model.Service) error {
	query := `
		UPDATE services
		SET category_id = $1,
			title = $2,
			description = $3,
			price = $4,
			duration_minutes = $5,
			image_url = $6,
			updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at`

	return r.pool.QueryRow(ctx, query,
		service.CategoryID,
		service.Title,
		nullableString(service.Description),
		service.Price,
		service.DurationMinutes,
		nullableString(service.ImageURL),
		service.ID,
	).Scan(&service.UpdatedAt)
}

func (r *serviceRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := `
		UPDATE services
		SET status = $1, updated_at = NOW()
		WHERE id = $2`

	_, err := r.pool.Exec(ctx, query, status, id)
	return err
}

func serviceViewSelect() string {
	return `
		SELECT s.id, s.title, COALESCE(s.description, ''), s.price, s.duration_minutes,
			COALESCE(s.image_url, ''), s.status, s.avg_rating, s.review_count,
			s.created_at, s.updated_at,
			sp.id, sp.business_name,
			c.id, c.name
		FROM services s
		JOIN service_providers sp ON sp.id = s.provider_id
		JOIN categories c ON c.id = s.category_id`
}

func buildServiceWhere(filter model.ServiceFilter) (string, []interface{}) {
	conditions := make([]string, 0)
	args := make([]interface{}, 0)

	addCondition := func(format string, value interface{}) {
		args = append(args, value)
		conditions = append(conditions, fmt.Sprintf(format, len(args)))
	}

	if filter.Keyword != "" {
		args = append(args, "%"+filter.Keyword+"%")
		conditions = append(conditions, fmt.Sprintf("(s.title ILIKE $%d OR s.description ILIKE $%d)", len(args), len(args)))
	}
	if filter.CategoryID != nil {
		addCondition("s.category_id = $%d", *filter.CategoryID)
	}
	if filter.ProviderID != nil {
		addCondition("s.provider_id = $%d", *filter.ProviderID)
	}
	if filter.MinPrice != nil {
		addCondition("s.price >= $%d", *filter.MinPrice)
	}
	if filter.MaxPrice != nil {
		addCondition("s.price <= $%d", *filter.MaxPrice)
	}
	if filter.Status != "" {
		addCondition("s.status = $%d", filter.Status)
	}

	if len(conditions) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(conditions, " AND "), args
}

func serviceOrderBy(sortBy, sortOrder string) string {
	field := map[string]string{
		"price":      "s.price",
		"rating":     "s.avg_rating",
		"created_at": "s.created_at",
	}[sortBy]
	if field == "" {
		field = "s.created_at"
	}

	direction := "DESC"
	if sortOrder == "asc" {
		direction = "ASC"
	}
	return " ORDER BY " + field + " " + direction
}

func scanServiceView(row rowScanner) (*model.ServiceView, error) {
	service := &model.ServiceView{}
	err := row.Scan(
		&service.ID,
		&service.Title,
		&service.Description,
		&service.Price,
		&service.DurationMinutes,
		&service.ImageURL,
		&service.Status,
		&service.AvgRating,
		&service.ReviewCount,
		&service.CreatedAt,
		&service.UpdatedAt,
		&service.Provider.ID,
		&service.Provider.BusinessName,
		&service.Category.ID,
		&service.Category.Name,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return service, nil
}
