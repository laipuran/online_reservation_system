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
	err                 error
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

func TestReservationService_CancelMine_Success(t *testing.T) {
	repo := &mockReservationRepo{
		userReservation: &model.Reservation{
			ID:     10,
			UserID: 1,
			Status: ReservationStatusPending,
		},
	}
	svc := NewReservationService(repo)

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
}

func TestReservationService_CancelMine_InvalidStatus(t *testing.T) {
	repo := &mockReservationRepo{
		userReservation: &model.Reservation{
			ID:     10,
			UserID: 1,
			Status: ReservationStatusCompleted,
		},
	}
	svc := NewReservationService(repo)

	_, err := svc.CancelMine(context.Background(), 1, 10)
	if !errors.Is(err, ErrReservationInvalidStatus) {
		t.Errorf("CancelMine() error = %v, want %v", err, ErrReservationInvalidStatus)
	}
	if repo.updatedStatus != "" {
		t.Errorf("UpdateStatus should not be called, got %s", repo.updatedStatus)
	}
}

func TestReservationService_ConfirmForProvider_Success(t *testing.T) {
	repo := &mockReservationRepo{
		providerReservation: &model.Reservation{
			ID:     10,
			Status: ReservationStatusPending,
		},
	}
	svc := NewReservationService(repo)

	reservation, err := svc.ConfirmForProvider(context.Background(), 20, 10)
	if err != nil {
		t.Fatalf("ConfirmForProvider() error = %v", err)
	}
	if reservation.Status != ReservationStatusConfirmed {
		t.Errorf("ConfirmForProvider() status = %s, want %s", reservation.Status, ReservationStatusConfirmed)
	}
}

func TestReservationService_ListMine_InvalidStatus(t *testing.T) {
	svc := NewReservationService(&mockReservationRepo{})

	_, err := svc.ListMine(context.Background(), 1, "unknown", 20, 0)
	if !errors.Is(err, ErrReservationInvalidStatus) {
		t.Errorf("ListMine() error = %v, want %v", err, ErrReservationInvalidStatus)
	}
}

func TestReservationService_GetMine_NotFound(t *testing.T) {
	svc := NewReservationService(&mockReservationRepo{})

	_, err := svc.GetMine(context.Background(), 1, 10)
	if !errors.Is(err, ErrReservationNotFound) {
		t.Errorf("GetMine() error = %v, want %v", err, ErrReservationNotFound)
	}
}
