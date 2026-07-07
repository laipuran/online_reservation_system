package service

import (
	"context"
	"errors"
	"strings"

	"ors-be/internal/model"
	"ors-be/internal/repository"
)

const (
	NotificationTypeReservationConfirmed = "reservation_confirmed"
	NotificationTypeReservationCancelled = "reservation_cancelled"
	NotificationTypeReservationReminder  = "reservation_reminder"
	NotificationTypeReviewReceived       = "review_received"
	NotificationTypeSystem               = "system"
)

var (
	ErrNotificationNotFound     = errors.New("notification not found")
	ErrNotificationInvalidInput = errors.New("notification input invalid")
	ErrNotificationInvalidType  = errors.New("notification type invalid")
)

type NotificationInput struct {
	UserID  int64  `json:"user_id"`
	Title   string `json:"title"`
	Content string `json:"content"`
	Type    string `json:"type"`
}

type NotificationService interface {
	Create(ctx context.Context, input NotificationInput) (*model.Notification, error)
	ListMine(ctx context.Context, userID int64, isRead *bool, limit, offset int) ([]*model.Notification, error)
	CountUnread(ctx context.Context, userID int64) (int64, error)
	MarkRead(ctx context.Context, userID, id int64) (*model.Notification, error)
	MarkAllRead(ctx context.Context, userID int64) (int64, error)
}

type notificationService struct {
	notificationRepo repository.NotificationRepository
}

func NewNotificationService(notificationRepo repository.NotificationRepository) NotificationService {
	return &notificationService{notificationRepo: notificationRepo}
}

func (s *notificationService) Create(ctx context.Context, input NotificationInput) (*model.Notification, error) {
	notification := normalizeNotificationInput(input)
	if notification.UserID <= 0 || notification.Title == "" || notification.Content == "" {
		return nil, ErrNotificationInvalidInput
	}
	if !isNotificationType(notification.Type) {
		return nil, ErrNotificationInvalidType
	}
	if err := s.notificationRepo.Create(ctx, notification); err != nil {
		return nil, err
	}
	return notification, nil
}

func (s *notificationService) ListMine(ctx context.Context, userID int64, isRead *bool, limit, offset int) ([]*model.Notification, error) {
	if userID <= 0 {
		return nil, ErrNotificationInvalidInput
	}
	return s.notificationRepo.ListByUserID(ctx, userID, isRead, limit, offset)
}

func (s *notificationService) CountUnread(ctx context.Context, userID int64) (int64, error) {
	if userID <= 0 {
		return 0, ErrNotificationInvalidInput
	}
	return s.notificationRepo.CountUnread(ctx, userID)
}

func (s *notificationService) MarkRead(ctx context.Context, userID, id int64) (*model.Notification, error) {
	if userID <= 0 || id <= 0 {
		return nil, ErrNotificationInvalidInput
	}
	notification, err := s.notificationRepo.MarkRead(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if notification == nil {
		return nil, ErrNotificationNotFound
	}
	return notification, nil
}

func (s *notificationService) MarkAllRead(ctx context.Context, userID int64) (int64, error) {
	if userID <= 0 {
		return 0, ErrNotificationInvalidInput
	}
	return s.notificationRepo.MarkAllRead(ctx, userID)
}

func normalizeNotificationInput(input NotificationInput) *model.Notification {
	return &model.Notification{
		UserID:  input.UserID,
		Title:   strings.TrimSpace(input.Title),
		Content: strings.TrimSpace(input.Content),
		Type:    strings.TrimSpace(input.Type),
	}
}

func isNotificationType(notificationType string) bool {
	switch notificationType {
	case NotificationTypeReservationConfirmed,
		NotificationTypeReservationCancelled,
		NotificationTypeReservationReminder,
		NotificationTypeReviewReceived,
		NotificationTypeSystem:
		return true
	default:
		return false
	}
}
