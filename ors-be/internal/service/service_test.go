package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"ors-be/internal/model"
)

type mockServiceRepo struct {
	services map[int64]*model.Service
	nextID   int64
}

type mockServiceTagRepo struct {
	tagsByService map[int64][]*model.Tag
}

func newMockServiceRepo() *mockServiceRepo {
	return &mockServiceRepo{
		services: make(map[int64]*model.Service),
		nextID:   1,
	}
}

func newMockServiceTagRepo() *mockServiceTagRepo {
	return &mockServiceTagRepo{tagsByService: make(map[int64][]*model.Tag)}
}

func (m *mockServiceRepo) Create(ctx context.Context, service *model.Service) error {
	service.ID = m.nextID
	service.AvgRating = 0
	service.ReviewCount = 0
	service.CreatedAt = time.Now()
	service.UpdatedAt = service.CreatedAt
	m.nextID++
	m.services[service.ID] = cloneService(service)
	return nil
}

func (m *mockServiceRepo) GetByID(ctx context.Context, id int64) (*model.Service, error) {
	service := m.services[id]
	if service == nil {
		return nil, nil
	}
	return cloneService(service), nil
}

func (m *mockServiceRepo) GetViewByID(ctx context.Context, id int64) (*model.ServiceView, error) {
	service := m.services[id]
	if service == nil {
		return nil, nil
	}
	return serviceToView(service), nil
}

func (m *mockServiceRepo) List(ctx context.Context, filter model.ServiceFilter) ([]*model.ServiceView, int, error) {
	items := make([]*model.ServiceView, 0)
	for _, service := range m.services {
		if filter.ProviderID != nil && service.ProviderID != *filter.ProviderID {
			continue
		}
		if filter.Status != "" && service.Status != filter.Status {
			continue
		}
		items = append(items, serviceToView(service))
	}
	return items, len(items), nil
}

func (m *mockServiceRepo) Update(ctx context.Context, service *model.Service) error {
	if m.services[service.ID] == nil {
		return errors.New("service not found")
	}
	service.UpdatedAt = time.Now()
	m.services[service.ID] = cloneService(service)
	return nil
}

func (m *mockServiceRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	service := m.services[id]
	if service == nil {
		return errors.New("service not found")
	}
	service.Status = status
	service.UpdatedAt = time.Now()
	return nil
}

func (m *mockServiceTagRepo) ReplaceByServiceID(ctx context.Context, serviceID int64, tagIDs []int64) error {
	tags := make([]*model.Tag, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		tags = append(tags, &model.Tag{ID: tagID, Name: "标签"})
	}
	m.tagsByService[serviceID] = tags
	return nil
}

func (m *mockServiceTagRepo) ListByServiceID(ctx context.Context, serviceID int64) ([]*model.Tag, error) {
	tags := make([]*model.Tag, 0, len(m.tagsByService[serviceID]))
	for _, tag := range m.tagsByService[serviceID] {
		tags = append(tags, cloneTag(tag))
	}
	return tags, nil
}

func newTestBusinessService() ServiceService {
	providerRepo := newMockServiceProviderRepo()
	_ = providerRepo.Create(context.Background(), &model.ServiceProvider{
		UserID:       1,
		BusinessName: "舒心养生馆",
	})
	tagRepo := newMockTagRepo()
	_ = tagRepo.Create(context.Background(), &model.Tag{Name: "放松"})
	_ = tagRepo.Create(context.Background(), &model.Tag{Name: "塑形"})
	return NewServiceService(newMockServiceRepo(), providerRepo, tagRepo, newMockServiceTagRepo())
}

func TestServiceService_Create_Success(t *testing.T) {
	svc := newTestBusinessService()

	result, err := svc.Create(context.Background(), 1, ServiceInput{
		CategoryID:      1,
		Title:           " 肩颈按摩 ",
		Description:     " 放松 ",
		Price:           199,
		DurationMinutes: 60,
		ImageURL:        " https://example.com/service.png ",
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if result.ID == 0 {
		t.Error("Create() ID should be non-zero")
	}
	if result.Title != "肩颈按摩" {
		t.Errorf("Create() title = %q", result.Title)
	}
	if result.Status != "active" {
		t.Errorf("Create() status = %q, want active", result.Status)
	}
	if result.Provider.ID != 1 {
		t.Errorf("Create() provider ID = %d, want 1", result.Provider.ID)
	}
}

func TestServiceService_Create_InvalidInput(t *testing.T) {
	svc := newTestBusinessService()

	tests := []struct {
		name  string
		input ServiceInput
		want  error
	}{
		{"empty title", ServiceInput{CategoryID: 1, Price: 1, DurationMinutes: 1}, ErrServiceTitleRequired},
		{"empty category", ServiceInput{Title: "服务", Price: 1, DurationMinutes: 1}, ErrServiceInvalidCategory},
		{"negative price", ServiceInput{CategoryID: 1, Title: "服务", Price: -1, DurationMinutes: 1}, ErrServiceInvalidPrice},
		{"invalid duration", ServiceInput{CategoryID: 1, Title: "服务", Price: 1}, ErrServiceInvalidDuration},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Create(context.Background(), 1, tt.input)
			if !errors.Is(err, tt.want) {
				t.Errorf("Create() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestServiceService_Create_ProviderNotFound(t *testing.T) {
	svc := newTestBusinessService()

	_, err := svc.Create(context.Background(), 99, ServiceInput{
		CategoryID:      1,
		Title:           "服务",
		Price:           1,
		DurationMinutes: 1,
	})
	if !errors.Is(err, ErrProviderNotFound) {
		t.Errorf("Create() error = %v, want %v", err, ErrProviderNotFound)
	}
}

func TestServiceService_Update_ForbidsOtherProvider(t *testing.T) {
	serviceRepo := newMockServiceRepo()
	providerRepo := newMockServiceProviderRepo()
	_ = providerRepo.Create(context.Background(), &model.ServiceProvider{UserID: 1, BusinessName: "商家A"})
	_ = providerRepo.Create(context.Background(), &model.ServiceProvider{UserID: 2, BusinessName: "商家B"})
	svc := NewServiceService(serviceRepo, providerRepo, newMockTagRepo(), newMockServiceTagRepo())

	created, err := svc.Create(context.Background(), 1, ServiceInput{
		CategoryID:      1,
		Title:           "服务",
		Price:           1,
		DurationMinutes: 1,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	_, err = svc.Update(context.Background(), 2, created.ID, ServiceInput{
		CategoryID:      1,
		Title:           "新服务",
		Price:           2,
		DurationMinutes: 2,
	})
	if !errors.Is(err, ErrServiceForbidden) {
		t.Errorf("Update() error = %v, want %v", err, ErrServiceForbidden)
	}
}

func TestServiceService_UpdateStatus_Success(t *testing.T) {
	svc := newTestBusinessService()

	created, err := svc.Create(context.Background(), 1, ServiceInput{
		CategoryID:      1,
		Title:           "服务",
		Price:           1,
		DurationMinutes: 1,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	updated, err := svc.UpdateStatus(context.Background(), 1, created.ID, "inactive")
	if err != nil {
		t.Fatalf("UpdateStatus() error = %v", err)
	}
	if updated.Status != "inactive" {
		t.Errorf("UpdateStatus() status = %q, want inactive", updated.Status)
	}
}

func TestServiceService_ReplaceTags_Success(t *testing.T) {
	svc := newTestBusinessService()

	created, err := svc.Create(context.Background(), 1, ServiceInput{
		CategoryID:      1,
		Title:           "服务",
		Price:           1,
		DurationMinutes: 1,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	tags, err := svc.ReplaceTags(context.Background(), 1, created.ID, ServiceTagsInput{TagIDs: []int64{1, 2, 1}})
	if err != nil {
		t.Fatalf("ReplaceTags() error = %v", err)
	}
	if len(tags) != 2 {
		t.Fatalf("ReplaceTags() len = %d, want 2", len(tags))
	}

	listed, err := svc.ListTags(context.Background(), created.ID)
	if err != nil {
		t.Fatalf("ListTags() error = %v", err)
	}
	if len(listed) != 2 {
		t.Errorf("ListTags() len = %d, want 2", len(listed))
	}
}

func TestServiceService_ReplaceTags_InvalidTag(t *testing.T) {
	svc := newTestBusinessService()

	created, err := svc.Create(context.Background(), 1, ServiceInput{
		CategoryID:      1,
		Title:           "服务",
		Price:           1,
		DurationMinutes: 1,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	_, err = svc.ReplaceTags(context.Background(), 1, created.ID, ServiceTagsInput{TagIDs: []int64{99}})
	if !errors.Is(err, ErrTagNotFound) {
		t.Errorf("ReplaceTags() error = %v, want %v", err, ErrTagNotFound)
	}
}

func TestServiceService_ReplaceTags_ForbidsOtherProvider(t *testing.T) {
	serviceRepo := newMockServiceRepo()
	providerRepo := newMockServiceProviderRepo()
	_ = providerRepo.Create(context.Background(), &model.ServiceProvider{UserID: 1, BusinessName: "商家A"})
	_ = providerRepo.Create(context.Background(), &model.ServiceProvider{UserID: 2, BusinessName: "商家B"})
	tagRepo := newMockTagRepo()
	_ = tagRepo.Create(context.Background(), &model.Tag{Name: "放松"})
	svc := NewServiceService(serviceRepo, providerRepo, tagRepo, newMockServiceTagRepo())

	created, err := svc.Create(context.Background(), 1, ServiceInput{
		CategoryID:      1,
		Title:           "服务",
		Price:           1,
		DurationMinutes: 1,
	})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	_, err = svc.ReplaceTags(context.Background(), 2, created.ID, ServiceTagsInput{TagIDs: []int64{1}})
	if !errors.Is(err, ErrServiceForbidden) {
		t.Errorf("ReplaceTags() error = %v, want %v", err, ErrServiceForbidden)
	}
}

func TestServiceService_List_NormalizesPagination(t *testing.T) {
	svc := newTestBusinessService()

	result, err := svc.List(context.Background(), model.ServiceFilter{PageSize: 1000})
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if result.Page != 1 {
		t.Errorf("List() page = %d, want 1", result.Page)
	}
	if result.PageSize != 50 {
		t.Errorf("List() pageSize = %d, want 50", result.PageSize)
	}
}

func cloneService(service *model.Service) *model.Service {
	cloned := *service
	return &cloned
}

func serviceToView(service *model.Service) *model.ServiceView {
	return &model.ServiceView{
		ID:              service.ID,
		Title:           service.Title,
		Description:     service.Description,
		Provider:        model.ServiceProviderSummary{ID: service.ProviderID, BusinessName: "舒心养生馆"},
		Category:        model.CategorySummary{ID: service.CategoryID, Name: "美容"},
		Price:           service.Price,
		DurationMinutes: service.DurationMinutes,
		ImageURL:        service.ImageURL,
		Status:          service.Status,
		AvgRating:       service.AvgRating,
		ReviewCount:     service.ReviewCount,
		CreatedAt:       service.CreatedAt,
		UpdatedAt:       service.UpdatedAt,
	}
}
