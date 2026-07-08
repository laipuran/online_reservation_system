package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"ors-be/internal/api/http/middleware"
	"ors-be/internal/auth"
	"ors-be/internal/model"
	"ors-be/internal/service"
)

type mockReviewService struct {
	createFn         func(ctx context.Context, userID int64, input service.ReviewInput) (*model.Review, error)
	listByServiceFn  func(ctx context.Context, serviceID int64, limit, offset int) ([]*model.Review, error)
	listMineFn       func(ctx context.Context, userID int64, limit, offset int) ([]*model.Review, error)
	listByProviderFn func(ctx context.Context, providerID int64, limit, offset int) ([]*model.Review, error)
}

func (m *mockReviewService) Create(ctx context.Context, userID int64, input service.ReviewInput) (*model.Review, error) {
	return m.createFn(ctx, userID, input)
}

func (m *mockReviewService) ListByService(ctx context.Context, serviceID int64, limit, offset int) ([]*model.Review, error) {
	return m.listByServiceFn(ctx, serviceID, limit, offset)
}

func (m *mockReviewService) ListMine(ctx context.Context, userID int64, limit, offset int) ([]*model.Review, error) {
	return m.listMineFn(ctx, userID, limit, offset)
}

func (m *mockReviewService) ListByProvider(ctx context.Context, providerID int64, limit, offset int) ([]*model.Review, error) {
	return m.listByProviderFn(ctx, providerID, limit, offset)
}

func TestReviewHandler_ListByService_Success(t *testing.T) {
	h := NewReviewHandler(&mockReviewService{
		listByServiceFn: func(ctx context.Context, serviceID int64, limit, offset int) ([]*model.Review, error) {
			if serviceID != 3 {
				t.Errorf("serviceID = %d, want 3", serviceID)
			}
			if limit != 10 || offset != 10 {
				t.Errorf("limit/offset = %d/%d, want 10/10", limit, offset)
			}
			return []*model.Review{{ID: 1, ServiceID: serviceID, Rating: 5}}, nil
		},
	})

	req := requestWithURLParam(http.MethodGet, "/api/v1/services/3/reviews?page=2&page_size=10", "id", "3", nil)
	w := httptest.NewRecorder()

	h.ListByService()(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	var resp struct {
		Code int `json:"code"`
		Data struct {
			Items    []model.Review `json:"items"`
			Page     int            `json:"page"`
			PageSize int            `json:"page_size"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Code != 200 || len(resp.Data.Items) != 1 || resp.Data.Page != 2 || resp.Data.PageSize != 10 {
		t.Fatalf("response = %+v", resp)
	}
}

func TestReviewHandler_ListMine_Success(t *testing.T) {
	h := NewReviewHandler(&mockReviewService{
		listMineFn: func(ctx context.Context, userID int64, limit, offset int) ([]*model.Review, error) {
			if userID != 7 {
				t.Errorf("userID = %d, want 7", userID)
			}
			return []*model.Review{{ID: 2, UserID: userID, Rating: 4}}, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me/reviews", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserCtxKey, &auth.Claims{UserID: 7, Role: "customer"}))
	w := httptest.NewRecorder()

	h.ListMine()(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestReviewHandler_Create_Success(t *testing.T) {
	h := NewReviewHandler(&mockReviewService{
		createFn: func(ctx context.Context, userID int64, input service.ReviewInput) (*model.Review, error) {
			if userID != 7 {
				t.Errorf("userID = %d, want 7", userID)
			}
			if input.ReservationID != 10 || input.Rating != 5 {
				t.Errorf("input = %+v", input)
			}
			return &model.Review{ID: 1, ReservationID: input.ReservationID, UserID: userID, ServiceID: 3, Rating: input.Rating, CreatedAt: time.Now()}, nil
		},
	})

	body := bytes.NewBufferString(`{"reservation_id":10,"rating":5,"comment":"good"}`)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reviews", body)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserCtxKey, &auth.Claims{UserID: 7, Role: "customer"}))
	w := httptest.NewRecorder()

	h.Create()(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusCreated)
	}
}

func TestReviewHandler_Create_InvalidRating(t *testing.T) {
	h := NewReviewHandler(&mockReviewService{
		createFn: func(ctx context.Context, userID int64, input service.ReviewInput) (*model.Review, error) {
			return nil, service.ErrReviewInvalidRating
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/reviews", bytes.NewBufferString(`{"reservation_id":10,"rating":6}`))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserCtxKey, &auth.Claims{UserID: 7, Role: "customer"}))
	w := httptest.NewRecorder()

	h.Create()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestReviewHandler_Create_ReservationNotFound(t *testing.T) {
	h := NewReviewHandler(&mockReviewService{
		createFn: func(ctx context.Context, userID int64, input service.ReviewInput) (*model.Review, error) {
			return nil, service.ErrReviewReservationNotFound
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/reviews", bytes.NewBufferString(`{"reservation_id":10,"rating":5}`))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserCtxKey, &auth.Claims{UserID: 7, Role: "customer"}))
	w := httptest.NewRecorder()

	h.Create()(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestReviewHandler_Create_AlreadyExists(t *testing.T) {
	h := NewReviewHandler(&mockReviewService{
		createFn: func(ctx context.Context, userID int64, input service.ReviewInput) (*model.Review, error) {
			return nil, service.ErrReviewAlreadyExists
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/reviews", bytes.NewBufferString(`{"reservation_id":10,"rating":5}`))
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserCtxKey, &auth.Claims{UserID: 7, Role: "customer"}))
	w := httptest.NewRecorder()

	h.Create()(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusConflict)
	}
}

func TestReviewHandler_ListByProvider_InvalidID(t *testing.T) {
	h := NewReviewHandler(&mockReviewService{})

	req := requestWithURLParam(http.MethodGet, "/api/v1/providers/abc/reviews", "id", "abc", nil)
	w := httptest.NewRecorder()

	h.ListByProvider()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestReviewHandler_ListByService_ServiceError(t *testing.T) {
	h := NewReviewHandler(&mockReviewService{
		listByServiceFn: func(ctx context.Context, serviceID int64, limit, offset int) ([]*model.Review, error) {
			return nil, errors.New("db failed")
		},
	})

	req := requestWithURLParam(http.MethodGet, "/api/v1/services/3/reviews", "id", "3", nil)
	w := httptest.NewRecorder()

	h.ListByService()(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func requestWithURLParam(method, target, key, value string, body *bytes.Buffer) *http.Request {
	var reader *bytes.Buffer
	if body == nil {
		reader = bytes.NewBuffer(nil)
	} else {
		reader = body
	}
	req := httptest.NewRequest(method, target, reader)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add(key, value)
	return req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, routeCtx))
}
