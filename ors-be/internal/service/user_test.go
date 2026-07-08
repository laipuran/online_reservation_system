package service

import (
	"context"
	"errors"
	"testing"

	"ors-be/internal/auth"
	"ors-be/internal/model"
)

func newTestUserService() (UserService, *mockUserRepo) {
	repo := newMockUserRepo()
	hasher := auth.NewHasher()
	hash, _ := hasher.Hash("password123")
	_ = repo.Create(context.Background(), &model.User{
		Name:         "测试用户",
		Email:        "user@example.com",
		PasswordHash: hash,
		Role:         "customer",
	})
	return NewUserService(repo, hasher), repo
}

func TestUserService_GetMine_Success(t *testing.T) {
	svc, _ := newTestUserService()

	user, err := svc.GetMine(context.Background(), 1)
	if err != nil {
		t.Fatalf("GetMine() error = %v", err)
	}
	if user.Email != "user@example.com" {
		t.Errorf("GetMine() email = %s", user.Email)
	}
}

func TestUserService_UpdateMine_Success(t *testing.T) {
	svc, _ := newTestUserService()

	user, err := svc.UpdateMine(context.Background(), 1, UserInput{
		Name:      " 新昵称 ",
		Phone:     " 13800000000 ",
		AvatarURL: " https://example.com/avatar.png ",
	})
	if err != nil {
		t.Fatalf("UpdateMine() error = %v", err)
	}
	if user.Name != "新昵称" {
		t.Errorf("UpdateMine() name = %q", user.Name)
	}
	if user.Phone != "13800000000" {
		t.Errorf("UpdateMine() phone = %q", user.Phone)
	}
	if user.AvatarURL != "https://example.com/avatar.png" {
		t.Errorf("UpdateMine() avatarURL = %q", user.AvatarURL)
	}
}

func TestUserService_UpdateMine_EmptyName(t *testing.T) {
	svc, _ := newTestUserService()

	_, err := svc.UpdateMine(context.Background(), 1, UserInput{Name: " "})
	if !errors.Is(err, ErrNameRequired) {
		t.Errorf("UpdateMine() error = %v, want %v", err, ErrNameRequired)
	}
}

func TestUserService_UpdatePassword_Success(t *testing.T) {
	svc, repo := newTestUserService()

	err := svc.UpdatePassword(context.Background(), 1, UserPasswordInput{
		CurrentPassword: "password123",
		NewPassword:     "newpass123",
	})
	if err != nil {
		t.Fatalf("UpdatePassword() error = %v", err)
	}

	user, _ := repo.GetByID(context.Background(), 1)
	if !auth.NewHasher().Verify("newpass123", user.PasswordHash) {
		t.Error("UpdatePassword() did not update password hash")
	}
}

func TestUserService_UpdatePassword_WrongCurrentPassword(t *testing.T) {
	svc, _ := newTestUserService()

	err := svc.UpdatePassword(context.Background(), 1, UserPasswordInput{
		CurrentPassword: "wrongpass",
		NewPassword:     "newpass123",
	})
	if !errors.Is(err, ErrCurrentPasswordWrong) {
		t.Errorf("UpdatePassword() error = %v, want %v", err, ErrCurrentPasswordWrong)
	}
}
