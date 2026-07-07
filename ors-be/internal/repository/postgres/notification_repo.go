package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"ors-be/internal/model"
	"ors-be/internal/repository"
)

type notificationRepo struct {
	pool *pgxpool.Pool
}

func NewNotificationRepo(pool *pgxpool.Pool) repository.NotificationRepository {
	return &notificationRepo{pool: pool}
}

func (r *notificationRepo) Create(ctx context.Context, notification *model.Notification) error {
	query := `
		INSERT INTO notifications (user_id, title, content, type, is_read, created_at)
		VALUES ($1, $2, $3, $4, $5, NOW())
		RETURNING id, is_read, created_at`

	return r.pool.QueryRow(ctx, query,
		notification.UserID,
		notification.Title,
		notification.Content,
		notification.Type,
		notification.IsRead,
	).Scan(&notification.ID, &notification.IsRead, &notification.CreatedAt)
}

func (r *notificationRepo) GetByID(ctx context.Context, id int64) (*model.Notification, error) {
	query := notificationSelectSQL() + ` WHERE id = $1`
	return scanNotification(r.pool.QueryRow(ctx, query, id))
}

func (r *notificationRepo) ListByUserID(ctx context.Context, userID int64, isRead *bool, limit, offset int) ([]*model.Notification, error) {
	limit, offset = normalizeLimitOffset(limit, offset)

	var (
		rows pgx.Rows
		err  error
	)
	if isRead == nil {
		query := notificationSelectSQL() + `
			WHERE user_id = $1
			ORDER BY created_at DESC, id DESC
			LIMIT $2 OFFSET $3`
		rows, err = r.pool.Query(ctx, query, userID, limit, offset)
	} else {
		query := notificationSelectSQL() + `
			WHERE user_id = $1 AND is_read = $2
			ORDER BY created_at DESC, id DESC
			LIMIT $3 OFFSET $4`
		rows, err = r.pool.Query(ctx, query, userID, *isRead, limit, offset)
	}
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanNotifications(rows)
}

func (r *notificationRepo) CountUnread(ctx context.Context, userID int64) (int64, error) {
	query := `SELECT COUNT(*) FROM notifications WHERE user_id = $1 AND is_read = FALSE`

	var count int64
	if err := r.pool.QueryRow(ctx, query, userID).Scan(&count); err != nil {
		return 0, err
	}
	return count, nil
}

func (r *notificationRepo) MarkRead(ctx context.Context, id, userID int64) (*model.Notification, error) {
	query := `
		UPDATE notifications
		SET is_read = TRUE
		WHERE id = $1 AND user_id = $2
		RETURNING id, user_id, title, content, type, is_read, created_at`

	return scanNotification(r.pool.QueryRow(ctx, query, id, userID))
}

func (r *notificationRepo) MarkAllRead(ctx context.Context, userID int64) (int64, error) {
	query := `UPDATE notifications SET is_read = TRUE WHERE user_id = $1 AND is_read = FALSE`

	tag, err := r.pool.Exec(ctx, query, userID)
	if err != nil {
		return 0, err
	}
	return tag.RowsAffected(), nil
}

func notificationSelectSQL() string {
	return `
		SELECT id, user_id, title, content, type, is_read, created_at
		FROM notifications`
}

func scanNotification(row rowScanner) (*model.Notification, error) {
	notification := &model.Notification{}
	err := row.Scan(
		&notification.ID,
		&notification.UserID,
		&notification.Title,
		&notification.Content,
		&notification.Type,
		&notification.IsRead,
		&notification.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return notification, nil
}

func scanNotifications(rows pgx.Rows) ([]*model.Notification, error) {
	notifications := make([]*model.Notification, 0)
	for rows.Next() {
		notification := &model.Notification{}
		if err := rows.Scan(
			&notification.ID,
			&notification.UserID,
			&notification.Title,
			&notification.Content,
			&notification.Type,
			&notification.IsRead,
			&notification.CreatedAt,
		); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return notifications, nil
}
