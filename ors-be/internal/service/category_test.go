package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"ors-be/internal/model"
)

type mockCategoryRepo struct {
	categories []*model.Category
	err        error
}

func (m *mockCategoryRepo) Create(ctx context.Context, category *model.Category) error {
	if m.err != nil {
		return m.err
	}
	category.ID = int64(len(m.categories) + 1)
	category.CreatedAt = time.Now()
	m.categories = append(m.categories, category)
	return nil
}

func (m *mockCategoryRepo) GetByID(ctx context.Context, id int64) (*model.Category, error) {
	if m.err != nil {
		return nil, m.err
	}
	for _, category := range m.categories {
		if category.ID == id {
			return category, nil
		}
	}
	return nil, nil
}

func (m *mockCategoryRepo) List(ctx context.Context) ([]*model.Category, error) {
	if m.err != nil {
		return nil, m.err
	}
	return m.categories, nil
}

func TestCategoryService_List_Success(t *testing.T) {
	parentID := int64(1)
	categories := []*model.Category{
		{ID: 1, Name: "医疗", Description: "医疗健康服务", CreatedAt: time.Now()},
		{ID: 2, Name: "口腔护理", ParentID: &parentID, CreatedAt: time.Now()},
	}
	svc := NewCategoryService(&mockCategoryRepo{categories: categories})

	result, err := svc.List(context.Background())
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(result) != 2 {
		t.Fatalf("List() len = %d, want 2", len(result))
	}
	if result[0].Name != "医疗" {
		t.Errorf("List()[0].Name = %s, want 医疗", result[0].Name)
	}
	if result[1].ParentID == nil || *result[1].ParentID != parentID {
		t.Fatalf("List()[1].ParentID = %v, want %d", result[1].ParentID, parentID)
	}
}

func TestCategoryService_List_RepoError(t *testing.T) {
	repoErr := errors.New("db failed")
	svc := NewCategoryService(&mockCategoryRepo{err: repoErr})

	_, err := svc.List(context.Background())
	if !errors.Is(err, repoErr) {
		t.Errorf("List() error = %v, want %v", err, repoErr)
	}
}
