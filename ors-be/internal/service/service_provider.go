package service

import (
	"context"
	"errors"
	"strings"

	"ors-be/internal/model"
	"ors-be/internal/repository"
)

var (
	ErrProviderNotFound      = errors.New("服务提供者不存在")
	ErrProviderAlreadyExists = errors.New("服务提供者资料已存在")
	ErrBusinessNameRequired  = errors.New("商家名称不能为空")
	ErrProviderForbidden     = errors.New("无权操作该服务提供者资料")
)

type ServiceProviderService interface {
	Create(ctx context.Context, userID int64, input ServiceProviderInput) (*model.ServiceProvider, error)
	GetByID(ctx context.Context, id int64) (*model.ServiceProvider, error)
	GetMine(ctx context.Context, userID int64) (*model.ServiceProvider, error)
	UpdateMine(ctx context.Context, userID int64, input ServiceProviderInput) (*model.ServiceProvider, error)
}

type ServiceProviderInput struct {
	BusinessName string `json:"business_name"`
	Description  string `json:"description"`
	Address      string `json:"address"`
	Phone        string `json:"phone"`
	Email        string `json:"email"`
	LogoURL      string `json:"logo_url"`
}

type serviceProviderService struct {
	providerRepo repository.ServiceProviderRepository
}

func NewServiceProviderService(providerRepo repository.ServiceProviderRepository) ServiceProviderService {
	return &serviceProviderService{providerRepo: providerRepo}
}

func (s *serviceProviderService) Create(ctx context.Context, userID int64, input ServiceProviderInput) (*model.ServiceProvider, error) {
	provider := normalizeProviderInput(userID, input)
	if provider.BusinessName == "" {
		return nil, ErrBusinessNameRequired
	}
	if provider.Email != "" && !emailRegex.MatchString(provider.Email) {
		return nil, ErrInvalidEmail
	}
	if provider.Phone != "" && !isPhoneValid(provider.Phone) {
		return nil, ErrInvalidPhone
	}

	existing, err := s.providerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrProviderAlreadyExists
	}

	if err := s.providerRepo.Create(ctx, provider); err != nil {
		return nil, err
	}
	return provider, nil
}

func (s *serviceProviderService) GetByID(ctx context.Context, id int64) (*model.ServiceProvider, error) {
	provider, err := s.providerRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if provider == nil {
		return nil, ErrProviderNotFound
	}
	return provider, nil
}

func (s *serviceProviderService) GetMine(ctx context.Context, userID int64) (*model.ServiceProvider, error) {
	provider, err := s.providerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if provider == nil {
		return nil, ErrProviderNotFound
	}
	return provider, nil
}

func (s *serviceProviderService) UpdateMine(ctx context.Context, userID int64, input ServiceProviderInput) (*model.ServiceProvider, error) {
	existing, err := s.providerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrProviderNotFound
	}

	updated := normalizeProviderInput(userID, input)
	if updated.BusinessName == "" {
		return nil, ErrBusinessNameRequired
	}
	if updated.Email != "" && !emailRegex.MatchString(updated.Email) {
		return nil, ErrInvalidEmail
	}
	if updated.Phone != "" && !isPhoneValid(updated.Phone) {
		return nil, ErrInvalidPhone
	}

	existing.BusinessName = updated.BusinessName
	existing.Description = updated.Description
	existing.Address = updated.Address
	existing.Phone = updated.Phone
	existing.Email = updated.Email
	existing.LogoURL = updated.LogoURL

	if err := s.providerRepo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func normalizeProviderInput(userID int64, input ServiceProviderInput) *model.ServiceProvider {
	return &model.ServiceProvider{
		UserID:       userID,
		BusinessName: strings.TrimSpace(input.BusinessName),
		Description:  strings.TrimSpace(input.Description),
		Address:      strings.TrimSpace(input.Address),
		Phone:        strings.TrimSpace(input.Phone),
		Email:        strings.TrimSpace(strings.ToLower(input.Email)),
		LogoURL:      strings.TrimSpace(input.LogoURL),
	}
}
