package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"ors-be/internal/api/http/middleware"
	"ors-be/internal/auth"
	"ors-be/internal/model"
	"ors-be/internal/service"
)

type mockNotificationService struct {
	createFn      func(ctx context.Context, input service.NotificationInput) (*model.Notification, error)
	listMineFn    func(ctx context.Context, userID int64, isRead *bool, limit, offset int) ([]*model.Notification, error)
	countUnreadFn func(ctx context.Context, userID int64) (int64, error)
	markReadFn    func(ctx context.Context, userID, id int64) (*model.Notification, error)
	markAllReadFn func(ctx context.Context, userID int64) (int64, error)
}

func (m *mockNotificationService) Create(ctx context.Context, input service.NotificationInput) (*model.Notification, error) {
	return m.createFn(ctx, input)
}

func (m *mockNotificationService) ListMine(ctx context.Context, userID int64, isRead *bool, limit, offset int) ([]*model.Notification, error) {
	return m.listMineFn(ctx, userID, isRead, limit, offset)
}

func (m *mockNotificationService) CountUnread(ctx context.Context, userID int64) (int64, error) {
	return m.countUnreadFn(ctx, userID)
}

func (m *mockNotificationService) MarkRead(ctx context.Context, userID, id int64) (*model.Notification, error) {
	return m.markReadFn(ctx, userID, id)
}

func (m *mockNotificationService) MarkAllRead(ctx context.Context, userID int64) (int64, error) {
	return m.markAllReadFn(ctx, userID)
}

func TestNotificationHandler_ListMine_Success(t *testing.T) {
	h := NewNotificationHandler(&mockNotificationService{
		listMineFn: func(ctx context.Context, userID int64, isRead *bool, limit, offset int) ([]*model.Notification, error) {
			if userID != 7 {
				t.Errorf("userID = %d, want 7", userID)
			}
			if isRead == nil || *isRead {
				t.Fatalf("isRead = %v, want false", isRead)
			}
			if limit != 5 || offset != 5 {
				t.Errorf("limit/offset = %d/%d, want 5/5", limit, offset)
			}
			return []*model.Notification{{ID: 1, UserID: userID, Title: "Notice", Type: service.NotificationTypeSystem}}, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/notifications?is_read=false&page=2&page_size=5", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserCtxKey, &auth.Claims{UserID: 7, Role: "customer"}))
	w := httptest.NewRecorder()

	h.ListMine()(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	var resp struct {
		Code int `json:"code"`
		Data struct {
			Items    []model.Notification `json:"items"`
			Page     int                  `json:"page"`
			PageSize int                  `json:"page_size"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if len(resp.Data.Items) != 1 || resp.Data.Page != 2 || resp.Data.PageSize != 5 {
		t.Fatalf("response = %+v", resp)
	}
}

func TestNotificationHandler_ListMine_InvalidIsRead(t *testing.T) {
	h := NewNotificationHandler(&mockNotificationService{})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/notifications?is_read=yes", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserCtxKey, &auth.Claims{UserID: 7, Role: "customer"}))
	w := httptest.NewRecorder()

	h.ListMine()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestNotificationHandler_CountUnread_Success(t *testing.T) {
	h := NewNotificationHandler(&mockNotificationService{
		countUnreadFn: func(ctx context.Context, userID int64) (int64, error) {
			if userID != 7 {
				t.Errorf("userID = %d, want 7", userID)
			}
			return 3, nil
		},
	})

	req := httptest.NewRequest(http.MethodGet, "/api/v1/notifications/unread-count", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserCtxKey, &auth.Claims{UserID: 7, Role: "customer"}))
	w := httptest.NewRecorder()

	h.CountUnread()(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	var resp struct {
		Data struct {
			Count int64 `json:"count"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Data.Count != 3 {
		t.Fatalf("count = %d, want 3", resp.Data.Count)
	}
}

func TestNotificationHandler_MarkRead_Success(t *testing.T) {
	h := NewNotificationHandler(&mockNotificationService{
		markReadFn: func(ctx context.Context, userID, id int64) (*model.Notification, error) {
			if userID != 7 || id != 11 {
				t.Errorf("userID/id = %d/%d, want 7/11", userID, id)
			}
			return &model.Notification{ID: id, UserID: userID, Title: "Notice", IsRead: true, CreatedAt: time.Now()}, nil
		},
	})

	req := requestWithURLParam(http.MethodPut, "/api/v1/notifications/11/read", "id", "11", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserCtxKey, &auth.Claims{UserID: 7, Role: "customer"}))
	w := httptest.NewRecorder()

	h.MarkRead()(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
}

func TestNotificationHandler_MarkRead_NotFound(t *testing.T) {
	h := NewNotificationHandler(&mockNotificationService{
		markReadFn: func(ctx context.Context, userID, id int64) (*model.Notification, error) {
			return nil, service.ErrNotificationNotFound
		},
	})

	req := requestWithURLParam(http.MethodPut, "/api/v1/notifications/11/read", "id", "11", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserCtxKey, &auth.Claims{UserID: 7, Role: "customer"}))
	w := httptest.NewRecorder()

	h.MarkRead()(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestNotificationHandler_MarkRead_InvalidID(t *testing.T) {
	h := NewNotificationHandler(&mockNotificationService{})

	req := requestWithURLParam(http.MethodPut, "/api/v1/notifications/0/read", "id", "0", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserCtxKey, &auth.Claims{UserID: 7, Role: "customer"}))
	w := httptest.NewRecorder()

	h.MarkRead()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestNotificationHandler_MarkAllRead_Success(t *testing.T) {
	h := NewNotificationHandler(&mockNotificationService{
		markAllReadFn: func(ctx context.Context, userID int64) (int64, error) {
			if userID != 7 {
				t.Errorf("userID = %d, want 7", userID)
			}
			return 4, nil
		},
	})

	req := httptest.NewRequest(http.MethodPut, "/api/v1/notifications/read-all", nil)
	req = req.WithContext(context.WithValue(req.Context(), middleware.UserCtxKey, &auth.Claims{UserID: 7, Role: "customer"}))
	w := httptest.NewRecorder()

	h.MarkAllRead()(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}
	var resp struct {
		Data struct {
			UpdatedCount int64 `json:"updated_count"`
		} `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	if resp.Data.UpdatedCount != 4 {
		t.Fatalf("updated_count = %d, want 4", resp.Data.UpdatedCount)
	}
}
