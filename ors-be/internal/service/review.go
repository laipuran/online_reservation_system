package service

import (
	"context"
	"errors"
	"strings"

	"ors-be/internal/model"
	"ors-be/internal/repository"
)

var (
	ErrReviewInvalidScope        = errors.New("review scope invalid")
	ErrReviewInvalidInput        = errors.New("review input invalid")
	ErrReviewInvalidRating       = errors.New("review rating invalid")
	ErrReviewReservationNotFound = errors.New("review reservation not found")
	ErrReviewForbidden           = errors.New("review reservation forbidden")
	ErrReviewReservationNotDone  = errors.New("review reservation not completed")
	ErrReviewAlreadyExists       = errors.New("review already exists")
)

type ReviewInput struct {
	ReservationID int64  `json:"reservation_id"`
	Rating        int16  `json:"rating"`
	Comment       string `json:"comment"`
}

type ReviewService interface {
	Create(ctx context.Context, userID int64, input ReviewInput) (*model.Review, error)
	ListByService(ctx context.Context, serviceID int64, limit, offset int) ([]*model.Review, error)
	ListMine(ctx context.Context, userID int64, limit, offset int) ([]*model.Review, error)
	ListByProvider(ctx context.Context, providerID int64, limit, offset int) ([]*model.Review, error)
}

type reviewService struct {
	reviewRepo      repository.ReviewRepository
	reservationRepo repository.ReservationRepository
}

func NewReviewService(reviewRepo repository.ReviewRepository, reservationRepo repository.ReservationRepository) ReviewService {
	return &reviewService{reviewRepo: reviewRepo, reservationRepo: reservationRepo}
}

func (s *reviewService) Create(ctx context.Context, userID int64, input ReviewInput) (*model.Review, error) {
	if userID <= 0 || input.ReservationID <= 0 {
		return nil, ErrReviewInvalidInput
	}
	if input.Rating < 1 || input.Rating > 5 {
		return nil, ErrReviewInvalidRating
	}

	reservation, err := s.reservationRepo.GetByID(ctx, input.ReservationID)
	if err != nil {
		return nil, err
	}
	if reservation == nil {
		return nil, ErrReviewReservationNotFound
	}
	if reservation.UserID != userID {
		return nil, ErrReviewForbidden
	}
	if reservation.Status != ReservationStatusCompleted {
		return nil, ErrReviewReservationNotDone
	}

	existing, err := s.reviewRepo.GetByReservationID(ctx, input.ReservationID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrReviewAlreadyExists
	}

	review := &model.Review{
		ReservationID: reservation.ID,
		UserID:        userID,
		ServiceID:     reservation.ServiceID,
		Rating:        input.Rating,
		Comment:       strings.TrimSpace(input.Comment),
	}
	if err := s.reviewRepo.CreateAndRefreshServiceRating(ctx, review); err != nil {
		return nil, err
	}
	return review, nil
}

func (s *reviewService) ListByService(ctx context.Context, serviceID int64, limit, offset int) ([]*model.Review, error) {
	if serviceID <= 0 {
		return nil, ErrReviewInvalidScope
	}
	return s.reviewRepo.ListByServiceID(ctx, serviceID, limit, offset)
}

func (s *reviewService) ListMine(ctx context.Context, userID int64, limit, offset int) ([]*model.Review, error) {
	if userID <= 0 {
		return nil, ErrReviewInvalidScope
	}
	return s.reviewRepo.ListByUserID(ctx, userID, limit, offset)
}

func (s *reviewService) ListByProvider(ctx context.Context, providerID int64, limit, offset int) ([]*model.Review, error) {
	if providerID <= 0 {
		return nil, ErrReviewInvalidScope
	}
	return s.reviewRepo.ListByProviderID(ctx, providerID, limit, offset)
}
