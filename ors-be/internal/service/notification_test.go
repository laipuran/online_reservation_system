package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"ors-be/internal/model"
)

type mockNotificationRepo struct {
	notifications []*model.Notification
	countUnread   int64
	markAllCount  int64
	markReadNil   bool
	err           error
}

func (m *mockNotificationRepo) Create(ctx context.Context, notification *model.Notification) error {
	if m.err != nil {
		return m.err
	}
	notification.ID = int64(len(m.notifications) + 1)
	notification.CreatedAt = time.Now()
	m.notifications = append(m.notifications, notification)
	return nil
}

func (m *mockNotificationRepo) GetByID(ctx context.Context, id int64) (*model.Notification, error) {
	if m.err != nil {
		return nil, m.err
	}
	return nil, nil
}

func (m *mockNotificationRepo) ListByUserID(ctx context.Context, userID int64, isRead *bool, limit, offset int) ([]*model.Notification, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.notifications, nil
}

func (m *mockNotificationRepo) CountUnread(ctx context.Context, userID int64) (int64, error) {
	if m.err != nil {
		return 0, m.err
	}
	return m.countUnread, nil
}

func (m *mockNotificationRepo) MarkRead(ctx context.Context, id, userID int64) (*model.Notification, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.markReadNil {
		return nil, nil
	}
	return &model.Notification{
		ID:     id,
		UserID: userID,
		Title:  "Reservation confirmed",
		Type:   NotificationTypeReservationConfirmed,
		IsRead: true,
	}, nil
}

func (m *mockNotificationRepo) MarkAllRead(ctx context.Context, userID int64) (int64, error) {
	if m.err != nil {
		return 0, m.err
	}
	return m.markAllCount, nil
}

func TestNotificationService_Create_Success(t *testing.T) {
	repo := &mockNotificationRepo{}
	svc := NewNotificationService(repo)

	notification, err := svc.Create(context.Background(), NotificationInput{
		UserID:  1,
		Title:   " Reservation confirmed ",
		Content: " Your reservation is confirmed ",
		Type:    NotificationTypeReservationConfirmed,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if notification.ID == 0 {
		t.Fatal("Create() ID should be non-zero")
	}
	if notification.Title != "Reservation confirmed" {
		t.Errorf("Create() title = %q, want Reservation confirmed", notification.Title)
	}
	if notification.Content != "Your reservation is confirmed" {
		t.Errorf("Create() content = %q, want trimmed content", notification.Content)
	}
	if notification.IsRead {
		t.Error("Create() IsRead should default to false")
	}
}

func TestNotificationService_Create_InvalidType(t *testing.T) {
	svc := NewNotificationService(&mockNotificationRepo{})

	_, err := svc.Create(context.Background(), NotificationInput{
		UserID:  1,
		Title:   "Title",
		Content: "Content",
		Type:    "unknown",
	})
	if !errors.Is(err, ErrNotificationInvalidType) {
		t.Errorf("Create() error = %v, want %v", err, ErrNotificationInvalidType)
	}
}

func TestNotificationService_MarkRead_NotFound(t *testing.T) {
	svc := NewNotificationService(&mockNotificationRepo{markReadNil: true})

	_, err := svc.MarkRead(context.Background(), 1, 10)
	if !errors.Is(err, ErrNotificationNotFound) {
		t.Errorf("MarkRead() error = %v, want %v", err, ErrNotificationNotFound)
	}
}

func TestNotificationService_CountUnread_InvalidUser(t *testing.T) {
	svc := NewNotificationService(&mockNotificationRepo{})

	_, err := svc.CountUnread(context.Background(), 0)
	if !errors.Is(err, ErrNotificationInvalidInput) {
		t.Errorf("CountUnread() error = %v, want %v", err, ErrNotificationInvalidInput)
	}
}
