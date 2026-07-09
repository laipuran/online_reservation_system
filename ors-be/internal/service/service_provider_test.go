package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"ors-be/internal/model"
)

type mockServiceProviderRepo struct {
	byID     map[int64]*model.ServiceProvider
	byUserID map[int64]*model.ServiceProvider
	nextID   int64
}

func newMockServiceProviderRepo() *mockServiceProviderRepo {
	return &mockServiceProviderRepo{
		byID:     make(map[int64]*model.ServiceProvider),
		byUserID: make(map[int64]*model.ServiceProvider),
		nextID:   1,
	}
}

func (m *mockServiceProviderRepo) Create(ctx context.Context, provider *model.ServiceProvider) error {
	if _, exists := m.byUserID[provider.UserID]; exists {
		return errors.New("duplicate user_id")
	}
	provider.ID = m.nextID
	provider.CreatedAt = time.Now()
	provider.UpdatedAt = provider.CreatedAt
	m.nextID++
	m.byID[provider.ID] = provider
	m.byUserID[provider.UserID] = provider
	return nil
}

func (m *mockServiceProviderRepo) GetByID(ctx context.Context, id int64) (*model.ServiceProvider, error) {
	return m.byID[id], nil
}

func (m *mockServiceProviderRepo) GetByUserID(ctx context.Context, userID int64) (*model.ServiceProvider, error) {
	return m.byUserID[userID], nil
}

func (m *mockServiceProviderRepo) Update(ctx context.Context, provider *model.ServiceProvider) error {
	if _, exists := m.byID[provider.ID]; !exists {
		return errors.New("not found")
	}
	provider.UpdatedAt = time.Now()
	m.byID[provider.ID] = provider
	m.byUserID[provider.UserID] = provider
	return nil
}

func newTestServiceProviderService() ServiceProviderService {
	return NewServiceProviderService(newMockServiceProviderRepo())
}

func TestServiceProviderService_Create_Success(t *testing.T) {
	svc := newTestServiceProviderService()

	provider, err := svc.Create(context.Background(), 1, ServiceProviderInput{
		BusinessName: " 舒心养生馆 ",
		Description:  " 专业按摩 ",
		Address:      " 上海市 ",
		Phone:        " 13800000000 ",
		Email:        " SHOP@EXAMPLE.COM ",
		LogoURL:      " https://example.com/logo.png ",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if provider.ID == 0 {
		t.Error("Create() ID should be non-zero")
	}
	if provider.UserID != 1 {
		t.Errorf("Create() userID = %d, want 1", provider.UserID)
	}
	if provider.BusinessName != "舒心养生馆" {
		t.Errorf("Create() businessName = %q", provider.BusinessName)
	}
	if provider.Email != "shop@example.com" {
		t.Errorf("Create() email = %q", provider.Email)
	}
}

func TestServiceProviderService_Create_EmptyBusinessName(t *testing.T) {
	svc := newTestServiceProviderService()

	_, err := svc.Create(context.Background(), 1, ServiceProviderInput{BusinessName: " "})
	if !errors.Is(err, ErrBusinessNameRequired) {
		t.Errorf("Create() error = %v, want %v", err, ErrBusinessNameRequired)
	}
}

func TestServiceProviderService_Create_InvalidEmail(t *testing.T) {
	svc := newTestServiceProviderService()

	_, err := svc.Create(context.Background(), 1, ServiceProviderInput{
		BusinessName: "商家A",
		Email:        "bad-email",
	})
	if !errors.Is(err, ErrInvalidEmail) {
		t.Errorf("Create() error = %v, want %v", err, ErrInvalidEmail)
	}
}

func TestServiceProviderService_Create_DuplicateUser(t *testing.T) {
	svc := newTestServiceProviderService()

	_, err := svc.Create(context.Background(), 1, ServiceProviderInput{BusinessName: "商家A"})
	if err != nil {
		t.Fatalf("first Create() error = %v", err)
	}

	_, err = svc.Create(context.Background(), 1, ServiceProviderInput{BusinessName: "商家B"})
	if !errors.Is(err, ErrProviderAlreadyExists) {
		t.Errorf("Create() error = %v, want %v", err, ErrProviderAlreadyExists)
	}
}

func TestServiceProviderService_GetMine_NotFound(t *testing.T) {
	svc := newTestServiceProviderService()

	_, err := svc.GetMine(context.Background(), 999)
	if !errors.Is(err, ErrProviderNotFound) {
		t.Errorf("GetMine() error = %v, want %v", err, ErrProviderNotFound)
	}
}

func TestServiceProviderService_UpdateMine_Success(t *testing.T) {
	svc := newTestServiceProviderService()

	_, err := svc.Create(context.Background(), 1, ServiceProviderInput{BusinessName: "旧商家"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	updated, err := svc.UpdateMine(context.Background(), 1, ServiceProviderInput{
		BusinessName: "新商家",
		Description:  "新简介",
	})
	if err != nil {
		t.Fatalf("UpdateMine() error = %v", err)
	}
	if updated.BusinessName != "新商家" {
		t.Errorf("UpdateMine() businessName = %q", updated.BusinessName)
	}
	if updated.Description != "新简介" {
		t.Errorf("UpdateMine() description = %q", updated.Description)
	}
}

func TestServiceProviderService_UpdateMine_InvalidEmail(t *testing.T) {
	svc := newTestServiceProviderService()

	_, err := svc.Create(context.Background(), 1, ServiceProviderInput{BusinessName: "旧商家"})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	_, err = svc.UpdateMine(context.Background(), 1, ServiceProviderInput{
		BusinessName: "新商家",
		Email:        "bad-email",
	})
	if !errors.Is(err, ErrInvalidEmail) {
		t.Errorf("UpdateMine() error = %v, want %v", err, ErrInvalidEmail)
	}
}

func TestServiceProviderService_GetByID_NotFound(t *testing.T) {
	svc := newTestServiceProviderService()

	_, err := svc.GetByID(context.Background(), 999)
	if !errors.Is(err, ErrProviderNotFound) {
		t.Errorf("GetByID() error = %v, want %v", err, ErrProviderNotFound)
	}
}
