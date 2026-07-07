package service

import (
	"context"
	"errors"

	"ors-be/internal/model"
	"ors-be/internal/repository"
)

const (
	ReservationStatusPending   = "pending"
	ReservationStatusConfirmed = "confirmed"
	ReservationStatusCompleted = "completed"
	ReservationStatusCancelled = "cancelled"
	ReservationStatusRejected  = "rejected"
)

var (
	ErrReservationNotFound      = errors.New("reservation not found")
	ErrReservationInvalidStatus = errors.New("reservation status invalid")
)

type ReservationService interface {
	GetMine(ctx context.Context, userID, id int64) (*model.Reservation, error)
	ListMine(ctx context.Context, userID int64, status string, limit, offset int) ([]*model.Reservation, error)
	GetForProvider(ctx context.Context, providerID, id int64) (*model.Reservation, error)
	ListForProvider(ctx context.Context, providerID int64, status string, limit, offset int) ([]*model.Reservation, error)
	CancelMine(ctx context.Context, userID, id int64) (*model.Reservation, error)
	ConfirmForProvider(ctx context.Context, providerID, id int64) (*model.Reservation, error)
	RejectForProvider(ctx context.Context, providerID, id int64) (*model.Reservation, error)
}

type reservationService struct {
	reservationRepo repository.ReservationRepository
}

func NewReservationService(reservationRepo repository.ReservationRepository) ReservationService {
	return &reservationService{reservationRepo: reservationRepo}
}

func (s *reservationService) GetMine(ctx context.Context, userID, id int64) (*model.Reservation, error) {
	reservation, err := s.reservationRepo.GetByIDForUser(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if reservation == nil {
		return nil, ErrReservationNotFound
	}
	return reservation, nil
}

func (s *reservationService) ListMine(ctx context.Context, userID int64, status string, limit, offset int) ([]*model.Reservation, error) {
	if !isReservationStatusFilter(status) {
		return nil, ErrReservationInvalidStatus
	}
	return s.reservationRepo.ListByUserID(ctx, userID, status, limit, offset)
}

func (s *reservationService) GetForProvider(ctx context.Context, providerID, id int64) (*model.Reservation, error) {
	reservation, err := s.reservationRepo.GetByIDForProvider(ctx, id, providerID)
	if err != nil {
		return nil, err
	}
	if reservation == nil {
		return nil, ErrReservationNotFound
	}
	return reservation, nil
}

func (s *reservationService) ListForProvider(ctx context.Context, providerID int64, status string, limit, offset int) ([]*model.Reservation, error) {
	if !isReservationStatusFilter(status) {
		return nil, ErrReservationInvalidStatus
	}
	return s.reservationRepo.ListByProviderID(ctx, providerID, status, limit, offset)
}

func (s *reservationService) CancelMine(ctx context.Context, userID, id int64) (*model.Reservation, error) {
	reservation, err := s.GetMine(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	if reservation.Status != ReservationStatusPending && reservation.Status != ReservationStatusConfirmed {
		return nil, ErrReservationInvalidStatus
	}
	return s.updateStatus(ctx, id, ReservationStatusCancelled)
}

func (s *reservationService) ConfirmForProvider(ctx context.Context, providerID, id int64) (*model.Reservation, error) {
	reservation, err := s.GetForProvider(ctx, providerID, id)
	if err != nil {
		return nil, err
	}
	if reservation.Status != ReservationStatusPending {
		return nil, ErrReservationInvalidStatus
	}
	return s.updateStatus(ctx, id, ReservationStatusConfirmed)
}

func (s *reservationService) RejectForProvider(ctx context.Context, providerID, id int64) (*model.Reservation, error) {
	reservation, err := s.GetForProvider(ctx, providerID, id)
	if err != nil {
		return nil, err
	}
	if reservation.Status != ReservationStatusPending {
		return nil, ErrReservationInvalidStatus
	}
	return s.updateStatus(ctx, id, ReservationStatusRejected)
}

func (s *reservationService) updateStatus(ctx context.Context, id int64, status string) (*model.Reservation, error) {
	reservation, err := s.reservationRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		return nil, err
	}
	if reservation == nil {
		return nil, ErrReservationNotFound
	}
	return reservation, nil
}

func isReservationStatusFilter(status string) bool {
	switch status {
	case "", ReservationStatusPending, ReservationStatusConfirmed, ReservationStatusCompleted, ReservationStatusCancelled, ReservationStatusRejected:
		return true
	default:
		return false
	}
}
