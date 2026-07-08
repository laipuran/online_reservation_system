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

type mockTagService struct {
	createFn  func(ctx context.Context, input service.TagInput) (*model.Tag, error)
	getByIDFn func(ctx context.Context, id int64) (*model.Tag, error)
	listFn    func(ctx context.Context) ([]*model.Tag, error)
}

func (m *mockTagService) Create(ctx context.Context, input service.TagInput) (*model.Tag, error) {
	return m.createFn(ctx, input)
}
func (m *mockTagService) GetByID(ctx context.Context, id int64) (*model.Tag, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockTagService) List(ctx context.Context) ([]*model.Tag, error) {
	return m.listFn(ctx)
}

// ────── Create ──────

func TestTagHandler_Create_Success(t *testing.T) {
	svc := &mockTagService{
		createFn: func(_ context.Context, input service.TagInput) (*model.Tag, error) {
			return &model.Tag{ID: 1, Name: input.Name, CreatedAt: time.Now()}, nil
		},
	}
	h := NewTagHandler(svc)

	body := `{"name":"放松"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tags", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create()(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusCreated)
	}
}

func TestTagHandler_Create_InvalidJSON(t *testing.T) {
	h := NewTagHandler(&mockTagService{})

	req := httptest.NewRequest(http.MethodPost, "/api/v1/tags", strings.NewReader(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create()(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestTagHandler_Create_NameRequired(t *testing.T) {
	svc := &mockTagService{
		createFn: func(_ context.Context, _ service.TagInput) (*model.Tag, error) {
			return nil, service.ErrTagNameRequired
		},
	}
	h := NewTagHandler(svc)

	body := `{"name":""}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tags", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create()(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestTagHandler_Create_NameTooLong(t *testing.T) {
	svc := &mockTagService{
		createFn: func(_ context.Context, _ service.TagInput) (*model.Tag, error) {
			return nil, service.ErrTagNameTooLong
		},
	}
	h := NewTagHandler(svc)

	body := `{"name":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tags", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create()(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestTagHandler_Create_AlreadyExists(t *testing.T) {
	svc := &mockTagService{
		createFn: func(_ context.Context, _ service.TagInput) (*model.Tag, error) {
			return nil, service.ErrTagAlreadyExists
		},
	}
	h := NewTagHandler(svc)

	body := `{"name":"重复标签"}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/tags", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create()(w, req)
	if w.Code != http.StatusConflict {
		t.Errorf("status = %d, want %d", w.Code, http.StatusConflict)
	}
}

// ────── GetByID ──────

func TestTagHandler_GetByID_Success(t *testing.T) {
	svc := &mockTagService{
		getByIDFn: func(_ context.Context, id int64) (*model.Tag, error) {
			return &model.Tag{ID: id, Name: "医疗", CreatedAt: time.Now()}, nil
		},
	}
	h := NewTagHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tags/1", nil)
	req = withChiURLParam(req, "id", "1")
	w := httptest.NewRecorder()

	h.GetByID()(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestTagHandler_GetByID_InvalidID(t *testing.T) {
	h := NewTagHandler(&mockTagService{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tags/abc", nil)
	req = withChiURLParam(req, "id", "abc")
	w := httptest.NewRecorder()

	h.GetByID()(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestTagHandler_GetByID_NotFound(t *testing.T) {
	svc := &mockTagService{
		getByIDFn: func(_ context.Context, _ int64) (*model.Tag, error) {
			return nil, service.ErrTagNotFound
		},
	}
	h := NewTagHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tags/999", nil)
	req = withChiURLParam(req, "id", "999")
	w := httptest.NewRecorder()

	h.GetByID()(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// ────── List ──────

func TestTagHandler_List_Success(t *testing.T) {
	svc := &mockTagService{
		listFn: func(_ context.Context) ([]*model.Tag, error) {
			return []*model.Tag{
				{ID: 1, Name: "放松"},
				{ID: 2, Name: "塑形"},
			}, nil
		},
	}
	h := NewTagHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tags", nil)
	w := httptest.NewRecorder()

	h.List()(w, req)

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

func TestTagHandler_List_Empty(t *testing.T) {
	svc := &mockTagService{
		listFn: func(_ context.Context) ([]*model.Tag, error) {
			return []*model.Tag{}, nil
		},
	}
	h := NewTagHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tags", nil)
	w := httptest.NewRecorder()

	h.List()(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestTagHandler_List_Error(t *testing.T) {
	svc := &mockTagService{
		listFn: func(_ context.Context) ([]*model.Tag, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewTagHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/tags", nil)
	w := httptest.NewRecorder()

	h.List()(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// ────── Status code coverage ──────

func TestTagHandler_StatusCodeCoverage(t *testing.T) {
	used := map[int]bool{}

	// 200, 201
	svcOK := &mockTagService{
		createFn: func(_ context.Context, _ service.TagInput) (*model.Tag, error) {
			return &model.Tag{ID: 1, Name: "t"}, nil
		},
		getByIDFn: func(_ context.Context, _ int64) (*model.Tag, error) {
			return &model.Tag{ID: 1, Name: "t"}, nil
		},
		listFn: func(_ context.Context) ([]*model.Tag, error) {
			return []*model.Tag{}, nil
		},
	}
	h := NewTagHandler(svcOK)

	// 200 List
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	h.List()(w, req)
	used[w.Code] = true

	// 200 GetByID
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req = withChiURLParam(req, "id", "1")
	w = httptest.NewRecorder()
	h.GetByID()(w, req)
	used[w.Code] = true

	// 201
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"t"}`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	h.Create()(w, req)
	used[w.Code] = true

	// 400
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	h.Create()(w, req)
	used[w.Code] = true

	// 404
	notFoundSvc := &mockTagService{
		getByIDFn: func(_ context.Context, _ int64) (*model.Tag, error) {
			return nil, service.ErrTagNotFound
		},
	}
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req = withChiURLParam(req, "id", "999")
	w = httptest.NewRecorder()
	NewTagHandler(notFoundSvc).GetByID()(w, req)
	used[w.Code] = true

	// 409
	conflictSvc := &mockTagService{
		createFn: func(_ context.Context, _ service.TagInput) (*model.Tag, error) {
			return nil, service.ErrTagAlreadyExists
		},
	}
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"name":"dup"}`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	NewTagHandler(conflictSvc).Create()(w, req)
	used[w.Code] = true

	// 500
	errSvc := &mockTagService{
		listFn: func(_ context.Context) ([]*model.Tag, error) {
			return nil, errors.New("db error")
		},
	}
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	w = httptest.NewRecorder()
	NewTagHandler(errSvc).List()(w, req)
	used[w.Code] = true

	expected := []int{200, 201, 400, 404, 409, 500}
	for _, code := range expected {
		if !used[code] {
			t.Errorf("status code %d is not covered", code)
		}
	}
}
