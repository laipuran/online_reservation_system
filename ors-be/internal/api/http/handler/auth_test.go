package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"ors-be/internal/model"
	"ors-be/internal/service"
)

type mockAuthService struct {
	registerFn func(ctx context.Context, email, password, name string) (*service.RegisterResult, error)
	loginFn    func(ctx context.Context, email, password string) (*service.LoginResult, error)
}

func (m *mockAuthService) Register(ctx context.Context, email, password, name string) (*service.RegisterResult, error) {
	return m.registerFn(ctx, email, password, name)
}

func (m *mockAuthService) Login(ctx context.Context, email, password string) (*service.LoginResult, error) {
	return m.loginFn(ctx, email, password)
}

func newDefaultMock() *mockAuthService {
	return &mockAuthService{
		registerFn: func(_ context.Context, _, _, _ string) (*service.RegisterResult, error) {
			return nil, service.ErrInvalidEmail
		},
		loginFn: func(_ context.Context, _, _ string) (*service.LoginResult, error) {
			return nil, service.ErrInvalidCredentials
		},
	}
}

type testResponse struct {
	Code    int              `json:"code"`
	Message string           `json:"message"`
	Data    *json.RawMessage `json:"data"`
}

func decodeResponse(t *testing.T, body []byte) *testResponse {
	t.Helper()
	var resp testResponse
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("decode response: %v, body: %s", err, string(body))
	}
	return &resp
}

func TestRegister_Success(t *testing.T) {
	mock := &mockAuthService{
		registerFn: func(_ context.Context, email, password, name string) (*service.RegisterResult, error) {
			return &service.RegisterResult{
				User: &model.User{
					ID: 1, Name: name, Email: email, Role: "customer",
				},
				AccessToken: "test-token",
			}, nil
		},
	}
	h := NewAuthHandler(mock)

	body := `{"email":"test@example.com","password":"password123","name":"测试用户"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register()(w, req)

	if w.Code != http.StatusCreated {
		t.Errorf("status = %d, want %d", w.Code, http.StatusCreated)
	}

	resp := decodeResponse(t, w.Body.Bytes())
	if resp.Code != 201 {
		t.Errorf("code = %d, want %d", resp.Code, 201)
	}
	if resp.Message != "created" {
		t.Errorf("message = %s, want %s", resp.Message, "created")
	}
	if resp.Data == nil {
		t.Fatal("data is nil")
	}
}

func TestRegister_InvalidJSON(t *testing.T) {
	h := NewAuthHandler(&mockAuthService{
		registerFn: func(_ context.Context, _, _, _ string) (*service.RegisterResult, error) {
			return nil, service.ErrInvalidEmail
		},
	})

	body := `{invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	resp := decodeResponse(t, w.Body.Bytes())
	if resp.Code != 400 {
		t.Errorf("code = %d, want %d", resp.Code, 400)
	}
}

func TestRegister_InvalidEmail(t *testing.T) {
	mock := &mockAuthService{
		registerFn: func(_ context.Context, _, _, _ string) (*service.RegisterResult, error) {
			return nil, service.ErrInvalidEmail
		},
	}
	h := NewAuthHandler(mock)

	body := `{"email":"bad-email","password":"password123","name":"用户"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	resp := decodeResponse(t, w.Body.Bytes())
	if resp.Code != 400 {
		t.Errorf("code = %d, want %d", resp.Code, 400)
	}
}

func TestRegister_WeakPassword(t *testing.T) {
	mock := &mockAuthService{
		registerFn: func(_ context.Context, _, _, _ string) (*service.RegisterResult, error) {
			return nil, service.ErrWeakPassword
		},
	}
	h := NewAuthHandler(mock)

	body := `{"email":"test@test.com","password":"123","name":"用户"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestRegister_EmptyName(t *testing.T) {
	mock := &mockAuthService{
		registerFn: func(_ context.Context, _, _, _ string) (*service.RegisterResult, error) {
			return nil, service.ErrNameRequired
		},
	}
	h := NewAuthHandler(mock)

	body := `{"email":"test@test.com","password":"password123","name":""}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestRegister_EmailAlreadyRegistered(t *testing.T) {
	mock := &mockAuthService{
		registerFn: func(_ context.Context, _, _, _ string) (*service.RegisterResult, error) {
			return nil, service.ErrEmailAlreadyRegistered
		},
	}
	h := NewAuthHandler(mock)

	body := `{"email":"dup@example.com","password":"password123","name":"用户"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register()(w, req)

	if w.Code != http.StatusConflict {
		t.Errorf("status = %d, want %d", w.Code, http.StatusConflict)
	}
}

func TestRegister_UnexpectedError(t *testing.T) {
	mock := &mockAuthService{
		registerFn: func(_ context.Context, _, _, _ string) (*service.RegisterResult, error) {
			return nil, errors.New("internal server error")
		},
	}
	h := NewAuthHandler(mock)

	body := `{"email":"test@example.com","password":"password123","name":"用户"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register()(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestAuthHandler_Register_MissingFields(t *testing.T) {
	mock := newDefaultMock()
	h := NewAuthHandler(mock)

	tests := []struct {
		name string
		body string
	}{
		{"empty object", `{}`},
		{"only email", `{"email":"test@test.com"}`},
		{"missing name", `{"email":"test@test.com","password":"password123"}`},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()
			h.Register()(w, req)
			if w.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
			}
		})
	}
}

func TestLogin_Success(t *testing.T) {
	mock := &mockAuthService{
		loginFn: func(_ context.Context, email, password string) (*service.LoginResult, error) {
			return &service.LoginResult{
				User: &model.User{
					ID: 1, Name: "用户", Email: email, Role: "customer",
				},
				AccessToken: "test-token",
			}, nil
		},
	}
	h := NewAuthHandler(mock)

	body := `{"email":"test@example.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Login()(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("status = %d, want %d", w.Code, http.StatusOK)
	}

	resp := decodeResponse(t, w.Body.Bytes())
	if resp.Code != 200 {
		t.Errorf("code = %d, want %d", resp.Code, 200)
	}
	if resp.Data == nil {
		t.Fatal("data is nil")
	}
}

func TestLogin_InvalidJSON(t *testing.T) {
	h := NewAuthHandler(&mockAuthService{
		loginFn: func(_ context.Context, _, _ string) (*service.LoginResult, error) {
			return nil, service.ErrInvalidCredentials
		},
	})

	body := `not json`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Login()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestLogin_InvalidCredentials(t *testing.T) {
	mock := &mockAuthService{
		loginFn: func(_ context.Context, _, _ string) (*service.LoginResult, error) {
			return nil, service.ErrInvalidCredentials
		},
	}
	h := NewAuthHandler(mock)

	body := `{"email":"wrong@test.com","password":"wrongpass"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Login()(w, req)

	if w.Code != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", w.Code, http.StatusUnauthorized)
	}

	resp := decodeResponse(t, w.Body.Bytes())
	if resp.Message != "邮箱或密码错误" {
		t.Errorf("message = %s, want %s", resp.Message, "邮箱或密码错误")
	}
}

func TestLogin_UnexpectedError(t *testing.T) {
	mock := &mockAuthService{
		loginFn: func(_ context.Context, _, _ string) (*service.LoginResult, error) {
			return nil, errors.New("db connection failed")
		},
	}
	h := NewAuthHandler(mock)

	body := `{"email":"test@test.com","password":"password123"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Login()(w, req)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("status = %d, want %d", w.Code, http.StatusInternalServerError)
	}
}

func TestAuthHandler_Login_EmptyBody(t *testing.T) {
	h := NewAuthHandler(&mockAuthService{
		loginFn: func(_ context.Context, _, _ string) (*service.LoginResult, error) {
			return nil, service.ErrInvalidCredentials
		},
	})

	req := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Login()(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}
}

func TestAuthHandler_ResponseFormat(t *testing.T) {
	mock := &mockAuthService{
		registerFn: func(_ context.Context, _, _, _ string) (*service.RegisterResult, error) {
			return &service.RegisterResult{
				User: &model.User{ID: 1, Name: "用户", Email: "test@test.com", Role: "customer"},
				AccessToken: "token",
			}, nil
		},
	}
	h := NewAuthHandler(mock)

	body := `{"email":"test@test.com","password":"password123","name":"用户"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register()(w, req)

	var raw map[string]interface{}
	if err := json.Unmarshal(w.Body.Bytes(), &raw); err != nil {
		t.Fatalf("unmarshal error: %v", err)
	}

	if _, ok := raw["code"]; !ok {
		t.Error("response missing 'code' field")
	}
	if _, ok := raw["message"]; !ok {
		t.Error("response missing 'message' field")
	}
	if _, ok := raw["data"]; !ok {
		t.Error("response missing 'data' field")
	}
}

func TestAuthHandler_ContentType(t *testing.T) {
	mock := &mockAuthService{
		registerFn: func(_ context.Context, _, _, _ string) (*service.RegisterResult, error) {
			return &service.RegisterResult{
				User:        &model.User{ID: 1, Name: "用户", Email: "test@test.com", Role: "customer"},
				AccessToken: "token",
			}, nil
		},
	}
	h := NewAuthHandler(mock)

	body := `{"email":"test@test.com","password":"password123","name":"用户"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Register()(w, req)

	ct := w.Header().Get("Content-Type")
	if ct != "application/json" {
		t.Errorf("Content-Type = %s, want application/json", ct)
	}
}

// 验证 handler 不使用 response 包以外的 HTTP 状态码
func TestAuthHandler_StatusCodeCoverage(t *testing.T) {
	mock := &mockAuthService{
		registerFn: func(_ context.Context, _, _, _ string) (*service.RegisterResult, error) {
			return nil, service.ErrInvalidEmail
		},
		loginFn: func(_ context.Context, _, _ string) (*service.LoginResult, error) {
			return nil, service.ErrInvalidCredentials
		},
	}
	h := NewAuthHandler(mock)

	usedStatuses := map[int]bool{}

	// Register 201
	svcOK := &mockAuthService{
		registerFn: func(_ context.Context, _, _, _ string) (*service.RegisterResult, error) {
			return &service.RegisterResult{
				User:        &model.User{ID: 1, Name: "u", Email: "a@b.com", Role: "customer"},
				AccessToken: "t",
			}, nil
		},
	}
	body := `{"email":"a@b.com","password":"password123","name":"u"}`
	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	NewAuthHandler(svcOK).Register()(w, req)
	usedStatuses[w.Code] = true

	// Register 400
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{bad`))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	h.Register()(w, req)
	usedStatuses[w.Code] = true

	// Register 409
	mock409 := &mockAuthService{
		registerFn: func(_ context.Context, _, _, _ string) (*service.RegisterResult, error) {
			return nil, service.ErrEmailAlreadyRegistered
		},
	}
	body = `{"email":"a@b.com","password":"password123","name":"u"}`
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	NewAuthHandler(mock409).Register()(w, req)
	usedStatuses[w.Code] = true

	// Login 200
	loginOK := &mockAuthService{
		loginFn: func(_ context.Context, _, _ string) (*service.LoginResult, error) {
			return &service.LoginResult{
				User:        &model.User{ID: 1, Name: "u", Email: "a@b.com", Role: "customer"},
				AccessToken: "t",
			}, nil
		},
	}
	body = `{"email":"a@b.com","password":"password123"}`
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	NewAuthHandler(loginOK).Login()(w, req)
	usedStatuses[w.Code] = true

	// Login 401
	req = httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()
	h.Login()(w, req)
	usedStatuses[w.Code] = true

	expected := []int{200, 201, 400, 401, 409}
	for _, code := range expected {
		if !usedStatuses[code] {
			t.Errorf("status code %d is not covered by any test case", code)
		}
	}
}
