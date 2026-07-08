package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"ors-be/internal/model"
	"ors-be/internal/repository"
)

type userRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) repository.UserRepository {
	return &userRepo{pool: pool}
}

func (r *userRepo) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (name, email, password_hash, role, created_at, updated_at)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
		RETURNING id, created_at, updated_at`

	return r.pool.QueryRow(ctx, query,
		user.Name, user.Email, user.PasswordHash, user.Role,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, name, email, password_hash, role, COALESCE(phone, ''), COALESCE(avatar_url, ''), created_at, updated_at
		FROM users WHERE email = $1`

	user := &model.User{}
	err := r.pool.QueryRow(ctx, query, email).Scan(
		&user.ID, &user.Name, &user.Email,
		&user.PasswordHash, &user.Role,
		&user.Phone, &user.AvatarURL,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	query := `
		SELECT id, name, email, password_hash, role, COALESCE(phone, ''), COALESCE(avatar_url, ''), created_at, updated_at
		FROM users WHERE id = $1`

	user := &model.User{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&user.ID, &user.Name, &user.Email,
		&user.PasswordHash, &user.Role,
		&user.Phone, &user.AvatarURL,
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return user, nil
}

func (r *userRepo) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET name = $1, phone = $2, avatar_url = $3, updated_at = NOW()
		WHERE id = $4
		RETURNING updated_at`

	return r.pool.QueryRow(ctx, query,
		user.Name,
		nullableString(user.Phone),
		nullableString(user.AvatarURL),
		user.ID,
	).Scan(&user.UpdatedAt)
}

func (r *userRepo) UpdatePassword(ctx context.Context, id int64, passwordHash string) error {
	query := `
		UPDATE users
		SET password_hash = $1, updated_at = NOW()
		WHERE id = $2`

	_, err := r.pool.Exec(ctx, query, passwordHash, id)
	return err
}
