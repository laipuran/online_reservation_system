package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"ors-be/internal/model"
)

type mockReservationRepo struct {
	userReservation     *model.Reservation
	providerReservation *model.Reservation
	listResult          []*model.Reservation
	updatedStatus       string
	completedDueAt      time.Time
	completedDueResult  []*model.Reservation
	err                 error
}

type mockReservationNotificationService struct {
	notifications []NotificationInput
	err           error
}

func (m *mockReservationNotificationService) Create(ctx context.Context, input NotificationInput) (*model.Notification, error) {
	if m.err != nil {
		return nil, m.err
	}
	m.notifications = append(m.notifications, input)
	return &model.Notification{
		ID:      int64(len(m.notifications)),
		UserID:  input.UserID,
		Title:   input.Title,
		Content: input.Content,
		Type:    input.Type,
	}, nil
}

func (m *mockReservationNotificationService) ListMine(ctx context.Context, userID int64, isRead *bool, limit, offset int) ([]*model.Notification, error) {
	return nil, nil
}

func (m *mockReservationNotificationService) CountUnread(ctx context.Context, userID int64) (int64, error) {
	return 0, nil
}

func (m *mockReservationNotificationService) MarkRead(ctx context.Context, userID, id int64) (*model.Notification, error) {
	return nil, nil
}

func (m *mockReservationNotificationService) MarkAllRead(ctx context.Context, userID int64) (int64, error) {
	return 0, nil
}

func (m *mockReservationRepo) Create(ctx context.Context, reservation *model.Reservation) error {
	if m.err != nil {
		return m.err
	}
	reservation.ID = 1
	return nil
}

func (m *mockReservationRepo) GetByID(ctx context.Context, id int64) (*model.Reservation, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.userReservation, nil
}

func (m *mockReservationRepo) GetByIDForUser(ctx context.Context, id, userID int64) (*model.Reservation, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.userReservation == nil || m.userReservation.ID != id || m.userReservation.UserID != userID {
		return nil, nil
	}
	return m.userReservation, nil
}

func (m *mockReservationRepo) GetByIDForProvider(ctx context.Context, id, providerID int64) (*model.Reservation, error) {
	if m.err != nil {
		return nil, m.err
	}
	if m.providerReservation == nil || m.providerReservation.ID != id {
		return nil, nil
	}
	return m.providerReservation, nil
}

func (m *mockReservationRepo) ListByUserID(ctx context.Context, userID int64, status string, limit, offset int) ([]*model.Reservation, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.listResult, nil
}

func (m *mockReservationRepo) ListByProviderID(ctx context.Context, providerID int64, status string, limit, offset int) ([]*model.Reservation, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.listResult, nil
}

func (m *mockReservationRepo) UpdateStatus(ctx context.Context, id int64, status string) (*model.Reservation, error) {
	if m.err != nil {
		return nil, m.err
	}
	m.updatedStatus = status
	return &model.Reservation{
		ID:        id,
		UserID:    1,
		ServiceID: 2,
		StartTime: time.Now(),
		EndTime:   time.Now().Add(time.Hour),
		Status:    status,
	}, nil
}

func (m *mockReservationRepo) CompleteDue(ctx context.Context, now time.Time) ([]*model.Reservation, error) {
	if m.err != nil {
		return nil, m.err
	}
	m.completedDueAt = now
	return m.completedDueResult, nil
}

func newTestReservationService(repo *mockReservationRepo) ReservationService {
	svc, _ := newTestReservationServiceWithNotifications(repo)
	return svc
}

func newTestReservationServiceWithNotifications(repo *mockReservationRepo) (ReservationService, *mockReservationNotificationService) {
	serviceRepo := newMockServiceRepo()
	serviceRepo.services[2] = &model.Service{
		ID:              2,
		ProviderID:      1,
		CategoryID:      1,
		Title:           "肩颈按摩 60 分钟",
		Price:           199,
		DurationMinutes: 60,
		Status:          "active",
	}

	providerRepo := newMockServiceProviderRepo()
	_ = providerRepo.Create(context.Background(), &model.ServiceProvider{
		UserID:       20,
		BusinessName: "舒心养生馆",
	})

	notificationSvc := &mockReservationNotificationService{}
	return NewReservationService(repo, serviceRepo, providerRepo, notificationSvc), notificationSvc
}

func TestReservationService_Create_Success(t *testing.T) {
	repo := &mockReservationRepo{}
	svc, notificationSvc := newTestReservationServiceWithNotifications(repo)
	start := time.Date(2026, 7, 10, 14, 0, 0, 0, time.UTC)

	reservation, err := svc.Create(context.Background(), 1, ReservationInput{
		ServiceID: 2,
		StartTime: start,
		Note:      " 请准备热水 ",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if reservation.ID == 0 {
		t.Error("Create() ID should be non-zero")
	}
	if reservation.Status != ReservationStatusPending {
		t.Errorf("Create() status = %s, want %s", reservation.Status, ReservationStatusPending)
	}
	if !reservation.EndTime.Equal(start.Add(time.Hour)) {
		t.Errorf("Create() endTime = %s, want %s", reservation.EndTime, start.Add(time.Hour))
	}
	if reservation.Note != "请准备热水" {
		t.Errorf("Create() note = %q", reservation.Note)
	}
	if len(notificationSvc.notifications) != 1 {
		t.Fatalf("notifications count = %d, want 1", len(notificationSvc.notifications))
	}
	notification := notificationSvc.notifications[0]
	if notification.UserID != 20 {
		t.Errorf("notification userID = %d, want provider user 20", notification.UserID)
	}
	if notification.Type != NotificationTypeSystem {
		t.Errorf("notification type = %s, want %s", notification.Type, NotificationTypeSystem)
	}
}

func TestReservationService_CancelMine_Success(t *testing.T) {
	repo := &mockReservationRepo{
		userReservation: &model.Reservation{
			ID:     10,
			UserID: 1,
			Status: ReservationStatusPending,
		},
	}
	svc, notificationSvc := newTestReservationServiceWithNotifications(repo)

	reservation, err := svc.CancelMine(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("CancelMine() error = %v", err)
	}
	if reservation.Status != ReservationStatusCancelled {
		t.Errorf("CancelMine() status = %s, want %s", reservation.Status, ReservationStatusCancelled)
	}
	if repo.updatedStatus != ReservationStatusCancelled {
		t.Errorf("UpdateStatus status = %s, want %s", repo.updatedStatus, ReservationStatusCancelled)
	}
	if len(notificationSvc.notifications) != 1 {
		t.Fatalf("notifications count = %d, want 1", len(notificationSvc.notifications))
	}
	notification := notificationSvc.notifications[0]
	if notification.UserID != 20 {
		t.Errorf("notification userID = %d, want provider user 20", notification.UserID)
	}
	if notification.Type != NotificationTypeReservationCancelled {
		t.Errorf("notification type = %s, want %s", notification.Type, NotificationTypeReservationCancelled)
	}
}

func TestReservationService_CancelMine_InvalidStatus(t *testing.T) {
	repo := &mockReservationRepo{
		userReservation: &model.Reservation{
			ID:     10,
			UserID: 1,
			Status: ReservationStatusCompleted,
		},
	}
	svc := newTestReservationService(repo)

	_, err := svc.CancelMine(context.Background(), 1, 10)
	if !errors.Is(err, ErrReservationCannotCancel) {
		t.Errorf("CancelMine() error = %v, want %v", err, ErrReservationCannotCancel)
	}
	if repo.updatedStatus != "" {
		t.Errorf("UpdateStatus should not be called, got %s", repo.updatedStatus)
	}
}

func TestReservationService_ConfirmForProvider_Success(t *testing.T) {
	repo := &mockReservationRepo{
		providerReservation: &model.Reservation{
			ID:        10,
			UserID:    1,
			ServiceID: 2,
			Status:    ReservationStatusPending,
		},
	}
	svc, notificationSvc := newTestReservationServiceWithNotifications(repo)

	reservation, err := svc.ConfirmForProvider(context.Background(), 20, 10)
	if err != nil {
		t.Fatalf("ConfirmForProvider() error = %v", err)
	}
	if reservation.Status != ReservationStatusConfirmed {
		t.Errorf("ConfirmForProvider() status = %s, want %s", reservation.Status, ReservationStatusConfirmed)
	}
	if len(notificationSvc.notifications) != 1 {
		t.Fatalf("notifications count = %d, want 1", len(notificationSvc.notifications))
	}
	notification := notificationSvc.notifications[0]
	if notification.UserID != 1 {
		t.Errorf("notification userID = %d, want customer user 1", notification.UserID)
	}
	if notification.Type != NotificationTypeReservationConfirmed {
		t.Errorf("notification type = %s, want %s", notification.Type, NotificationTypeReservationConfirmed)
	}
}

func TestReservationService_RejectForProvider_Success(t *testing.T) {
	repo := &mockReservationRepo{
		providerReservation: &model.Reservation{
			ID:        10,
			UserID:    1,
			ServiceID: 2,
			Status:    ReservationStatusPending,
		},
	}
	svc, notificationSvc := newTestReservationServiceWithNotifications(repo)

	reservation, err := svc.RejectForProvider(context.Background(), 20, 10)
	if err != nil {
		t.Fatalf("RejectForProvider() error = %v", err)
	}
	if reservation.Status != ReservationStatusRejected {
		t.Errorf("RejectForProvider() status = %s, want %s", reservation.Status, ReservationStatusRejected)
	}
	if len(notificationSvc.notifications) != 1 {
		t.Fatalf("notifications count = %d, want 1", len(notificationSvc.notifications))
	}
	notification := notificationSvc.notifications[0]
	if notification.UserID != 1 {
		t.Errorf("notification userID = %d, want customer user 1", notification.UserID)
	}
	if notification.Type != NotificationTypeSystem {
		t.Errorf("notification type = %s, want %s", notification.Type, NotificationTypeSystem)
	}
}

func TestReservationService_CompleteDue_Success(t *testing.T) {
	now := time.Date(2026, 7, 10, 15, 0, 0, 0, time.UTC)
	repo := &mockReservationRepo{
		completedDueResult: []*model.Reservation{
			{ID: 1, UserID: 1, ServiceID: 2, Status: ReservationStatusCompleted},
			{ID: 2, UserID: 2, ServiceID: 2, Status: ReservationStatusCompleted},
			{ID: 3, UserID: 3, ServiceID: 2, Status: ReservationStatusCompleted},
		},
	}
	svc, notificationSvc := newTestReservationServiceWithNotifications(repo)

	count, err := svc.CompleteDue(context.Background(), now)
	if err != nil {
		t.Fatalf("CompleteDue() error = %v", err)
	}
	if count != 3 {
		t.Errorf("CompleteDue() count = %d, want 3", count)
	}
	if !repo.completedDueAt.Equal(now) {
		t.Errorf("CompleteDue() now = %s, want %s", repo.completedDueAt, now)
	}
	if len(notificationSvc.notifications) != 3 {
		t.Fatalf("notifications count = %d, want 3", len(notificationSvc.notifications))
	}
	for i, notification := range notificationSvc.notifications {
		if notification.UserID != 20 {
			t.Errorf("notification[%d] userID = %d, want provider user 20", i, notification.UserID)
		}
		if notification.Type != NotificationTypeSystem {
			t.Errorf("notification[%d] type = %s, want %s", i, notification.Type, NotificationTypeSystem)
		}
	}
}

func TestReservationService_ListMine_InvalidStatus(t *testing.T) {
	svc := newTestReservationService(&mockReservationRepo{})

	_, err := svc.ListMine(context.Background(), 1, "unknown", 1, 20)
	if !errors.Is(err, ErrReservationInvalidStatus) {
		t.Errorf("ListMine() error = %v, want %v", err, ErrReservationInvalidStatus)
	}
}

func TestReservationService_GetMine_NotFound(t *testing.T) {
	svc := newTestReservationService(&mockReservationRepo{})

	_, err := svc.GetMine(context.Background(), 1, 10)
	if !errors.Is(err, ErrReservationNotFound) {
		t.Errorf("GetMine() error = %v, want %v", err, ErrReservationNotFound)
	}
}
