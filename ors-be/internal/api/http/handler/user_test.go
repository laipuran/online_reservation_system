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

type mockUserService struct {
	getMineFn       func(ctx context.Context, userID int64) (*model.User, error)
	updateMineFn    func(ctx context.Context, userID int64, input service.UserInput) (*model.User, error)
	updatePasswordFn func(ctx context.Context, userID int64, input service.UserPasswordInput) error
}

func (m *mockUserService) GetMine(ctx context.Context, userID int64) (*model.User, error) {
	return m.getMineFn(ctx, userID)
}
func (m *mockUserService) UpdateMine(ctx context.Context, userID int64, input service.UserInput) (*model.User, error) {
	return m.updateMineFn(ctx, userID, input)
}
func (m *mockUserService) UpdatePassword(ctx context.Context, userID int64, input service.UserPasswordInput) error {
	return m.updatePasswordFn(ctx, userID, input)
}

func stubUser() *model.User {
	return &model.User{
		ID:        1,
		Name:      "测试用户",
		Email:     "test@example.com",
		Role:      "customer",
		Phone:     "13800138000",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// ────── GetMine ──────

func TestUserHandler_GetMine_Success(t *testing.T) {
	svc := &mockUserService{
		getMineFn: func(_ context.Context, userID int64) (*model.User, error) {
			return stubUser(), nil
		},
	}
	h := NewUserHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.GetMine()(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp struct {
		Code int         `json:"code"`
		Data model.User  `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Code != 200 {
		t.Errorf("code = %d, want 200", resp.Code)
	}
	if resp.Data.Name != "测试用户" {
		t.Errorf("Name = %s", resp.Data.Name)
	}
	if resp.Data.Email != "test@example.com" {
		t.Errorf("Email = %s", resp.Data.Email)
	}
	if resp.Data.PasswordHash != "" {
		t.Error("PasswordHash should not be serialized (json:\"-\")")
	}
}

func TestUserHandler_GetMine_Unauthorized(t *testing.T) {
	h := NewUserHandler(&mockUserService{})
	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	w := httptest.NewRecorder()

	h.GetMine()(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestUserHandler_GetMine_NotFound(t *testing.T) {
	svc := &mockUserService{
		getMineFn: func(_ context.Context, _ int64) (*model.User, error) {
			return nil, service.ErrUserNotFound
		},
	}
	h := NewUserHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req = withCustomerClaims(req, 999)
	w := httptest.NewRecorder()

	h.GetMine()(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestUserHandler_GetMine_UnexpectedError(t *testing.T) {
	svc := &mockUserService{
		getMineFn: func(_ context.Context, _ int64) (*model.User, error) {
			return nil, errors.New("db connection lost")
		},
	}
	h := NewUserHandler(svc)

	req := httptest.NewRequest(http.MethodGet, "/api/v1/users/me", nil)
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.GetMine()(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// ────── UpdateMine ──────

func TestUserHandler_UpdateMine_Success(t *testing.T) {
	svc := &mockUserService{
		updateMineFn: func(_ context.Context, _ int64, input service.UserInput) (*model.User, error) {
			user := stubUser()
			user.Name = input.Name
			user.Phone = input.Phone
			return user, nil
		},
	}
	h := NewUserHandler(svc)

	body := `{"name":"新名字","phone":"13900139000"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.UpdateMine()(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp struct {
		Code int         `json:"code"`
		Data model.User  `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Data.Name != "新名字" {
		t.Errorf("Name = %s, want 新名字", resp.Data.Name)
	}
	if resp.Data.Phone != "13900139000" {
		t.Errorf("Phone = %s, want 13900139000", resp.Data.Phone)
	}
}

func TestUserHandler_UpdateMine_Unauthorized(t *testing.T) {
	h := NewUserHandler(&mockUserService{})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me", strings.NewReader(`{"name":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UpdateMine()(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestUserHandler_UpdateMine_InvalidJSON(t *testing.T) {
	h := NewUserHandler(&mockUserService{})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me", strings.NewReader(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.UpdateMine()(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestUserHandler_UpdateMine_EmptyName(t *testing.T) {
	svc := &mockUserService{
		updateMineFn: func(_ context.Context, _ int64, _ service.UserInput) (*model.User, error) {
			return nil, service.ErrNameRequired
		},
	}
	h := NewUserHandler(svc)

	body := `{"name":""}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.UpdateMine()(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestUserHandler_UpdateMine_NotFound(t *testing.T) {
	svc := &mockUserService{
		updateMineFn: func(_ context.Context, _ int64, _ service.UserInput) (*model.User, error) {
			return nil, service.ErrUserNotFound
		},
	}
	h := NewUserHandler(svc)

	body := `{"name":"test"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 999)
	w := httptest.NewRecorder()

	h.UpdateMine()(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

// ────── UpdatePassword ──────

func TestUserHandler_UpdatePassword_Success(t *testing.T) {
	svc := &mockUserService{
		updatePasswordFn: func(_ context.Context, _ int64, input service.UserPasswordInput) error {
			if input.CurrentPassword != "oldpass" || input.NewPassword != "newpass123" {
				return service.ErrCurrentPasswordWrong
			}
			return nil
		},
	}
	h := NewUserHandler(svc)

	body := `{"current_password":"oldpass","new_password":"newpass123"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me/password", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.UpdatePassword()(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp struct {
		Code int             `json:"code"`
		Data json.RawMessage `json:"data"`
	}
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp.Code != 200 {
		t.Errorf("code = %d, want 200", resp.Code)
	}
}

func TestUserHandler_UpdatePassword_Unauthorized(t *testing.T) {
	h := NewUserHandler(&mockUserService{})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me/password", strings.NewReader(`{}`))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.UpdatePassword()(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestUserHandler_UpdatePassword_InvalidJSON(t *testing.T) {
	h := NewUserHandler(&mockUserService{})
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me/password", strings.NewReader(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.UpdatePassword()(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestUserHandler_UpdatePassword_WrongPassword(t *testing.T) {
	svc := &mockUserService{
		updatePasswordFn: func(_ context.Context, _ int64, _ service.UserPasswordInput) error {
			return service.ErrCurrentPasswordWrong
		},
	}
	h := NewUserHandler(svc)

	body := `{"current_password":"wrongpass","new_password":"newpass123"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me/password", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.UpdatePassword()(w, req)
	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}
}

func TestUserHandler_UpdatePassword_WeakPassword(t *testing.T) {
	svc := &mockUserService{
		updatePasswordFn: func(_ context.Context, _ int64, _ service.UserPasswordInput) error {
			return service.ErrWeakPassword
		},
	}
	h := NewUserHandler(svc)

	body := `{"current_password":"oldpass","new_password":"123"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me/password", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.UpdatePassword()(w, req)
	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestUserHandler_UpdatePassword_UserNotFound(t *testing.T) {
	svc := &mockUserService{
		updatePasswordFn: func(_ context.Context, _ int64, _ service.UserPasswordInput) error {
			return service.ErrUserNotFound
		},
	}
	h := NewUserHandler(svc)

	body := `{"current_password":"oldpass","new_password":"newpass123"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me/password", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 999)
	w := httptest.NewRecorder()

	h.UpdatePassword()(w, req)
	if w.Code != http.StatusNotFound {
		t.Errorf("status = %d, want %d", w.Code, http.StatusNotFound)
	}
}

func TestUserHandler_UpdatePassword_UnexpectedError(t *testing.T) {
	svc := &mockUserService{
		updatePasswordFn: func(_ context.Context, _ int64, _ service.UserPasswordInput) error {
			return errors.New("db error")
		},
	}
	h := NewUserHandler(svc)

	body := `{"current_password":"oldpass","new_password":"newpass123"}`
	req := httptest.NewRequest(http.MethodPut, "/api/v1/users/me/password", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()

	h.UpdatePassword()(w, req)
	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

// ────── Status code coverage ──────

func TestUserHandler_StatusCodeCoverage(t *testing.T) {
	used := map[int]bool{}

	svcOK := &mockUserService{
		getMineFn: func(_ context.Context, _ int64) (*model.User, error) {
			return stubUser(), nil
		},
		updateMineFn: func(_ context.Context, _ int64, _ service.UserInput) (*model.User, error) {
			return stubUser(), nil
		},
		updatePasswordFn: func(_ context.Context, _ int64, _ service.UserPasswordInput) error {
			return nil
		},
	}
	h := NewUserHandler(svcOK)

	// 200 GetMine
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = withCustomerClaims(req, 1)
	w := httptest.NewRecorder()
	h.GetMine()(w, req)
	used[w.Code] = true

	// 200 UpdateMine
	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{"name":"x"}`))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w = httptest.NewRecorder()
	h.UpdateMine()(w, req)
	used[w.Code] = true

	// 200 UpdatePassword
	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{"current_password":"old","new_password":"newpass123"}`))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w = httptest.NewRecorder()
	h.UpdatePassword()(w, req)
	used[w.Code] = true

	// 400 Invalid JSON
	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w = httptest.NewRecorder()
	h.UpdateMine()(w, req)
	used[w.Code] = true

	// 401 Missing claims
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	w = httptest.NewRecorder()
	h.GetMine()(w, req)
	used[w.Code] = true

	// 401 Wrong password
	wrongPwdSvc := &mockUserService{
		updatePasswordFn: func(_ context.Context, _ int64, _ service.UserPasswordInput) error {
			return service.ErrCurrentPasswordWrong
		},
	}
	req = httptest.NewRequest(http.MethodPut, "/", strings.NewReader(`{"current_password":"wrong","new_password":"newpass123"}`))
	req.Header.Set("Content-Type", "application/json")
	req = withCustomerClaims(req, 1)
	w = httptest.NewRecorder()
	NewUserHandler(wrongPwdSvc).UpdatePassword()(w, req)
	used[w.Code] = true

	// 404 Not found
	notFoundSvc := &mockUserService{
		getMineFn: func(_ context.Context, _ int64) (*model.User, error) {
			return nil, service.ErrUserNotFound
		},
	}
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req = withCustomerClaims(req, 999)
	w = httptest.NewRecorder()
	NewUserHandler(notFoundSvc).GetMine()(w, req)
	used[w.Code] = true

	// 500
	errSvc := &mockUserService{
		getMineFn: func(_ context.Context, _ int64) (*model.User, error) {
			return nil, errors.New("db error")
		},
	}
	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req = withCustomerClaims(req, 1)
	w = httptest.NewRecorder()
	NewUserHandler(errSvc).GetMine()(w, req)
	used[w.Code] = true

	expected := []int{200, 400, 401, 404, 500}
	for _, code := range expected {
		if !used[code] {
			t.Errorf("status code %d is not covered", code)
		}
	}
}
