package service

import (
	"context"
	"errors"
	"strings"

	"ors-be/internal/model"
	"ors-be/internal/repository"
)

var (
	ErrServiceNotFound        = errors.New("服务不存在")
	ErrServiceTitleRequired   = errors.New("服务标题不能为空")
	ErrServiceInvalidCategory = errors.New("服务分类不能为空")
	ErrServiceInvalidPrice    = errors.New("服务价格不能小于0")
	ErrServiceInvalidDuration = errors.New("服务时长必须大于0")
	ErrServiceInvalidStatus   = errors.New("服务状态不正确")
	ErrServiceForbidden       = errors.New("无权操作该服务")
	ErrServiceInvalidTag      = errors.New("标签ID不正确")
)

type ServiceInput struct {
	CategoryID      int64   `json:"category_id"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	Price           float64 `json:"price"`
	DurationMinutes int     `json:"duration_minutes"`
	ImageURL        string  `json:"image_url"`
}

type ServiceStatusInput struct {
	Status string `json:"status"`
}

type ServiceTagsInput struct {
	TagIDs []int64 `json:"tag_ids"`
}

type ServiceListResult struct {
	Items    []*model.ServiceView `json:"items"`
	Total    int                  `json:"total"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"page_size"`
}

type ServiceService interface {
	Create(ctx context.Context, userID int64, input ServiceInput) (*model.ServiceView, error)
	GetByID(ctx context.Context, id int64) (*model.ServiceView, error)
	List(ctx context.Context, filter model.ServiceFilter) (*ServiceListResult, error)
	Update(ctx context.Context, userID, id int64, input ServiceInput) (*model.ServiceView, error)
	UpdateStatus(ctx context.Context, userID, id int64, status string) (*model.ServiceView, error)
	ListTags(ctx context.Context, id int64) ([]*model.Tag, error)
	ReplaceTags(ctx context.Context, userID, id int64, input ServiceTagsInput) ([]*model.Tag, error)
}

type serviceService struct {
	serviceRepo    repository.ServiceRepository
	providerRepo   repository.ServiceProviderRepository
	tagRepo        repository.TagRepository
	serviceTagRepo repository.ServiceTagRepository
}

func NewServiceService(
	serviceRepo repository.ServiceRepository,
	providerRepo repository.ServiceProviderRepository,
	tagRepo repository.TagRepository,
	serviceTagRepo repository.ServiceTagRepository,
) ServiceService {
	return &serviceService{
		serviceRepo:    serviceRepo,
		providerRepo:   providerRepo,
		tagRepo:        tagRepo,
		serviceTagRepo: serviceTagRepo,
	}
}

func (s *serviceService) Create(ctx context.Context, userID int64, input ServiceInput) (*model.ServiceView, error) {
	provider, err := s.providerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if provider == nil {
		return nil, ErrProviderNotFound
	}

	service, err := normalizeServiceInput(input)
	if err != nil {
		return nil, err
	}
	service.ProviderID = provider.ID
	service.Status = "active"

	if err := s.serviceRepo.Create(ctx, service); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, service.ID)
}

func (s *serviceService) GetByID(ctx context.Context, id int64) (*model.ServiceView, error) {
	service, err := s.serviceRepo.GetViewByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if service == nil {
		return nil, ErrServiceNotFound
	}
	return service, nil
}

func (s *serviceService) List(ctx context.Context, filter model.ServiceFilter) (*ServiceListResult, error) {
	filter = normalizeServiceFilter(filter)
	items, total, err := s.serviceRepo.List(ctx, filter)
	if err != nil {
		return nil, err
	}
	return &ServiceListResult{
		Items:    items,
		Total:    total,
		Page:     filter.Page,
		PageSize: filter.PageSize,
	}, nil
}

func (s *serviceService) Update(ctx context.Context, userID, id int64, input ServiceInput) (*model.ServiceView, error) {
	existing, err := s.authorizeServiceOwner(ctx, userID, id)
	if err != nil {
		return nil, err
	}

	updated, err := normalizeServiceInput(input)
	if err != nil {
		return nil, err
	}

	existing.CategoryID = updated.CategoryID
	existing.Title = updated.Title
	existing.Description = updated.Description
	existing.Price = updated.Price
	existing.DurationMinutes = updated.DurationMinutes
	existing.ImageURL = updated.ImageURL

	if err := s.serviceRepo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

func (s *serviceService) UpdateStatus(ctx context.Context, userID, id int64, status string) (*model.ServiceView, error) {
	if _, err := s.authorizeServiceOwner(ctx, userID, id); err != nil {
		return nil, err
	}

	status = strings.TrimSpace(strings.ToLower(status))
	if status != "active" && status != "inactive" {
		return nil, ErrServiceInvalidStatus
	}

	if err := s.serviceRepo.UpdateStatus(ctx, id, status); err != nil {
		return nil, err
	}
	return s.GetByID(ctx, id)
}

func (s *serviceService) ListTags(ctx context.Context, id int64) ([]*model.Tag, error) {
	existing, err := s.serviceRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrServiceNotFound
	}
	return s.serviceTagRepo.ListByServiceID(ctx, id)
}

func (s *serviceService) ReplaceTags(ctx context.Context, userID, id int64, input ServiceTagsInput) ([]*model.Tag, error) {
	if _, err := s.authorizeServiceOwner(ctx, userID, id); err != nil {
		return nil, err
	}

	tagIDs, err := s.normalizeTagIDs(ctx, input.TagIDs)
	if err != nil {
		return nil, err
	}

	if err := s.serviceTagRepo.ReplaceByServiceID(ctx, id, tagIDs); err != nil {
		return nil, err
	}
	return s.serviceTagRepo.ListByServiceID(ctx, id)
}

func (s *serviceService) authorizeServiceOwner(ctx context.Context, userID, serviceID int64) (*model.Service, error) {
	provider, err := s.providerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if provider == nil {
		return nil, ErrProviderNotFound
	}

	existing, err := s.serviceRepo.GetByID(ctx, serviceID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrServiceNotFound
	}
	if existing.ProviderID != provider.ID {
		return nil, ErrServiceForbidden
	}
	return existing, nil
}

func (s *serviceService) normalizeTagIDs(ctx context.Context, tagIDs []int64) ([]int64, error) {
	seen := make(map[int64]struct{}, len(tagIDs))
	normalized := make([]int64, 0, len(tagIDs))

	for _, tagID := range tagIDs {
		if tagID <= 0 {
			return nil, ErrServiceInvalidTag
		}
		if _, ok := seen[tagID]; ok {
			continue
		}

		tag, err := s.tagRepo.GetByID(ctx, tagID)
		if err != nil {
			return nil, err
		}
		if tag == nil {
			return nil, ErrTagNotFound
		}

		seen[tagID] = struct{}{}
		normalized = append(normalized, tagID)
	}
	return normalized, nil
}

func normalizeServiceInput(input ServiceInput) (*model.Service, error) {
	service := &model.Service{
		CategoryID:      input.CategoryID,
		Title:           strings.TrimSpace(input.Title),
		Description:     strings.TrimSpace(input.Description),
		Price:           input.Price,
		DurationMinutes: input.DurationMinutes,
		ImageURL:        strings.TrimSpace(input.ImageURL),
	}

	if service.Title == "" {
		return nil, ErrServiceTitleRequired
	}
	if service.CategoryID <= 0 {
		return nil, ErrServiceInvalidCategory
	}
	if service.Price < 0 {
		return nil, ErrServiceInvalidPrice
	}
	if service.DurationMinutes <= 0 {
		return nil, ErrServiceInvalidDuration
	}
	return service, nil
}

func normalizeServiceFilter(filter model.ServiceFilter) model.ServiceFilter {
	filter.Keyword = strings.TrimSpace(filter.Keyword)
	filter.SortBy = strings.TrimSpace(strings.ToLower(filter.SortBy))
	filter.SortOrder = strings.TrimSpace(strings.ToLower(filter.SortOrder))
	filter.Status = strings.TrimSpace(strings.ToLower(filter.Status))

	if filter.Page <= 0 {
		filter.Page = 1
	}
	if filter.PageSize <= 0 {
		filter.PageSize = 20
	}
	if filter.PageSize > 50 {
		filter.PageSize = 50
	}
	if filter.SortBy != "price" && filter.SortBy != "rating" && filter.SortBy != "created_at" {
		filter.SortBy = "created_at"
	}
	if filter.SortOrder != "asc" {
		filter.SortOrder = "desc"
	}
	return filter
}
