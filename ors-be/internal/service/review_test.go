package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"ors-be/internal/model"
)

type mockReviewRepo struct {
	reviews        []*model.Review
	existingReview *model.Review
	createdReview  *model.Review
	err            error
}

func (m *mockReviewRepo) Create(ctx context.Context, review *model.Review) error {
	if m.err != nil {
		return m.err
	}
	review.ID = 1
	m.createdReview = review
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
	return m.existingReview, nil
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

type mockReviewReservationRepo struct {
	reservation *model.Reservation
	err         error
}

func (m *mockReviewReservationRepo) Create(ctx context.Context, reservation *model.Reservation) error {
	return nil
}

func (m *mockReviewReservationRepo) GetByID(ctx context.Context, id int64) (*model.Reservation, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.reservation, nil
}

func (m *mockReviewReservationRepo) GetByIDForUser(ctx context.Context, id, userID int64) (*model.Reservation, error) {
	return nil, nil
}

func (m *mockReviewReservationRepo) GetByIDForProvider(ctx context.Context, id, providerID int64) (*model.Reservation, error) {
	return nil, nil
}

func (m *mockReviewReservationRepo) ListByUserID(ctx context.Context, userID int64, status string, limit, offset int) ([]*model.Reservation, error) {
	return nil, nil
}

func (m *mockReviewReservationRepo) ListByProviderID(ctx context.Context, providerID int64, status string, limit, offset int) ([]*model.Reservation, error) {
	return nil, nil
}

func (m *mockReviewReservationRepo) UpdateStatus(ctx context.Context, id int64, status string) (*model.Reservation, error) {
	return nil, nil
}

func (m *mockReviewReservationRepo) CompleteDue(ctx context.Context, now time.Time) ([]*model.Reservation, error) {
	return nil, nil
}

func TestReviewService_ListByService_Success(t *testing.T) {
	reviews := []*model.Review{{ID: 1, ServiceID: 2, Rating: 5}}
	svc := NewReviewService(&mockReviewRepo{reviews: reviews}, nil)

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
	svc := NewReviewService(&mockReviewRepo{}, nil)

	_, err := svc.ListByService(context.Background(), 0, 20, 0)
	if !errors.Is(err, ErrReviewInvalidScope) {
		t.Errorf("ListByService() error = %v, want %v", err, ErrReviewInvalidScope)
	}
}

func TestReviewService_ListMine_RepoError(t *testing.T) {
	repoErr := errors.New("db failed")
	svc := NewReviewService(&mockReviewRepo{err: repoErr}, nil)

	_, err := svc.ListMine(context.Background(), 1, 20, 0)
	if !errors.Is(err, repoErr) {
		t.Errorf("ListMine() error = %v, want %v", err, repoErr)
	}
}

func TestReviewService_Create_Success(t *testing.T) {
	reviewRepo := &mockReviewRepo{}
	svc := NewReviewService(reviewRepo, &mockReviewReservationRepo{
		reservation: &model.Reservation{
			ID:        10,
			UserID:    1,
			ServiceID: 3,
			Status:    ReservationStatusCompleted,
			StartTime: time.Now(),
			EndTime:   time.Now().Add(time.Hour),
		},
	})

	review, err := svc.Create(context.Background(), 1, ReviewInput{
		ReservationID: 10,
		Rating:        5,
		Comment:       " great service ",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if review.ID != 1 {
		t.Errorf("Create() ID = %d, want 1", review.ID)
	}
	if review.ServiceID != 3 {
		t.Errorf("Create() ServiceID = %d, want 3", review.ServiceID)
	}
	if review.Comment != "great service" {
		t.Errorf("Create() Comment = %q, want trimmed comment", review.Comment)
	}
	if reviewRepo.createdReview == nil {
		t.Fatal("Create() did not call review repo")
	}
}

func TestReviewService_Create_ReservationNotFound(t *testing.T) {
	svc := NewReviewService(&mockReviewRepo{}, &mockReviewReservationRepo{})

	_, err := svc.Create(context.Background(), 1, ReviewInput{ReservationID: 10, Rating: 5})
	if !errors.Is(err, ErrReviewReservationNotFound) {
		t.Errorf("Create() error = %v, want %v", err, ErrReviewReservationNotFound)
	}
}

func TestReviewService_Create_ReservationForbidden(t *testing.T) {
	svc := NewReviewService(&mockReviewRepo{}, &mockReviewReservationRepo{
		reservation: &model.Reservation{ID: 10, UserID: 2, ServiceID: 3, Status: ReservationStatusCompleted},
	})

	_, err := svc.Create(context.Background(), 1, ReviewInput{ReservationID: 10, Rating: 5})
	if !errors.Is(err, ErrReviewForbidden) {
		t.Errorf("Create() error = %v, want %v", err, ErrReviewForbidden)
	}
}

func TestReviewService_Create_ReservationNotCompleted(t *testing.T) {
	svc := NewReviewService(&mockReviewRepo{}, &mockReviewReservationRepo{
		reservation: &model.Reservation{ID: 10, UserID: 1, ServiceID: 3, Status: ReservationStatusConfirmed},
	})

	_, err := svc.Create(context.Background(), 1, ReviewInput{ReservationID: 10, Rating: 5})
	if !errors.Is(err, ErrReviewReservationNotDone) {
		t.Errorf("Create() error = %v, want %v", err, ErrReviewReservationNotDone)
	}
}

func TestReviewService_Create_AlreadyExists(t *testing.T) {
	svc := NewReviewService(&mockReviewRepo{
		existingReview: &model.Review{ID: 99, ReservationID: 10},
	}, &mockReviewReservationRepo{
		reservation: &model.Reservation{ID: 10, UserID: 1, ServiceID: 3, Status: ReservationStatusCompleted},
	})

	_, err := svc.Create(context.Background(), 1, ReviewInput{ReservationID: 10, Rating: 5})
	if !errors.Is(err, ErrReviewAlreadyExists) {
		t.Errorf("Create() error = %v, want %v", err, ErrReviewAlreadyExists)
	}
}

func TestReviewService_Create_InvalidRating(t *testing.T) {
	svc := NewReviewService(&mockReviewRepo{}, &mockReviewReservationRepo{})

	_, err := svc.Create(context.Background(), 1, ReviewInput{ReservationID: 10, Rating: 6})
	if !errors.Is(err, ErrReviewInvalidRating) {
		t.Errorf("Create() error = %v, want %v", err, ErrReviewInvalidRating)
	}
}
