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

type mockServiceService struct {
	createFn      func(ctx context.Context, userID int64, input service.ServiceInput) (*model.ServiceView, error)
	getByIDFn     func(ctx context.Context, id int64) (*model.ServiceView, error)
	listFn        func(ctx context.Context, filter model.ServiceFilter) (*service.ServiceListResult, error)
	updateFn      func(ctx context.Context, userID, id int64, input service.ServiceInput) (*model.ServiceView, error)
	updateStatusFn func(ctx context.Context, userID, id int64, status string) (*model.ServiceView, error)
	listTagsFn    func(ctx context.Context, id int64) ([]*model.Tag, error)
	replaceTagsFn func(ctx context.Context, userID, id int64, input service.ServiceTagsInput) ([]*model.Tag, error)
}

func (m *mockServiceService) Create(ctx context.Context, userID int64, input service.ServiceInput) (*model.ServiceView, error) {
	return m.createFn(ctx, userID, input)
}
func (m *mockServiceService) GetByID(ctx context.Context, id int64) (*model.ServiceView, error) {
	return m.getByIDFn(ctx, id)
}
func (m *mockServiceService) List(ctx context.Context, filter model.ServiceFilter) (*service.ServiceListResult, error) {
	return m.listFn(ctx, filter)
}
func (m *mockServiceService) Update(ctx context.Context, userID, id int64, input service.ServiceInput) (*model.ServiceView, error) {
	return m.updateFn(ctx, userID, id, input)
}
func (m *mockServiceService) UpdateStatus(ctx context.Context, userID, id int64, status string) (*model.ServiceView, error) {
	return m.updateStatusFn(ctx, userID, id, status)
}
func (m *mockServiceService) ListTags(ctx context.Context, id int64) ([]*model.Tag, error) {
	return m.listTagsFn(ctx, id)
}
func (m *mockServiceService) ReplaceTags(ctx context.Context, userID, id int64, input service.ServiceTagsInput) ([]*model.Tag, error) {
	return m.replaceTagsFn(ctx, userID, id, input)
}

func stubServiceView() *model.ServiceView {
	return &model.ServiceView{
		ID:              1,
		Title:           "肩颈按摩 60 分钟",
		Price:           199,
		DurationMinutes: 60,
		Status:          "active",
		Provider: model.ServiceProviderSummary{
			ID: 1, BusinessName: "舒心养生馆",
		},
		Category: model.CategorySummary{
			ID: 1, Name: "医疗",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// ────── List ──────

func TestServiceHandler_List_Success(t *testing.T) {
	svc := &mockServiceService{
		listFn: func(_ context.Context, filter model.ServiceFilter) (*service.ServiceListResult, error) {
			if filter.Status != "active" {
				t.Errorf("List() status = %s, want active", filter.Status)
			}
			return &service.ServiceListResult{
				Items:    []*model.ServiceView{stubServiceView()},
				Total:    1,
				Page:     1,
				PageSize: 20,
			}, nil
		},
	}
	h := NewServiceHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/services", nil)
	w := httptest.NewRecorder()

	h.List()(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp struct {
		Code int                         `json:"code"`
		Data service.ServiceListResult   `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Data.Total != 1 {
		t.Errorf("Total = %d, want 1", resp.Data.Total)
	}
	if len(resp.Data.Items) != 1 {
		t.Fatalf("items len = %d", len(resp.Data.Items))
	}
	if resp.Data.Items[0].Provider.BusinessName != "舒心养生馆" {
		t.Errorf("BusinessName = %s", resp.Data.Items[0].Provider.BusinessName)
	}
}

func TestServiceHandler_List_BadCategory(t *testing.T) {
	h := NewServiceHandler(&mockServiceService{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/services?category_id=abc", nil)
	w := httptest.NewRecorder()

	h.List()(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestServiceHandler_List_BadPage(t *testing.T) {
	h := NewServiceHandler(&mockServiceService{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/services?page=abc", nil)
	w := httptest.NewRecorder()

	h.List()(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestServiceHandler_List_ServiceError(t *testing.T) {
	svc := &mockServiceService{
		listFn: func(_ context.Context, _ model.ServiceFilter) (*service.ServiceListResult, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewServiceHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/services", nil)
	w := httptest.NewRecorder()

	h.List()(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// ────── ListByProvider ──────

func TestServiceHandler_ListByProvider_Success(t *testing.T) {
	svc := &mockServiceService{
		listFn: func(_ context.Context, filter model.ServiceFilter) (*service.ServiceListResult, error) {
			if filter.ProviderID == nil || *filter.ProviderID != 1 {
				t.Errorf("ProviderID = %v, want 1", filter.ProviderID)
			}
			return &service.ServiceListResult{Items: []*model.ServiceView{}, Total: 0, Page: 1, PageSize: 20}, nil
		},
	}
	h := NewServiceHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/providers/1/services", nil)
	req = withChiURLParam(req, "id", "1")
	w := httptest.NewRecorder()

	h.ListByProvider()(w, req)
	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestServiceHandler_ListByProvider_InvalidID(t *testing.T) {
	h := NewServiceHandler(&mockServiceService{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/providers/abc/services", nil)
	req = withChiURLParam(req, "id", "abc")
	w := httptest.NewRecorder()

	h.ListByProvider()(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

// ────── GetByID ──────

func TestServiceHandler_GetByID_Success(t *testing.T) {
	svc := &mockServiceService{
		getByIDFn: func(_ context.Context, id int64) (*model.ServiceView, error) {
			sv := stubServiceView()
			sv.ID = id
			return sv, nil
		},
	}
	h := NewServiceHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/services/1", nil)
	req = withChiURLParam(req, "id", "1")
	w := httptest.NewRecorder()

	h.GetByID()(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestServiceHandler_GetByID_InvalidID(t *testing.T) {
	h := NewServiceHandler(&mockServiceService{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/services/abc", nil)
	req = withChiURLParam(req, "id", "abc")
	w := httptest.NewRecorder()

	h.GetByID()(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestServiceHandler_GetByID_NotFound(t *testing.T) {
	svc := &mockServiceService{
		getByIDFn: func(_ context.Context, _ int64) (*model.ServiceView, error) {
			return nil, service.ErrServiceNotFound
		},
	}
	h := NewServiceHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/services/999", nil)
	req = withChiURLParam(req, "id", "999")
	w := httptest.NewRecorder()

	h.GetByID()(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// ────── Create ──────

func TestServiceHandler_Create_Success(t *testing.T) {
	svc := &mockServiceService{
		createFn: func(_ context.Context, userID int64, input service.ServiceInput) (*model.ServiceView, error) {
			return stubServiceView(), nil
		},
	}
	h := NewServiceHandler(svc)

	body := `{"category_id":1,"title":"肩颈按摩 60 分钟","price":199,"duration_minutes":60}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/services", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	w := httptest.NewRecorder()

	h.Create()(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusCreated)
	}
}

func TestServiceHandler_Create_Unauthorized(t *testing.T) {
	h := NewServiceHandler(&mockServiceService{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/services", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Create()(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestServiceHandler_Create_InvalidJSON(t *testing.T) {
	h := NewServiceHandler(&mockServiceService{})
	req := httptest.NewRequest(http.MethodPost, "/api/v1/services", strings.NewReader(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	w := httptest.NewRecorder()

	h.Create()(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestServiceHandler_Create_ValidationError(t *testing.T) {
	svc := &mockServiceService{
		createFn: func(_ context.Context, _ int64, _ service.ServiceInput) (*model.ServiceView, error) {
			return nil, service.ErrServiceTitleRequired
		},
	}
	h := NewServiceHandler(svc)

	body := `{"category_id":1,"title":"","price":199}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/services", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	w := httptest.NewRecorder()

	h.Create()(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestServiceHandler_Create_ProviderNotFound(t *testing.T) {
	svc := &mockServiceService{
		createFn: func(_ context.Context, _ int64, _ service.ServiceInput) (*model.ServiceView, error) {
			return nil, service.ErrProviderNotFound
		},
	}
	h := NewServiceHandler(svc)

	body := `{"category_id":1,"title":"test","price":100,"duration_minutes":30}`
	req := httptest.NewRequest(http.MethodPost, "/api/v1/services", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 999)
	w := httptest.NewRecorder()

	h.Create()(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// ────── Update ──────

func TestServiceHandler_Update_Success(t *testing.T) {
	svc := &mockServiceService{
		updateFn: func(_ context.Context, _, _ int64, _ service.ServiceInput) (*model.ServiceView, error) {
			return stubServiceView(), nil
		},
	}
	h := NewServiceHandler(svc)

	body := `{"category_id":1,"title":"更新标题","price":299,"duration_minutes":45}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/services/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	req = withChiURLParam(req, "id", "1")
	w := httptest.NewRecorder()

	h.Update()(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestServiceHandler_Update_Unauthorized(t *testing.T) {
	h := NewServiceHandler(&mockServiceService{})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/services/1", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	req = withChiURLParam(req, "id", "1")
	w := httptest.NewRecorder()

	h.Update()(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestServiceHandler_Update_Forbidden(t *testing.T) {
	svc := &mockServiceService{
		updateFn: func(_ context.Context, _, _ int64, _ service.ServiceInput) (*model.ServiceView, error) {
			return nil, service.ErrServiceForbidden
		},
	}
	h := NewServiceHandler(svc)

	body := `{"category_id":1,"title":"test","price":100,"duration_minutes":30}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/services/1", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	req = withChiURLParam(req, "id", "1")
	w := httptest.NewRecorder()

	h.Update()(w, req)
	if w.Code != http.StatusForbidden {
		t.Errorf("status = %d, want %d", w.Code, http.StatusForbidden)
	}
}

// ────── UpdateStatus ──────

func TestServiceHandler_UpdateStatus_Success(t *testing.T) {
	svc := &mockServiceService{
		updateStatusFn: func(_ context.Context, _, _ int64, status string) (*model.ServiceView, error) {
			sv := stubServiceView()
			sv.Status = status
			return sv, nil
		},
	}
	h := NewServiceHandler(svc)

	body := `{"status":"inactive"}`
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/services/1/status", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	req = withChiURLParam(req, "id", "1")
	w := httptest.NewRecorder()

	h.UpdateStatus()(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestServiceHandler_UpdateStatus_InvalidStatus(t *testing.T) {
	svc := &mockServiceService{
		updateStatusFn: func(_ context.Context, _, _ int64, _ string) (*model.ServiceView, error) {
			return nil, service.ErrServiceInvalidStatus
		},
	}
	h := NewServiceHandler(svc)

	body := `{"status":"invalid"}`
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/services/1/status", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	req = withChiURLParam(req, "id", "1")
	w := httptest.NewRecorder()

	h.UpdateStatus()(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestServiceHandler_UpdateStatus_Unauthorized(t *testing.T) {
	h := NewServiceHandler(&mockServiceService{})
	req := httptest.NewRequest(http.MethodPatch, "/api/v1/services/1/status", strings.NewReader(`{"status":"inactive"}`))
	req.Header.Set("Content-Type", "application/json")
	req = withChiURLParam(req, "id", "1")
	w := httptest.NewRecorder()

	h.UpdateStatus()(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

// ────── ListTags ──────

func TestServiceHandler_ListTags_Success(t *testing.T) {
	svc := &mockServiceService{
		listTagsFn: func(_ context.Context, _ int64) ([]*model.Tag, error) {
			return []*model.Tag{
				{ID: 1, Name: "放松"},
				{ID: 2, Name: "塑形"},
			}, nil
		},
	}
	h := NewServiceHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/services/1/tags", nil)
	req = withChiURLParam(req, "id", "1")
	w := httptest.NewRecorder()

	h.ListTags()(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestServiceHandler_ListTags_ServiceNotFound(t *testing.T) {
	svc := &mockServiceService{
		listTagsFn: func(_ context.Context, _ int64) ([]*model.Tag, error) {
			return nil, service.ErrServiceNotFound
		},
	}
	h := NewServiceHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/services/999/tags", nil)
	req = withChiURLParam(req, "id", "999")
	w := httptest.NewRecorder()

	h.ListTags()(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// ────── ReplaceTags ──────

func TestServiceHandler_ReplaceTags_Success(t *testing.T) {
	svc := &mockServiceService{
		replaceTagsFn: func(_ context.Context, _, _ int64, _ service.ServiceTagsInput) ([]*model.Tag, error) {
			return []*model.Tag{{ID: 1, Name: "医疗"}}, nil
		},
	}
	h := NewServiceHandler(svc)

	body := `{"tag_ids":[1,2,3]}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/services/1/tags", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	req = withChiURLParam(req, "id", "1")
	w := httptest.NewRecorder()

	h.ReplaceTags()(w, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestServiceHandler_ReplaceTags_Unauthorized(t *testing.T) {
	h := NewServiceHandler(&mockServiceService{})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/services/1/tags", strings.NewReader(`{"tag_ids":[1]}`))
	req.Header.Set("Content-Type", "application/json")
	req = withChiURLParam(req, "id", "1")
	w := httptest.NewRecorder()

	h.ReplaceTags()(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestServiceHandler_ReplaceTags_InvalidTag(t *testing.T) {
	svc := &mockServiceService{
		replaceTagsFn: func(_ context.Context, _, _ int64, _ service.ServiceTagsInput) ([]*model.Tag, error) {
			return nil, service.ErrServiceInvalidTag
		},
	}
	h := NewServiceHandler(svc)

	body := `{"tag_ids":[999]}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/services/1/tags", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	req = withChiURLParam(req, "id", "1")
	w := httptest.NewRecorder()

	h.ReplaceTags()(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestServiceHandler_ReplaceTags_TagNotFound(t *testing.T) {
	svc := &mockServiceService{
		replaceTagsFn: func(_ context.Context, _, _ int64, _ service.ServiceTagsInput) ([]*model.Tag, error) {
			return nil, service.ErrTagNotFound
		},
	}
	h := NewServiceHandler(svc)

	body := `{"tag_ids":[999]}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/services/1/tags", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	req = withChiURLParam(req, "id", "1")
	w := httptest.NewRecorder()

	h.ReplaceTags()(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// ────── Status code coverage ──────

func TestServiceHandler_StatusCodeCoverage(t *testing.T) {
	used := map[int]bool{}

	svcOK := &mockServiceService{
		listFn: func(_ context.Context, _ model.ServiceFilter) (*service.ServiceListResult, error) {
			return &service.ServiceListResult{Items: []*model.ServiceView{}, Total: 0, Page: 1, PageSize: 20}, nil
		},
		getByIDFn: func(_ context.Context, _ int64) (*model.ServiceView, error) {
			return stubServiceView(), nil
		},
		createFn: func(_ context.Context, _ int64, _ service.ServiceInput) (*model.ServiceView, error) {
			return stubServiceView(), nil
		},
		updateFn: func(_ context.Context, _, _ int64, _ service.ServiceInput) (*model.ServiceView, error) {
			return stubServiceView(), nil
		},
		updateStatusFn: func(_ context.Context, _, _ int64, _ string) (*model.ServiceView, error) {
			return stubServiceView(), nil
		},
		listTagsFn: func(_ context.Context, _ int64) ([]*model.Tag, error) {
			return []*model.Tag{}, nil
		},
		replaceTagsFn: func(_ context.Context, _, _ int64, _ service.ServiceTagsInput) ([]*model.Tag, error) {
			return []*model.Tag{}, nil
		},
	}
	h := NewServiceHandler(svcOK)

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

	// 200 Update
	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{"category_id":1,"title":"t","price":100,"duration_minutes":30}`))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	req = withChiURLParam(req, "id", "1")
	w = httptest.NewRecorder()
	h.Update()(w, req)
	used[w.Code] = true

	// 200 UpdateStatus
	req = httptest.NewRequest(http.MethodPatch, "/", strings.NewReader(`{"status":"inactive"}`))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	req = withChiURLParam(req, "id", "1")
	w = httptest.NewRecorder()
	h.UpdateStatus()(w, req)
	used[w.Code] = true

	// 200 ListTags
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req = withChiURLParam(req, "id", "1")
	w = httptest.NewRecorder()
	h.ListTags()(w, req)
	used[w.Code] = true

	// 200 ReplaceTags
	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{"tag_ids":[1]}`))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	req = withChiURLParam(req, "id", "1")
	w = httptest.NewRecorder()
	h.ReplaceTags()(w, req)
	used[w.Code] = true

	// 201 Create
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"category_id":1,"title":"t","price":100,"duration_minutes":30}`))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	w = httptest.NewRecorder()
	h.Create()(w, req)
	used[w.Code] = true

	// 400 Bad category_id
	req = httptest.NewRequest(http.MethodGet, "/?category_id=abc", nil)
	w = httptest.NewRecorder()
	h.List()(w, req)
	used[w.Code] = true

	// 401
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	h.Create()(w, req)
	used[w.Code] = true

	// 403
	forbiddenSvc := &mockServiceService{
		updateFn: func(_ context.Context, _, _ int64, _ service.ServiceInput) (*model.ServiceView, error) {
			return nil, service.ErrServiceForbidden
		},
	}
	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{"category_id":1,"title":"t","price":100,"duration_minutes":30}`))
	req.Header.Set("Content-Type", "application/json")
	req = withProviderClaims(req, 2)
	req = withChiURLParam(req, "id", "1")
	w = httptest.NewRecorder()
	NewServiceHandler(forbiddenSvc).Update()(w, req)
	used[w.Code] = true

	// 404
	notFoundSvc := &mockServiceService{
		getByIDFn: func(_ context.Context, _ int64) (*model.ServiceView, error) {
			return nil, service.ErrServiceNotFound
		},
	}
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req = withChiURLParam(req, "id", "999")
	w = httptest.NewRecorder()
	NewServiceHandler(notFoundSvc).GetByID()(w, req)
	used[w.Code] = true

	// 500
	errSvc := &mockServiceService{
		listFn: func(_ context.Context, _ model.ServiceFilter) (*service.ServiceListResult, error) {
			return nil, errors.New("db error")
		},
	}
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	w = httptest.NewRecorder()
	NewServiceHandler(errSvc).List()(w, req)
	used[w.Code] = true

	expected := []int{200, 201, 400, 401, 403, 404, 500}
	for _, code := range expected {
		if !used[code] {
			t.Errorf("status code %d is not covered", code)
		}
	}
}
