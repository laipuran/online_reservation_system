package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"ors-be/internal/model"
	"ors-be/internal/repository"
)

type serviceProviderRepo struct {
	pool *pgxpool.Pool
}

func NewServiceProviderRepo(pool *pgxpool.Pool) repository.ServiceProviderRepository {
	return &serviceProviderRepo{pool: pool}
}

func (r *serviceProviderRepo) Create(ctx context.Context, provider *model.ServiceProvider) error {
	query := `
		INSERT INTO service_providers (
			user_id, business_name, description, address, phone, email, logo_url, created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	return r.pool.QueryRow(ctx, query,
		provider.UserID,
		provider.BusinessName,
		nullableString(provider.Description),
		nullableString(provider.Address),
		nullableString(provider.Phone),
		nullableString(provider.Email),
		nullableString(provider.LogoURL),
	).Scan(&provider.ID, &provider.CreatedAt, &provider.UpdatedAt)
}

func (r *serviceProviderRepo) GetByID(ctx context.Context, id int64) (*model.ServiceProvider, error) {
	query := `
		SELECT id, user_id, business_name,
			COALESCE(description, ''), COALESCE(address, ''), COALESCE(phone, ''),
			COALESCE(email, ''), COALESCE(logo_url, ''), created_at, updated_at
		FROM service_providers WHERE id = $1`

	return scanServiceProvider(r.pool.QueryRow(ctx, query, id))
}

func (r *serviceProviderRepo) GetByUserID(ctx context.Context, userID int64) (*model.ServiceProvider, error) {
	query := `
		SELECT id, user_id, business_name,
			COALESCE(description, ''), COALESCE(address, ''), COALESCE(phone, ''),
			COALESCE(email, ''), COALESCE(logo_url, ''), created_at, updated_at
		FROM service_providers WHERE user_id = $1`

	return scanServiceProvider(r.pool.QueryRow(ctx, query, userID))
}

func (r *serviceProviderRepo) Update(ctx context.Context, provider *model.ServiceProvider) error {
	query := `
		UPDATE service_providers
		SET business_name = $1,
			description = $2,
			address = $3,
			phone = $4,
			email = $5,
			logo_url = $6,
			updated_at = NOW()
		WHERE id = $7
		RETURNING updated_at`

	return r.pool.QueryRow(ctx, query,
		provider.BusinessName,
		nullableString(provider.Description),
		nullableString(provider.Address),
		nullableString(provider.Phone),
		nullableString(provider.Email),
		nullableString(provider.LogoURL),
		provider.ID,
	).Scan(&provider.UpdatedAt)
}

type rowScanner interface {
	Scan(dest ...interface{}) error
}

func scanServiceProvider(row rowScanner) (*model.ServiceProvider, error) {
	provider := &model.ServiceProvider{}
	err := row.Scan(
		&provider.ID,
		&provider.UserID,
		&provider.BusinessName,
		&provider.Description,
		&provider.Address,
		&provider.Phone,
		&provider.Email,
		&provider.LogoURL,
		&provider.CreatedAt,
		&provider.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return provider, nil
}

func nullableString(value string) interface{} {
	if value == "" {
		return nil
	}
	return value
}
