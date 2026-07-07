package service

import (
	"context"
	"errors"
	"regexp"
	"strings"

	"ors-be/internal/auth"
	"ors-be/internal/model"
	"ors-be/internal/repository"
)

var (
	ErrEmailAlreadyRegistered = errors.New("邮箱已注册")
	ErrInvalidCredentials     = errors.New("邮箱或密码错误")
	ErrInvalidEmail           = errors.New("邮箱格式不正确")
	ErrWeakPassword           = errors.New("密码长度至少8位")
	ErrNameRequired           = errors.New("昵称不能为空")
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)

type RegisterResult struct {
	User        *model.User `json:"user"`
	AccessToken string      `json:"access_token"`
}

type LoginResult struct {
	User        *model.User `json:"user"`
	AccessToken string      `json:"access_token"`
}

type AuthService interface {
	Register(ctx context.Context, email, password, name string) (*RegisterResult, error)
	Login(ctx context.Context, email, password string) (*LoginResult, error)
}

type authService struct {
	userRepo repository.UserRepository
	hasher   auth.Hasher
	tokenGen auth.TokenGenerator
}

func NewAuthService(userRepo repository.UserRepository, hasher auth.Hasher, tokenGen auth.TokenGenerator) AuthService {
	return &authService{
		userRepo: userRepo,
		hasher:   hasher,
		tokenGen: tokenGen,
	}
}

func (s *authService) Register(ctx context.Context, email, password, name string) (*RegisterResult, error) {
	email = strings.TrimSpace(strings.ToLower(email))
	name = strings.TrimSpace(name)

	if email == "" || !emailRegex.MatchString(email) {
		return nil, ErrInvalidEmail
	}
	if len(password) < 8 {
		return nil, ErrWeakPassword
	}
	if name == "" {
		return nil, ErrNameRequired
	}

	existing, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrEmailAlreadyRegistered
	}

	hash, err := s.hasher.Hash(password)
	if err != nil {
		return nil, err
	}

	user := &model.User{
		Name:         name,
		Email:        email,
		PasswordHash: hash,
		Role:         "customer",
	}

	if err := s.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	token, err := s.tokenGen.Generate(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &RegisterResult{User: user, AccessToken: token}, nil
}

func (s *authService) Login(ctx context.Context, email, password string) (*LoginResult, error) {
	email = strings.TrimSpace(strings.ToLower(email))

	if email == "" {
		return nil, ErrInvalidCredentials
	}
	if password == "" {
		return nil, ErrInvalidCredentials
	}

	user, err := s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrInvalidCredentials
	}

	if !s.hasher.Verify(password, user.PasswordHash) {
		return nil, ErrInvalidCredentials
	}

	token, err := s.tokenGen.Generate(user.ID, user.Role)
	if err != nil {
		return nil, err
	}

	return &LoginResult{User: user, AccessToken: token}, nil
}
