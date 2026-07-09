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

	"ors-be/internal/model"
	"ors-be/internal/service"
)

type mockServiceProviderService struct {
	createFn     func(ctx context.Context, userID int64, input service.ServiceProviderInput) (*model.ServiceProvider, error)
	getByIDFn    func(ctx context.Context, id int64) (*model.ServiceProvider, error)
	getMineFn    func(ctx context.Context, userID int64) (*model.ServiceProvider, error)
	updateMineFn func(ctx context.Context, userID int64, input service.ServiceProviderInput) (*model.ServiceProvider, error)
}

func (m *mockServiceProviderService) Create(ctx context.Context, userID int64, input service.ServiceProviderInput) (*model.ServiceProvider, error) {
	return m.createFn(ctx, userID, input)
}
func (m *mockServiceProviderService) GetByID(ctx context.Context, id int64) (*model.ServiceProvider, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockServiceProviderService) GetMine(ctx context.Context, userID int64) (*model.ServiceProvider, error) {
	return m.getMineFn(ctx, userID)
}
func (m *mockServiceProviderService) UpdateMine(ctx context.Context, userID int64, input service.ServiceProviderInput) (*model.ServiceProvider, error) {
	return m.updateMineFn(ctx, userID, input)
}

// ────── CreateMine ──────

func TestProviderHandler_CreateMine_Success(t *testing.T) {
	svc := &mockServiceProviderService{
		createFn: func(_ context.Context, userID int64, input service.ServiceProviderInput) (*model.ServiceProvider, error) {
			return &model.ServiceProvider{
				ID: 1, UserID: userID, BusinessName: input.BusinessName,
				CreatedAt: time.Now(), UpdatedAt: time.Now(),
			}, nil
		},
	}
	h := NewServiceProviderHandler(svc)

	body := `{"business_name":"舒心养生馆","description":"专业按摩","address":"北京市","phone":"13800138000","email":"info@test.com"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/providers/me", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	w := httptest.NewRecorder()

	h.CreateMine()(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusCreated)
	}
	var resp struct {
		Code int                   `json:"code"`
		Data model.ServiceProvider `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Code != 201 {
		t.Errorf("code = %d, want 201", resp.Code)
	}
	if resp.Data.BusinessName != "舒心养生馆" {
		t.Errorf("BusinessName = %s", resp.Data.BusinessName)
	}
}

func TestProviderHandler_CreateMine_Unauthorized(t *testing.T) {
	h := NewServiceProviderHandler(&mockServiceProviderService{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/providers/me", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateMine()(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestProviderHandler_CreateMine_InvalidJSON(t *testing.T) {
	h := NewServiceProviderHandler(&mockServiceProviderService{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/providers/me", strings.NewReader(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	w := httptest.NewRecorder()

	h.CreateMine()(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestProviderHandler_CreateMine_BusinessNameRequired(t *testing.T) {
	svc := &mockServiceProviderService{
		createFn: func(_ context.Context, _ int64, _ service.ServiceProviderInput) (*model.ServiceProvider, error) {
			return nil, service.ErrBusinessNameRequired
		},
	}
	h := NewServiceProviderHandler(svc)

	body := `{"business_name":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/providers/me", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	w := httptest.NewRecorder()

	h.CreateMine()(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestProviderHandler_CreateMine_InvalidEmail(t *testing.T) {
	svc := &mockServiceProviderService{
		createFn: func(_ context.Context, _ int64, _ service.ServiceProviderInput) (*model.ServiceProvider, error) {
			return nil, service.ErrInvalidEmail
		},
	}
	h := NewServiceProviderHandler(svc)

	body := `{"business_name":"test","email":"bad-email"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/providers/me", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	w := httptest.NewRecorder()

	h.CreateMine()(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestProviderHandler_CreateMine_AlreadyExists(t *testing.T) {
	svc := &mockServiceProviderService{
		createFn: func(_ context.Context, _ int64, _ service.ServiceProviderInput) (*model.ServiceProvider, error) {
			return nil, service.ErrProviderAlreadyExists
		},
	}
	h := NewServiceProviderHandler(svc)

	body := `{"business_name":"test"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/providers/me", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	w := httptest.NewRecorder()

	h.CreateMine()(w, req)
	if w.Code != http.StatusConflict {
		t.Errorf("status = %d, want %d", w.Code, http.StatusConflict)
	}
}

// ────── GetMine ──────

func TestProviderHandler_GetMine_Success(t *testing.T) {
	svc := &mockServiceProviderService{
		getMineFn: func(_ context.Context, userID int64) (*model.ServiceProvider, error) {
			return &model.ServiceProvider{ID: 1, UserID: userID, BusinessName: "我的店铺"}, nil
		},
	}
	h := NewServiceProviderHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/providers/me", nil)
	req = withProviderClaims(req, 2)
	w := httptest.NewRecorder()

	h.GetMine()(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestProviderHandler_GetMine_Unauthorized(t *testing.T) {
	h := NewServiceProviderHandler(&mockServiceProviderService{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/providers/me", nil)
	w := httptest.NewRecorder()

	h.GetMine()(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestProviderHandler_GetMine_NotFound(t *testing.T) {
	svc := &mockServiceProviderService{
		getMineFn: func(_ context.Context, _ int64) (*model.ServiceProvider, error) {
			return nil, service.ErrProviderNotFound
		},
	}
	h := NewServiceProviderHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/providers/me", nil)
	req = withProviderClaims(req, 999)
	w := httptest.NewRecorder()

	h.GetMine()(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// ────── UpdateMine ──────

func TestProviderHandler_UpdateMine_Success(t *testing.T) {
	svc := &mockServiceProviderService{
		updateMineFn: func(_ context.Context, _ int64, input service.ServiceProviderInput) (*model.ServiceProvider, error) {
			return &model.ServiceProvider{ID: 1, BusinessName: input.BusinessName, Address: input.Address}, nil
		},
	}
	h := NewServiceProviderHandler(svc)

	body := `{"business_name":"新店名","address":"新地址"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/providers/me", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	w := httptest.NewRecorder()

	h.UpdateMine()(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestProviderHandler_UpdateMine_Unauthorized(t *testing.T) {
	h := NewServiceProviderHandler(&mockServiceProviderService{})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/providers/me", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UpdateMine()(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestProviderHandler_UpdateMine_NotFound(t *testing.T) {
	svc := &mockServiceProviderService{
		updateMineFn: func(_ context.Context, _ int64, _ service.ServiceProviderInput) (*model.ServiceProvider, error) {
			return nil, service.ErrProviderNotFound
		},
	}
	h := NewServiceProviderHandler(svc)

	body := `{"business_name":"test"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/providers/me", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 999)
	w := httptest.NewRecorder()

	h.UpdateMine()(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestProviderHandler_UpdateMine_InvalidEmail(t *testing.T) {
	svc := &mockServiceProviderService{
		updateMineFn: func(_ context.Context, _ int64, _ service.ServiceProviderInput) (*model.ServiceProvider, error) {
			return nil, service.ErrInvalidEmail
		},
	}
	h := NewServiceProviderHandler(svc)

	body := `{"business_name":"test","email":"bad-email"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/providers/me", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	w := httptest.NewRecorder()

	h.UpdateMine()(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// ────── GetByID ──────

func TestProviderHandler_GetByID_Success(t *testing.T) {
	svc := &mockServiceProviderService{
		getByIDFn: func(_ context.Context, id int64) (*model.ServiceProvider, error) {
			return &model.ServiceProvider{ID: id, BusinessName: "舒心养生馆"}, nil
		},
	}
	h := NewServiceProviderHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/providers/1", nil)
	req = withChiURLParam(req, "id", "1")
	w := httptest.NewRecorder()

	h.GetByID()(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestProviderHandler_GetByID_InvalidID(t *testing.T) {
	h := NewServiceProviderHandler(&mockServiceProviderService{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/providers/abc", nil)
	req = withChiURLParam(req, "id", "abc")
	w := httptest.NewRecorder()

	h.GetByID()(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestProviderHandler_GetByID_NotFound(t *testing.T) {
	svc := &mockServiceProviderService{
		getByIDFn: func(_ context.Context, _ int64) (*model.ServiceProvider, error) {
			return nil, service.ErrProviderNotFound
		},
	}
	h := NewServiceProviderHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/providers/999", nil)
	req = withChiURLParam(req, "id", "999")
	w := httptest.NewRecorder()

	h.GetByID()(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// ────── Status code coverage ──────

func TestProviderHandler_StatusCodeCoverage(t *testing.T) {
	used := map[int]bool{}

	// 201
	svcOK := &mockServiceProviderService{
		createFn: func(_ context.Context, _ int64, _ service.ServiceProviderInput) (*model.ServiceProvider, error) {
			return &model.ServiceProvider{ID: 1, BusinessName: "t"}, nil
		},
		getMineFn: func(_ context.Context, _ int64) (*model.ServiceProvider, error) {
			return &model.ServiceProvider{ID: 1, BusinessName: "t"}, nil
		},
		getByIDFn: func(_ context.Context, _ int64) (*model.ServiceProvider, error) {
			return &model.ServiceProvider{ID: 1, BusinessName: "t"}, nil
		},
		updateMineFn: func(_ context.Context, _ int64, _ service.ServiceProviderInput) (*model.ServiceProvider, error) {
			return &model.ServiceProvider{ID: 1, BusinessName: "t"}, nil
		},
	}
	h := NewServiceProviderHandler(svcOK)

	// 201
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"business_name":"t"}`))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	w := httptest.NewRecorder()
	h.CreateMine()(w, req)
	used[w.Code] = true

	// 200 GetMine
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req = withProviderClaims(req, 2)
	w = httptest.NewRecorder()
	h.GetMine()(w, req)
	used[w.Code] = true

	// 200 UpdateMine
	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{"business_name":"t"}`))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	w = httptest.NewRecorder()
	h.UpdateMine()(w, req)
	used[w.Code] = true

	// 200 GetByID
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req = withChiURLParam(req, "id", "1")
	w = httptest.NewRecorder()
	h.GetByID()(w, req)
	used[w.Code] = true

	// 400
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	w = httptest.NewRecorder()
	h.CreateMine()(w, req)
	used[w.Code] = true

	// 401
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	h.CreateMine()(w, req)
	used[w.Code] = true

	// 404
	notFoundSvc := &mockServiceProviderService{
		getByIDFn: func(_ context.Context, _ int64) (*model.ServiceProvider, error) {
			return nil, service.ErrProviderNotFound
		},
	}
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req = withChiURLParam(req, "id", "999")
	w = httptest.NewRecorder()
	NewServiceProviderHandler(notFoundSvc).GetByID()(w, req)
	used[w.Code] = true

	// 409
	conflictSvc := &mockServiceProviderService{
		createFn: func(_ context.Context, _ int64, _ service.ServiceProviderInput) (*model.ServiceProvider, error) {
			return nil, service.ErrProviderAlreadyExists
		},
	}
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"business_name":"t"}`))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	w = httptest.NewRecorder()
	NewServiceProviderHandler(conflictSvc).CreateMine()(w, req)
	used[w.Code] = true

	// 500
	errSvc := &mockServiceProviderService{
		getMineFn: func(_ context.Context, _ int64) (*model.ServiceProvider, error) {
			return nil, errors.New("internal error")
		},
	}
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req = withProviderClaims(req, 2)
	w = httptest.NewRecorder()
	NewServiceProviderHandler(errSvc).GetMine()(w, req)
	used[w.Code] = true

	expected := []int{200, 201, 400, 401, 404, 409, 500}
	for _, code := range expected {
		if !used[code] {
			t.Errorf("status code %d is not covered", code)
		}
	}
}
