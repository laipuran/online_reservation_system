package service

import (
	"context"
	"errors"
	"strings"

	"ors-be/internal/auth"
	"ors-be/internal/model"
	"ors-be/internal/repository"
)

var (
	ErrUserNotFound         = errors.New("用户不存在")
	ErrCurrentPasswordWrong = errors.New("当前密码错误")
)

type UserInput struct {
	Name      string `json:"name"`
	Phone     string `json:"phone"`
	AvatarURL string `json:"avatar_url"`
}

type UserPasswordInput struct {
	CurrentPassword string `json:"current_password"`
	NewPassword     string `json:"new_password"`
}

type UserService interface {
	GetByID(ctx context.Context, id int64) (*model.User, error)
	GetMine(ctx context.Context, userID int64) (*model.User, error)
	UpdateMine(ctx context.Context, userID int64, input UserInput) (*model.User, error)
	UpdatePassword(ctx context.Context, userID int64, input UserPasswordInput) error
}

type userService struct {
	userRepo repository.UserRepository
	hasher   auth.Hasher
}

func NewUserService(userRepo repository.UserRepository, hasher auth.Hasher) UserService {
	return &userService{userRepo: userRepo, hasher: hasher}
}

func (s *userService) GetByID(ctx context.Context, id int64) (*model.User, error) {
	user, err := s.userRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}
	return user, nil
}

func (s *userService) GetMine(ctx context.Context, userID int64) (*model.User, error) {
	return s.GetByID(ctx, userID)
}

func (s *userService) UpdateMine(ctx context.Context, userID int64, input UserInput) (*model.User, error) {
	user, err := s.GetMine(ctx, userID)
	if err != nil {
		return nil, err
	}

	user.Name = strings.TrimSpace(input.Name)
	user.Phone = strings.TrimSpace(input.Phone)
	user.AvatarURL = strings.TrimSpace(input.AvatarURL)
	if user.Name == "" {
		return nil, ErrNameRequired
	}
	if user.Phone != "" && !isPhoneValid(user.Phone) {
		return nil, ErrInvalidPhone
	}

	if err := s.userRepo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *userService) UpdatePassword(ctx context.Context, userID int64, input UserPasswordInput) error {
	user, err := s.GetMine(ctx, userID)
	if err != nil {
		return err
	}

	if !s.hasher.Verify(input.CurrentPassword, user.PasswordHash) {
		return ErrCurrentPasswordWrong
	}
	if len(input.NewPassword) < 8 {
		return ErrWeakPassword
	}

	hash, err := s.hasher.Hash(input.NewPassword)
	if err != nil {
		return err
	}
	return s.userRepo.UpdatePassword(ctx, userID, hash)
}
