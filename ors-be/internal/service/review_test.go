package service

import (
	"context"
	"errors"
	"testing"

	"ors-be/internal/model"
)

type mockReviewRepo struct {
	reviews []*model.Review
	err     error
}

func (m *mockReviewRepo) Create(ctx context.Context, review *model.Review) error {
	if m.err != nil {
		return m.err
	}
	review.ID = 1
	return nil
}

func (m *mockReviewRepo) GetByID(ctx context.Context, id int64) (*model.Review, error) {
	if m.err != nil {
		return nil, m.err
	}
	return nil, nil
}

func (m *mockReviewRepo) GetByReservationID(ctx context.Context, reservationID int64) (*model.Review, error) {
	if m.err != nil {
		return nil, m.err
	}
	return nil, nil
}

func (m *mockReviewRepo) ListByServiceID(ctx context.Context, serviceID int64, limit, offset int) ([]*model.Review, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.reviews, nil
}

func (m *mockReviewRepo) ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]*model.Review, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.reviews, nil
}

func (m *mockReviewRepo) ListByProviderID(ctx context.Context, providerID int64, limit, offset int) ([]*model.Review, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.reviews, nil
}

func TestReviewService_ListByService_Success(t *testing.T) {
	reviews := []*model.Review{{ID: 1, ServiceID: 2, Rating: 5}}
	svc := NewReviewService(&mockReviewRepo{reviews: reviews})

	result, err := svc.ListByService(context.Background(), 2, 20, 0)
	if err != nil {
		t.Fatalf("ListByService() error = %v", err)
	}
	if len(result) != 1 {
		t.Fatalf("ListByService() len = %d, want 1", len(result))
	}
	if result[0].Rating != 5 {
		t.Errorf("ListByService()[0].Rating = %d, want 5", result[0].Rating)
	}
}

func TestReviewService_ListByService_InvalidScope(t *testing.T) {
	svc := NewReviewService(&mockReviewRepo{})

	_, err := svc.ListByService(context.Background(), 0, 20, 0)
	if !errors.Is(err, ErrReviewInvalidScope) {
		t.Errorf("ListByService() error = %v, want %v", err, ErrReviewInvalidScope)
	}
}

func TestReviewService_ListMine_RepoError(t *testing.T) {
	repoErr := errors.New("db failed")
	svc := NewReviewService(&mockReviewRepo{err: repoErr})

	_, err := svc.ListMine(context.Background(), 1, 20, 0)
	if !errors.Is(err, repoErr) {
		t.Errorf("ListMine() error = %v, want %v", err, repoErr)
	}
}
