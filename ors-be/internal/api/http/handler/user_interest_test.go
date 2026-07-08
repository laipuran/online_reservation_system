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

type mockUserInterestService struct {
	listFn    func(ctx context.Context, userID int64) ([]*model.Tag, error)
	replaceFn func(ctx context.Context, userID int64, input service.UserInterestsInput) ([]*model.Tag, error)
}

func (m *mockUserInterestService) List(ctx context.Context, userID int64) ([]*model.Tag, error) {
	return m.listFn(ctx, userID)
}
func (m *mockUserInterestService) Replace(ctx context.Context, userID int64, input service.UserInterestsInput) ([]*model.Tag, error) {
	return m.replaceFn(ctx, userID, input)
}

// ────── ListMine ──────

func TestUserInterestHandler_ListMine_Success(t *testing.T) {
	svc := &mockUserInterestService{
		listFn: func(_ context.Context, userID int64) ([]*model.Tag, error) {
			return []*model.Tag{
				{ID: 1, Name: "放松", CreatedAt: time.Now()},
				{ID: 2, Name: "塑形", CreatedAt: time.Now()},
			}, nil
		},
	}
	h := NewUserInterestHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me/interests", nil)
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.ListMine()(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp struct {
		Code int          `json:"code"`
		Data []model.Tag  `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if len(resp.Data) != 2 {
		t.Errorf("len = %d, want 2", len(resp.Data))
	}
}

func TestUserInterestHandler_ListMine_Unauthorized(t *testing.T) {
	h := NewUserInterestHandler(&mockUserInterestService{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me/interests", nil)
	w := httptest.NewRecorder()

	h.ListMine()(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestUserInterestHandler_ListMine_Empty(t *testing.T) {
	svc := &mockUserInterestService{
		listFn: func(_ context.Context, _ int64) ([]*model.Tag, error) {
			return []*model.Tag{}, nil
		},
	}
	h := NewUserInterestHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me/interests", nil)
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.ListMine()(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestUserInterestHandler_ListMine_Error(t *testing.T) {
	svc := &mockUserInterestService{
		listFn: func(_ context.Context, _ int64) ([]*model.Tag, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewUserInterestHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me/interests", nil)
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.ListMine()(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// ────── ReplaceMine ──────

func TestUserInterestHandler_ReplaceMine_Success(t *testing.T) {
	svc := &mockUserInterestService{
		replaceFn: func(_ context.Context, userID int64, input service.UserInterestsInput) ([]*model.Tag, error) {
			return []*model.Tag{
				{ID: 1, Name: "医疗"},
				{ID: 3, Name: "健身"},
			}, nil
		},
	}
	h := NewUserInterestHandler(svc)

	body := `{"tag_ids":[1,3]}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me/interests", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.ReplaceMine()(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestUserInterestHandler_ReplaceMine_Unauthorized(t *testing.T) {
	h := NewUserInterestHandler(&mockUserInterestService{})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me/interests", strings.NewReader(`{"tag_ids":[1]}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.ReplaceMine()(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestUserInterestHandler_ReplaceMine_InvalidJSON(t *testing.T) {
	h := NewUserInterestHandler(&mockUserInterestService{})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me/interests", strings.NewReader(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.ReplaceMine()(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestUserInterestHandler_ReplaceMine_InvalidTag(t *testing.T) {
	svc := &mockUserInterestService{
		replaceFn: func(_ context.Context, _ int64, _ service.UserInterestsInput) ([]*model.Tag, error) {
			return nil, service.ErrUserInterestInvalidTag
		},
	}
	h := NewUserInterestHandler(svc)

	body := `{"tag_ids":[999]}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me/interests", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.ReplaceMine()(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestUserInterestHandler_ReplaceMine_TagNotFound(t *testing.T) {
	svc := &mockUserInterestService{
		replaceFn: func(_ context.Context, _ int64, _ service.UserInterestsInput) ([]*model.Tag, error) {
			return nil, service.ErrTagNotFound
		},
	}
	h := NewUserInterestHandler(svc)

	body := `{"tag_ids":[999]}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me/interests", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.ReplaceMine()(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// ────── Status code coverage ──────

func TestUserInterestHandler_StatusCodeCoverage(t *testing.T) {
	used := map[int]bool{}

	svcOK := &mockUserInterestService{
		listFn: func(_ context.Context, _ int64) ([]*model.Tag, error) {
			return []*model.Tag{}, nil
		},
		replaceFn: func(_ context.Context, _ int64, _ service.UserInterestsInput) ([]*model.Tag, error) {
			return []*model.Tag{{ID: 1, Name: "t"}}, nil
		},
	}
	h := NewUserInterestHandler(svcOK)

	// 200 List
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()
	h.ListMine()(w, req)
	used[w.Code] = true

	// 200 Replace
	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{"tag_ids":[1]}`))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w = httptest.NewRecorder()
	h.ReplaceMine()(w, req)
	used[w.Code] = true

	// 400 Invalid JSON
	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w = httptest.NewRecorder()
	h.ReplaceMine()(w, req)
	used[w.Code] = true

	// 401
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	w = httptest.NewRecorder()
	h.ListMine()(w, req)
	used[w.Code] = true

	// 404
	notFoundSvc := &mockUserInterestService{
		replaceFn: func(_ context.Context, _ int64, _ service.UserInterestsInput) ([]*model.Tag, error) {
			return nil, service.ErrTagNotFound
		},
	}
	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{"tag_ids":[999]}`))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w = httptest.NewRecorder()
	NewUserInterestHandler(notFoundSvc).ReplaceMine()(w, req)
	used[w.Code] = true

	// 500
	errSvc := &mockUserInterestService{
		listFn: func(_ context.Context, _ int64) ([]*model.Tag, error) {
			return nil, errors.New("db error")
		},
	}
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req = withCustomerClaims(req, 1)
	w = httptest.NewRecorder()
	NewUserInterestHandler(errSvc).ListMine()(w, req)
	used[w.Code] = true

	expected := []int{200, 400, 401, 404, 500}
	for _, code := range expected {
		if !used[code] {
			t.Errorf("status code %d is not covered", code)
		}
	}
}
