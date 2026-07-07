package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ors-be/internal/model"
)

type mockCategoryService struct {
	listFn func(ctx context.Context) ([]*model.Category, error)
}

func (m *mockCategoryService) List(ctx context.Context) ([]*model.Category, error) {
	return m.listFn(ctx)
}

func TestCategoryHandler_List_Success(t *testing.T) {
	parentID := int64(1)
	categories := []*model.Category{
		{ID: 1, Name: "医疗", Description: "医疗健康服务", CreatedAt: time.Now()},
		{ID: 2, Name: "口腔护理", ParentID: &parentID, CreatedAt: time.Now()},
	}
	h := NewCategoryHandler(&mockCategoryService{
		listFn: func(_ context.Context) ([]*model.Category, error) {
			return categories, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)
	w := httptest.NewRecorder()

	h.List()(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp struct {
		Code    int             `json:"code"`
		Message string          `json:"message"`
		Data    json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v, body: %s", err, w.Body.String())
	}
	if resp.Code != 200 {
		t.Errorf("code = %d, want 200", resp.Code)
	}
	if resp.Message != "ok" {
		t.Errorf("message = %s, want ok", resp.Message)
	}

	var data []model.Category
	if err := json.Unmarshal(resp.Data, &data); err != nil {
		t.Fatalf("decode data: %v", err)
	}
	if len(data) != 2 {
		t.Fatalf("data len = %d, want 2", len(data))
	}
	if data[0].Name != "医疗" {
		t.Errorf("data[0].Name = %s, want 医疗", data[0].Name)
	}
	if data[1].ParentID == nil || *data[1].ParentID != parentID {
		t.Fatalf("data[1].ParentID = %v, want %d", data[1].ParentID, parentID)
	}
}

func TestCategoryHandler_List_ServiceError(t *testing.T) {
	h := NewCategoryHandler(&mockCategoryService{
		listFn: func(_ context.Context) ([]*model.Category, error) {
			return nil, errors.New("db failed")
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/categories", nil)
	w := httptest.NewRecorder()

	h.List()(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}

	var resp struct {
		Code    int    `json:"code"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v, body: %s", err, w.Body.String())
	}
	if resp.Code != 500 {
		t.Errorf("code = %d, want 500", resp.Code)
	}
	if resp.Message != "分类列表查询失败" {
		t.Errorf("message = %s, want 分类列表查询失败", resp.Message)
	}
}
