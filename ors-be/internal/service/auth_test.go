package service

import (
	"context"
	"errors"
	"testing"

	"ors-be/internal/auth"
	"ors-be/internal/model"
)

type mockUserRepo struct {
	users map[string]*model.User
}

func newMockUserRepo() *mockUserRepo {
	return &mockUserRepo{users: make(map[string]*model.User)}
}

func (m *mockUserRepo) Create(ctx context.Context, user *model.User) error {
	if _, exists := m.users[user.Email]; exists {
		return errors.New("duplicate key")
	}
	user.ID = int64(len(m.users) + 1)
	m.users[user.Email] = user
	return nil
}

func (m *mockUserRepo) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	user, exists := m.users[email]
	if !exists {
		return nil, nil
	}
	return user, nil
}

func (m *mockUserRepo) GetByID(ctx context.Context, id int64) (*model.User, error) {
	for _, u := range m.users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil
}

func newTestService() AuthService {
	return NewAuthService(
		newMockUserRepo(),
		auth.NewHasher(),
		auth.NewTokenGenerator("test-secret", 24),
	)
}

func TestAuthService_Register_Success(t *testing.T) {
	svc := newTestService()

	result, err := svc.Register(context.Background(), "test@example.com", "password123", "测试用户", "")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if result.User == nil {
		t.Fatal("Register() result.User is nil")
	}
	if result.User.Email != "test@example.com" {
		t.Errorf("Register() email = %s, want %s", result.User.Email, "test@example.com")
	}
	if result.User.Name != "测试用户" {
		t.Errorf("Register() name = %s, want %s", result.User.Name, "测试用户")
	}
	if result.User.Role != "customer" {
		t.Errorf("Register() role = %s, want %s", result.User.Role, "customer")
	}
	if result.User.ID == 0 {
		t.Error("Register() ID should be non-zero")
	}
	if result.AccessToken == "" {
		t.Error("Register() AccessToken should not be empty")
	}
}

func TestAuthService_RegisterProvider_ThenLoginReturnsProviderRole(t *testing.T) {
	svc := newTestService()

	result, err := svc.Register(context.Background(), "provider@example.com", "password123", "商家用户", "provider")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}
	if result.User.Role != "provider" {
		t.Errorf("Register() role = %s, want provider", result.User.Role)
	}

	loginResult, err := svc.Login(context.Background(), "provider@example.com", "password123")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if loginResult.User.Role != "provider" {
		t.Errorf("Login() role = %s, want provider", loginResult.User.Role)
	}
}

func TestAuthService_Register_InvalidRole(t *testing.T) {
	svc := newTestService()

	_, err := svc.Register(context.Background(), "bad-role@example.com", "password123", "用户", "unknown")
	if !errors.Is(err, ErrInvalidRole) {
		t.Errorf("Register() error = %v, want %v", err, ErrInvalidRole)
	}
}

func TestAuthService_Register_InvalidEmail(t *testing.T) {
	svc := newTestService()
	tests := []struct {
		name  string
		email string
	}{
		{"empty", ""},
		{"no_at", "invalid"},
		{"no_domain", "user@"},
		{"no_tld", "user@domain"},
		{"spaces", "  "},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Register(context.Background(), tt.email, "password123", "用户", "")
			if !errors.Is(err, ErrInvalidEmail) {
				t.Errorf("Register() error = %v, want %v", err, ErrInvalidEmail)
			}
		})
	}
}

func TestAuthService_Register_WeakPassword(t *testing.T) {
	svc := newTestService()
	tests := []struct {
		name     string
		password string
	}{
		{"empty", ""},
		{"too_short", "abc"},
		{"seven_chars", "1234567"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Register(context.Background(), "test@example.com", tt.password, "用户", "")
			if !errors.Is(err, ErrWeakPassword) {
				t.Errorf("Register() error = %v, want %v", err, ErrWeakPassword)
			}
		})
	}
}

func TestAuthService_Register_EmptyName(t *testing.T) {
	svc := newTestService()
	_, err := svc.Register(context.Background(), "test@example.com", "password123", "", "")
	if !errors.Is(err, ErrNameRequired) {
		t.Errorf("Register() error = %v, want %v", err, ErrNameRequired)
	}
}

func TestAuthService_Register_DuplicateEmail(t *testing.T) {
	svc := newTestService()

	_, err := svc.Register(context.Background(), "dup@example.com", "password123", "用户A", "")
	if err != nil {
		t.Fatalf("first Register() error = %v", err)
	}

	_, err = svc.Register(context.Background(), "dup@example.com", "password456", "用户B", "")
	if !errors.Is(err, ErrEmailAlreadyRegistered) {
		t.Errorf("Register() error = %v, want %v", err, ErrEmailAlreadyRegistered)
	}
}

func TestAuthService_Register_EmailCaseInsensitive(t *testing.T) {
	svc := newTestService()

	_, err := svc.Register(context.Background(), "Case@Test.com", "password123", "用户A", "")
	if err != nil {
		t.Fatalf("first Register() error = %v", err)
	}

	_, err = svc.Register(context.Background(), "case@test.com", "password456", "用户B", "")
	if !errors.Is(err, ErrEmailAlreadyRegistered) {
		t.Errorf("Register() error = %v, want %v", err, ErrEmailAlreadyRegistered)
	}
}

func TestAuthService_Login_Success(t *testing.T) {
	svc := newTestService()

	_, err := svc.Register(context.Background(), "login@test.com", "password123", "用户", "")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	result, err := svc.Login(context.Background(), "login@test.com", "password123")
	if err != nil {
		t.Fatalf("Login() error = %v", err)
	}
	if result.User == nil {
		t.Fatal("Login() result.User is nil")
	}
	if result.User.Email != "login@test.com" {
		t.Errorf("Login() email = %s, want %s", result.User.Email, "login@test.com")
	}
	if result.AccessToken == "" {
		t.Error("Login() AccessToken should not be empty")
	}
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	svc := newTestService()

	_, err := svc.Register(context.Background(), "login@test.com", "password123", "用户", "")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	_, err = svc.Login(context.Background(), "login@test.com", "wrong-password")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("Login() error = %v, want %v", err, ErrInvalidCredentials)
	}
}

func TestAuthService_Login_NonExistentEmail(t *testing.T) {
	svc := newTestService()

	_, err := svc.Login(context.Background(), "nobody@test.com", "password123")
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Errorf("Login() error = %v, want %v", err, ErrInvalidCredentials)
	}
}

func TestAuthService_Login_EmptyInput(t *testing.T) {
	svc := newTestService()
	tests := []struct {
		name     string
		email    string
		password string
	}{
		{"empty_email", "", "password123"},
		{"empty_password", "test@test.com", ""},
		{"both_empty", "", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Login(context.Background(), tt.email, tt.password)
			if !errors.Is(err, ErrInvalidCredentials) {
				t.Errorf("Login() error = %v, want %v", err, ErrInvalidCredentials)
			}
		})
	}
}

func TestAuthService_Login_CaseInsensitiveEmail(t *testing.T) {
	svc := newTestService()

	_, err := svc.Register(context.Background(), "user@test.com", "password123", "用户", "")
	if err != nil {
		t.Fatalf("Register() error = %v", err)
	}

	result, err := svc.Login(context.Background(), "USER@test.com", "password123")
	if err != nil {
		t.Errorf("Login() error = %v, want nil (case insensitive)", err)
	}
	if result == nil || result.User == nil {
		t.Fatal("Login() result is nil")
	}
}
