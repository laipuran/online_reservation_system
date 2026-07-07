package service

import (
	"context"
	"errors"

	"ors-be/internal/model"
	"ors-be/internal/repository"
)

var ErrReviewInvalidScope = errors.New("review scope invalid")

type ReviewService interface {
	ListByService(ctx context.Context, serviceID int64, limit, offset int) ([]*model.Review, error)
	ListMine(ctx context.Context, userID int64, limit, offset int) ([]*model.Review, error)
	ListByProvider(ctx context.Context, providerID int64, limit, offset int) ([]*model.Review, error)
}

type reviewService struct {
	reviewRepo repository.ReviewRepository
}

func NewReviewService(reviewRepo repository.ReviewRepository) ReviewService {
	return &reviewService{reviewRepo: reviewRepo}
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
