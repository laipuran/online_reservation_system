package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"

	"ors-be/internal/api/http/middleware"
	"ors-be/internal/auth"
	"ors-be/internal/model"
	"ors-be/internal/repository"
	"ors-be/internal/service"
)

type mockReservationService struct {
	createFn             func(ctx context.Context, userID int64, input service.ReservationInput) (*model.ReservationView, error)
	getMineFn            func(ctx context.Context, userID, id int64) (*model.Reservation, error)
	listMineFn           func(ctx context.Context, userID int64, status string, page, pageSize int) (*service.ReservationListResult, error)
	cancelMineFn         func(ctx context.Context, userID, id int64) (*model.Reservation, error)
	listForProviderFn    func(ctx context.Context, userID int64, status string, page, pageSize int) (*service.ReservationListResult, error)
	confirmForProviderFn func(ctx context.Context, userID, id int64) (*model.Reservation, error)
	rejectForProviderFn  func(ctx context.Context, userID, id int64) (*model.Reservation, error)
	completeDueFn        func(ctx context.Context, now time.Time) (int64, error)
}

func (m *mockReservationService) Create(ctx context.Context, userID int64, input service.ReservationInput) (*model.ReservationView, error) {
	return m.createFn(ctx, userID, input)
}

func (m *mockReservationService) GetMine(ctx context.Context, userID, id int64) (*model.Reservation, error) {
	return m.getMineFn(ctx, userID, id)
}

func (m *mockReservationService) ListMine(ctx context.Context, userID int64, status string, page, pageSize int) (*service.ReservationListResult, error) {
	return m.listMineFn(ctx, userID, status, page, pageSize)
}

func (m *mockReservationService) CancelMine(ctx context.Context, userID, id int64) (*model.Reservation, error) {
	return m.cancelMineFn(ctx, userID, id)
}

func (m *mockReservationService) ListForProvider(ctx context.Context, userID int64, status string, page, pageSize int) (*service.ReservationListResult, error) {
	return m.listForProviderFn(ctx, userID, status, page, pageSize)
}

func (m *mockReservationService) ConfirmForProvider(ctx context.Context, userID, id int64) (*model.Reservation, error) {
	return m.confirmForProviderFn(ctx, userID, id)
}

func (m *mockReservationService) RejectForProvider(ctx context.Context, userID, id int64) (*model.Reservation, error) {
	return m.rejectForProviderFn(ctx, userID, id)
}

func (m *mockReservationService) CompleteDue(ctx context.Context, now time.Time) (int64, error) {
	return m.completeDueFn(ctx, now)
}

// withCustomerClaims injects customer JWT claims into the request context.
func withCustomerClaims(r *http.Request, userID int64) *http.Request {
	ctx := context.WithValue(r.Context(), middleware.UserCtxKey, &auth.Claims{
		UserID: userID,
		Role:   "customer",
	})
	return r.WithContext(ctx)
}

// withProviderClaims injects provider JWT claims into the request context.
func withProviderClaims(r *http.Request, userID int64) *http.Request {
	ctx := context.WithValue(r.Context(), middleware.UserCtxKey, &auth.Claims{
		UserID: userID,
		Role:   "provider",
	})
	return r.WithContext(ctx)
}

// withChiURLParam injects a chi URL parameter into the request context.
func withChiURLParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	return r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))
}

// ────── Create ──────

func TestReservationHandler_Create_Success(t *testing.T) {
	startTime := time.Date(2026, 7, 10, 14, 0, 0, 0, time.UTC)
	endTime := startTime.Add(time.Hour)

	svc := &mockReservationService{
		createFn: func(_ context.Context, userID int64, input service.ReservationInput) (*model.ReservationView, error) {
			if input.ServiceID != 1 || !input.StartTime.Equal(startTime) || input.Note != "请准备热水" {
				return nil, service.ErrReservationInvalidInput
			}
			return &model.ReservationView{
				ID: 1001,
				Service: model.ReservationServiceSummary{
					ID:    1,
					Title: "肩颈按摩 60 分钟",
					Provider: model.ReservationServiceProviderSummary{
						ID:           1,
						BusinessName: "舒心养生馆",
					},
				},
				StartTime: startTime,
				EndTime:   endTime,
				Status:    "pending",
				Note:      "请准备热水",
				CreatedAt: startTime,
			}, nil
		},
	}

	h := NewReservationHandler(svc)

	body := `{"service_id":1,"start_time":"2026-07-10T14:00:00Z","note":"请准备热水"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reservations", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.Create()(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusCreated)
	}

	var resp struct {
		Code    int              `json:"code"`
		Message string           `json:"message"`
		Data    *json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v, body: %s", err, w.Body.String())
	}
	if resp.Code != 201 {
		t.Errorf("code = %d, want 201", resp.Code)
	}
	if resp.Message != "created" {
		t.Errorf("message = %s, want created", resp.Message)
	}

	var view model.ReservationView
	if err := json.Unmarshal(*resp.Data, &view); err != nil {
		t.Fatalf("decode data: %v", err)
	}
	if view.ID != 1001 {
		t.Errorf("ID = %d, want 1001", view.ID)
	}
	if view.Service.ID != 1 {
		t.Errorf("Service.ID = %d, want 1", view.Service.ID)
	}
	if view.Service.Provider.BusinessName != "舒心养生馆" {
		t.Errorf("Provider.BusinessName = %s", view.Service.Provider.BusinessName)
	}
	if view.Status != "pending" {
		t.Errorf("Status = %s, want pending", view.Status)
	}
	if !view.EndTime.Equal(endTime) {
		t.Errorf("EndTime = %s, want %s", view.EndTime, endTime)
	}
}

func TestReservationHandler_Create_Unauthorized(t *testing.T) {
	svc := &mockReservationService{
		createFn: func(_ context.Context, _ int64, _ service.ReservationInput) (*model.ReservationView, error) {
			t.Fatal("Create should not be called without claims")
			return nil, nil
		},
	}
	h := NewReservationHandler(svc)

	body := `{"service_id":1,"start_time":"2026-07-10T14:00:00Z"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reservations", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	// no claims injected
	w := httptest.NewRecorder()

	h.Create()(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestReservationHandler_Create_InvalidJSON(t *testing.T) {
	svc := &mockReservationService{
		createFn: func(_ context.Context, _ int64, _ service.ReservationInput) (*model.ReservationView, error) {
			t.Fatal("Create should not be called with invalid JSON")
			return nil, nil
		},
	}
	h := NewReservationHandler(svc)

	body := `{bad json`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reservations", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.Create()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestReservationHandler_Create_InvalidInput(t *testing.T) {
	svc := &mockReservationService{
		createFn: func(_ context.Context, _ int64, _ service.ReservationInput) (*model.ReservationView, error) {
			return nil, service.ErrReservationInvalidInput
		},
	}
	h := NewReservationHandler(svc)

	body := `{"service_id":0,"start_time":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reservations", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.Create()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestReservationHandler_Create_TimeConflict(t *testing.T) {
	svc := &mockReservationService{
		createFn: func(_ context.Context, _ int64, _ service.ReservationInput) (*model.ReservationView, error) {
			return nil, repository.ErrReservationTimeConflict
		},
	}
	h := NewReservationHandler(svc)

	body := `{"service_id":1,"start_time":"2026-07-10T14:00:00Z"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reservations", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.Create()(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("status = %d, want %d", w.Code, http.StatusConflict)
	}
}

func TestReservationHandler_Create_ServiceNotFound(t *testing.T) {
	svc := &mockReservationService{
		createFn: func(_ context.Context, _ int64, _ service.ReservationInput) (*model.ReservationView, error) {
			return nil, service.ErrServiceNotFound
		},
	}
	h := NewReservationHandler(svc)

	body := `{"service_id":999,"start_time":"2026-07-10T14:00:00Z"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reservations", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.Create()(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestReservationHandler_Create_UnexpectedError(t *testing.T) {
	svc := &mockReservationService{
		createFn: func(_ context.Context, _ int64, _ service.ReservationInput) (*model.ReservationView, error) {
			return nil, errors.New("db connection lost")
		},
	}
	h := NewReservationHandler(svc)

	body := `{"service_id":1,"start_time":"2026-07-10T14:00:00Z"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/reservations", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.Create()(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// ────── ListMine ──────

func TestReservationHandler_ListMine_Success(t *testing.T) {
	svc := &mockReservationService{
		listMineFn: func(_ context.Context, userID int64, status string, page, pageSize int) (*service.ReservationListResult, error) {
			return &service.ReservationListResult{
				Items: []*model.Reservation{
					{ID: 10, UserID: 1, ServiceID: 2, Status: "pending"},
				},
				Page:     1,
				PageSize: 20,
			}, nil
		},
	}
	h := NewReservationHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reservations?status=pending&page=1&page_size=20", nil)
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.ListMine()(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp struct {
		Code    int              `json:"code"`
		Message string           `json:"message"`
		Data    *json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Code != 200 {
		t.Errorf("code = %d, want 200", resp.Code)
	}

	var result service.ReservationListResult
	if err := json.Unmarshal(*resp.Data, &result); err != nil {
		t.Fatalf("decode data: %v", err)
	}
	if len(result.Items) != 1 {
		t.Fatalf("items len = %d, want 1", len(result.Items))
	}
	if result.Items[0].ID != 10 {
		t.Errorf("items[0].ID = %d, want 10", result.Items[0].ID)
	}
	if result.Page != 1 {
		t.Errorf("page = %d, want 1", result.Page)
	}
}

func TestReservationHandler_ListMine_Unauthorized(t *testing.T) {
	h := NewReservationHandler(&mockReservationService{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reservations", nil)
	// no claims
	w := httptest.NewRecorder()

	h.ListMine()(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestReservationHandler_ListMine_InvalidPage(t *testing.T) {
	h := NewReservationHandler(&mockReservationService{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reservations?page=abc", nil)
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.ListMine()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestReservationHandler_ListMine_InvalidStatus(t *testing.T) {
	svc := &mockReservationService{
		listMineFn: func(_ context.Context, _ int64, status string, _, _ int) (*service.ReservationListResult, error) {
			return nil, service.ErrReservationInvalidStatus
		},
	}
	h := NewReservationHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reservations?status=invalid_status", nil)
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.ListMine()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// ────── GetMine ──────

func TestReservationHandler_GetMine_Success(t *testing.T) {
	startTime := time.Date(2026, 7, 10, 14, 0, 0, 0, time.UTC)
	svc := &mockReservationService{
		getMineFn: func(_ context.Context, userID, id int64) (*model.Reservation, error) {
			return &model.Reservation{
				ID:        10,
				UserID:    userID,
				ServiceID: 2,
				StartTime: startTime,
				EndTime:   startTime.Add(time.Hour),
				Status:    "pending",
			}, nil
		},
	}
	h := NewReservationHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reservations/10", nil)
	req = withCustomerClaims(req, 1)
	req = withChiURLParam(req, "id", "10")
	w := httptest.NewRecorder()

	h.GetMine()(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestReservationHandler_GetMine_Unauthorized(t *testing.T) {
	h := NewReservationHandler(&mockReservationService{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reservations/10", nil)
	w := httptest.NewRecorder()

	h.GetMine()(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestReservationHandler_GetMine_BadID(t *testing.T) {
	h := NewReservationHandler(&mockReservationService{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reservations/abc", nil)
	req = withCustomerClaims(req, 1)
	req = withChiURLParam(req, "id", "abc")
	w := httptest.NewRecorder()

	h.GetMine()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestReservationHandler_GetMine_NotFound(t *testing.T) {
	svc := &mockReservationService{
		getMineFn: func(_ context.Context, _, _ int64) (*model.Reservation, error) {
			return nil, service.ErrReservationNotFound
		},
	}
	h := NewReservationHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/reservations/999", nil)
	req = withCustomerClaims(req, 1)
	req = withChiURLParam(req, "id", "999")
	w := httptest.NewRecorder()

	h.GetMine()(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// ────── CancelMine ──────

func TestReservationHandler_CancelMine_Success(t *testing.T) {
	svc := &mockReservationService{
		cancelMineFn: func(_ context.Context, _, _ int64) (*model.Reservation, error) {
			return &model.Reservation{ID: 10, Status: "cancelled"}, nil
		},
	}
	h := NewReservationHandler(svc)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/reservations/10/cancel", nil)
	req = withCustomerClaims(req, 1)
	req = withChiURLParam(req, "id", "10")
	w := httptest.NewRecorder()

	h.CancelMine()(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestReservationHandler_CancelMine_Unauthorized(t *testing.T) {
	h := NewReservationHandler(&mockReservationService{})

	req := httptest.NewRequest(http.MethodPut, "/api/v1/reservations/10/cancel", nil)
	w := httptest.NewRecorder()

	h.CancelMine()(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestReservationHandler_CancelMine_CannotCancel(t *testing.T) {
	svc := &mockReservationService{
		cancelMineFn: func(_ context.Context, _, _ int64) (*model.Reservation, error) {
			return nil, service.ErrReservationCannotCancel
		},
	}
	h := NewReservationHandler(svc)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/reservations/10/cancel", nil)
	req = withCustomerClaims(req, 1)
	req = withChiURLParam(req, "id", "10")
	w := httptest.NewRecorder()

	h.CancelMine()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestReservationHandler_CancelMine_NotFound(t *testing.T) {
	svc := &mockReservationService{
		cancelMineFn: func(_ context.Context, _, _ int64) (*model.Reservation, error) {
			return nil, service.ErrReservationNotFound
		},
	}
	h := NewReservationHandler(svc)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/reservations/999/cancel", nil)
	req = withCustomerClaims(req, 1)
	req = withChiURLParam(req, "id", "999")
	w := httptest.NewRecorder()

	h.CancelMine()(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// ────── ListForProvider ──────

func TestReservationHandler_ListForProvider_Success(t *testing.T) {
	svc := &mockReservationService{
		listForProviderFn: func(_ context.Context, userID int64, status string, page, pageSize int) (*service.ReservationListResult, error) {
			return &service.ReservationListResult{
				Items: []*model.Reservation{
					{ID: 20, UserID: 1, ServiceID: 2, Status: "confirmed"},
				},
				Page:     1,
				PageSize: 20,
			}, nil
		},
	}
	h := NewReservationHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/provider/reservations", nil)
	req = withProviderClaims(req, 20)
	w := httptest.NewRecorder()

	h.ListForProvider()(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestReservationHandler_ListForProvider_Unauthorized(t *testing.T) {
	h := NewReservationHandler(&mockReservationService{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/provider/reservations", nil)
	w := httptest.NewRecorder()

	h.ListForProvider()(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestReservationHandler_ListForProvider_ProviderNotFound(t *testing.T) {
	svc := &mockReservationService{
		listForProviderFn: func(_ context.Context, _ int64, _ string, _, _ int) (*service.ReservationListResult, error) {
			return nil, service.ErrProviderNotFound
		},
	}
	h := NewReservationHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/provider/reservations", nil)
	req = withProviderClaims(req, 999)
	w := httptest.NewRecorder()

	h.ListForProvider()(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// ────── ConfirmForProvider ──────

func TestReservationHandler_ConfirmForProvider_Success(t *testing.T) {
	svc := &mockReservationService{
		confirmForProviderFn: func(_ context.Context, _, _ int64) (*model.Reservation, error) {
			return &model.Reservation{ID: 10, Status: "confirmed"}, nil
		},
	}
	h := NewReservationHandler(svc)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/provider/reservations/10/confirm", nil)
	req = withProviderClaims(req, 20)
	req = withChiURLParam(req, "id", "10")
	w := httptest.NewRecorder()

	h.ConfirmForProvider()(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestReservationHandler_ConfirmForProvider_Unauthorized(t *testing.T) {
	h := NewReservationHandler(&mockReservationService{})

	req := httptest.NewRequest(http.MethodPut, "/api/v1/provider/reservations/10/confirm", nil)
	w := httptest.NewRecorder()

	h.ConfirmForProvider()(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestReservationHandler_ConfirmForProvider_CannotConfirm(t *testing.T) {
	svc := &mockReservationService{
		confirmForProviderFn: func(_ context.Context, _, _ int64) (*model.Reservation, error) {
			return nil, service.ErrReservationCannotConfirm
		},
	}
	h := NewReservationHandler(svc)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/provider/reservations/10/confirm", nil)
	req = withProviderClaims(req, 20)
	req = withChiURLParam(req, "id", "10")
	w := httptest.NewRecorder()

	h.ConfirmForProvider()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestReservationHandler_ConfirmForProvider_NotFound(t *testing.T) {
	svc := &mockReservationService{
		confirmForProviderFn: func(_ context.Context, _, _ int64) (*model.Reservation, error) {
			return nil, service.ErrReservationNotFound
		},
	}
	h := NewReservationHandler(svc)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/provider/reservations/999/confirm", nil)
	req = withProviderClaims(req, 20)
	req = withChiURLParam(req, "id", "999")
	w := httptest.NewRecorder()

	h.ConfirmForProvider()(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// ────── RejectForProvider ──────

func TestReservationHandler_RejectForProvider_Success(t *testing.T) {
	svc := &mockReservationService{
		rejectForProviderFn: func(_ context.Context, _, _ int64) (*model.Reservation, error) {
			return &model.Reservation{ID: 10, Status: "rejected"}, nil
		},
	}
	h := NewReservationHandler(svc)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/provider/reservations/10/reject", nil)
	req = withProviderClaims(req, 20)
	req = withChiURLParam(req, "id", "10")
	w := httptest.NewRecorder()

	h.RejectForProvider()(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestReservationHandler_RejectForProvider_Unauthorized(t *testing.T) {
	h := NewReservationHandler(&mockReservationService{})

	req := httptest.NewRequest(http.MethodPut, "/api/v1/provider/reservations/10/reject", nil)
	w := httptest.NewRecorder()

	h.RejectForProvider()(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestReservationHandler_RejectForProvider_CannotReject(t *testing.T) {
	svc := &mockReservationService{
		rejectForProviderFn: func(_ context.Context, _, _ int64) (*model.Reservation, error) {
			return nil, service.ErrReservationCannotReject
		},
	}
	h := NewReservationHandler(svc)

	req := httptest.NewRequest(http.MethodPut, "/api/v1/provider/reservations/10/reject", nil)
	req = withProviderClaims(req, 20)
	req = withChiURLParam(req, "id", "10")
	w := httptest.NewRecorder()

	h.RejectForProvider()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// ────── Status code coverage ──────

func TestReservationHandler_StatusCodeCoverage(t *testing.T) {
	startTime := time.Date(2026, 7, 10, 14, 0, 0, 0, time.UTC)

	// 200: GetMine success
	svc := &mockReservationService{
		getMineFn: func(_ context.Context, _, _ int64) (*model.Reservation, error) {
			return &model.Reservation{ID: 1, Status: "pending", StartTime: startTime, EndTime: startTime.Add(time.Hour)}, nil
		},
		listMineFn: func(_ context.Context, _ int64, _ string, _, _ int) (*service.ReservationListResult, error) {
			return &service.ReservationListResult{Items: []*model.Reservation{}, Page: 1, PageSize: 20}, nil
		},
		cancelMineFn: func(_ context.Context, _, _ int64) (*model.Reservation, error) {
			return &model.Reservation{ID: 1, Status: "cancelled"}, nil
		},
		listForProviderFn: func(_ context.Context, _ int64, _ string, _, _ int) (*service.ReservationListResult, error) {
			return &service.ReservationListResult{Items: []*model.Reservation{}, Page: 1, PageSize: 20}, nil
		},
		confirmForProviderFn: func(_ context.Context, _, _ int64) (*model.Reservation, error) {
			return &model.Reservation{ID: 1, Status: "confirmed"}, nil
		},
		rejectForProviderFn: func(_ context.Context, _, _ int64) (*model.Reservation, error) {
			return &model.Reservation{ID: 1, Status: "rejected"}, nil
		},
	}
	h := NewReservationHandler(svc)
	used := map[int]bool{}

	// 200 - GetMine
	req := httptest.NewRequest(http.MethodGet, "/api/v1/reservations/1", nil)
	req = withCustomerClaims(req, 1)
	req = withChiURLParam(req, "id", "1")
	w := httptest.NewRecorder()
	h.GetMine()(w, req)
	used[w.Code] = true

	// 200 - ListMine
	req = httptest.NewRequest(http.MethodGet, "/api/v1/reservations", nil)
	req = withCustomerClaims(req, 1)
	w = httptest.NewRecorder()
	h.ListMine()(w, req)
	used[w.Code] = true

	// 200 - CancelMine
	req = httptest.NewRequest(http.MethodPut, "/api/v1/reservations/1/cancel", nil)
	req = withCustomerClaims(req, 1)
	req = withChiURLParam(req, "id", "1")
	w = httptest.NewRecorder()
	h.CancelMine()(w, req)
	used[w.Code] = true

	// 200 - ListForProvider
	req = httptest.NewRequest(http.MethodGet, "/api/v1/provider/reservations", nil)
	req = withProviderClaims(req, 20)
	w = httptest.NewRecorder()
	h.ListForProvider()(w, req)
	used[w.Code] = true

	// 200 - ConfirmForProvider
	req = httptest.NewRequest(http.MethodPut, "/api/v1/provider/reservations/1/confirm", nil)
	req = withProviderClaims(req, 20)
	req = withChiURLParam(req, "id", "1")
	w = httptest.NewRecorder()
	h.ConfirmForProvider()(w, req)
	used[w.Code] = true

	// 200 - RejectForProvider
	req = httptest.NewRequest(http.MethodPut, "/api/v1/provider/reservations/1/reject", nil)
	req = withProviderClaims(req, 20)
	req = withChiURLParam(req, "id", "1")
	w = httptest.NewRecorder()
	h.RejectForProvider()(w, req)
	used[w.Code] = true

	// 201 - Create
	createSvc := &mockReservationService{
		createFn: func(_ context.Context, _ int64, _ service.ReservationInput) (*model.ReservationView, error) {
			return &model.ReservationView{
				ID: 1,
				Service: model.ReservationServiceSummary{
					ID:    1,
					Title: "Test",
					Provider: model.ReservationServiceProviderSummary{
						ID: 1, BusinessName: "Test",
					},
				},
				StartTime: startTime,
				EndTime:   startTime.Add(time.Hour),
				Status:    "pending",
				CreatedAt: startTime,
			}, nil
		},
	}
	body := `{"service_id":1,"start_time":"2026-07-10T14:00:00Z"}`
	req = httptest.NewRequest(http.MethodPost, "/api/v1/reservations", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w = httptest.NewRecorder()
	NewReservationHandler(createSvc).Create()(w, req)
	used[w.Code] = true

	// 400 - Bad JSON
	req = httptest.NewRequest(http.MethodPost, "/api/v1/reservations", strings.NewReader(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w = httptest.NewRecorder()
	h.Create()(w, req)
	used[w.Code] = true

	// 401 - Missing claims
	req = httptest.NewRequest(http.MethodGet, "/api/v1/reservations", nil)
	w = httptest.NewRecorder()
	h.ListMine()(w, req)
	used[w.Code] = true

	// 404 - Not found
	notFoundSvc := &mockReservationService{
		getMineFn: func(_ context.Context, _, _ int64) (*model.Reservation, error) {
			return nil, service.ErrReservationNotFound
		},
	}
	req = httptest.NewRequest(http.MethodGet, "/api/v1/reservations/999", nil)
	req = withCustomerClaims(req, 1)
	req = withChiURLParam(req, "id", "999")
	w = httptest.NewRecorder()
	NewReservationHandler(notFoundSvc).GetMine()(w, req)
	used[w.Code] = true

	// 409 - Time conflict
	conflictSvc := &mockReservationService{
		createFn: func(_ context.Context, _ int64, _ service.ReservationInput) (*model.ReservationView, error) {
			return nil, repository.ErrReservationTimeConflict
		},
	}
	body = `{"service_id":1,"start_time":"2026-07-10T14:00:00Z"}`
	req = httptest.NewRequest(http.MethodPost, "/api/v1/reservations", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w = httptest.NewRecorder()
	NewReservationHandler(conflictSvc).Create()(w, req)
	used[w.Code] = true

	// 500 - Unexpected error
	errSvc := &mockReservationService{
		listMineFn: func(_ context.Context, _ int64, _ string, _, _ int) (*service.ReservationListResult, error) {
			return nil, errors.New("internal error")
		},
	}
	req = httptest.NewRequest(http.MethodGet, "/api/v1/reservations", nil)
	req = withCustomerClaims(req, 1)
	w = httptest.NewRecorder()
	NewReservationHandler(errSvc).ListMine()(w, req)
	used[w.Code] = true

	expected := []int{200, 201, 400, 401, 404, 409, 500}
	for _, code := range expected {
		if !used[code] {
			t.Errorf("status code %d is not covered by any test case", code)
		}
	}
}
