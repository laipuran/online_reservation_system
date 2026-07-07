package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"ors-be/internal/model"
)

type mockTagRepo struct {
	tags   map[int64]*model.Tag
	nextID int64
}

func newMockTagRepo() *mockTagRepo {
	return &mockTagRepo{
		tags:   make(map[int64]*model.Tag),
		nextID: 1,
	}
}

func (m *mockTagRepo) Create(ctx context.Context, tag *model.Tag) error {
	tag.ID = m.nextID
	tag.CreatedAt = time.Now()
	m.nextID++
	m.tags[tag.ID] = cloneTag(tag)
	return nil
}

func (m *mockTagRepo) GetByID(ctx context.Context, id int64) (*model.Tag, error) {
	tag := m.tags[id]
	if tag == nil {
		return nil, nil
	}
	return cloneTag(tag), nil
}

func (m *mockTagRepo) GetByName(ctx context.Context, name string) (*model.Tag, error) {
	for _, tag := range m.tags {
		if tag.Name == name {
			return cloneTag(tag), nil
		}
	}
	return nil, nil
}

func (m *mockTagRepo) List(ctx context.Context) ([]*model.Tag, error) {
	tags := make([]*model.Tag, 0, len(m.tags))
	for _, tag := range m.tags {
		tags = append(tags, cloneTag(tag))
	}
	return tags, nil
}

func TestTagService_Create_Success(t *testing.T) {
	svc := NewTagService(newMockTagRepo())

	tag, err := svc.Create(context.Background(), TagInput{Name: " 放松 "})
	if err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	if tag.ID == 0 {
		t.Error("Create() ID should be non-zero")
	}
	if tag.Name != "放松" {
		t.Errorf("Create() name = %q, want 放松", tag.Name)
	}
}

func TestTagService_Create_InvalidInput(t *testing.T) {
	svc := NewTagService(newMockTagRepo())
	longName := "123456789012345678901234567890123456789012345678901"

	tests := []struct {
		name  string
		input TagInput
		want  error
	}{
		{"empty name", TagInput{Name: " "}, ErrTagNameRequired},
		{"too long", TagInput{Name: longName}, ErrTagNameTooLong},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.Create(context.Background(), tt.input)
			if !errors.Is(err, tt.want) {
				t.Errorf("Create() error = %v, want %v", err, tt.want)
			}
		})
	}
}

func TestTagService_Create_Duplicate(t *testing.T) {
	svc := NewTagService(newMockTagRepo())

	if _, err := svc.Create(context.Background(), TagInput{Name: "放松"}); err != nil {
		t.Fatalf("Create() error = %v", err)
	}
	_, err := svc.Create(context.Background(), TagInput{Name: "放松"})
	if !errors.Is(err, ErrTagAlreadyExists) {
		t.Errorf("Create() error = %v, want %v", err, ErrTagAlreadyExists)
	}
}

func TestTagService_GetByID_NotFound(t *testing.T) {
	svc := NewTagService(newMockTagRepo())

	_, err := svc.GetByID(context.Background(), 1)
	if !errors.Is(err, ErrTagNotFound) {
		t.Errorf("GetByID() error = %v, want %v", err, ErrTagNotFound)
	}
}

func cloneTag(tag *model.Tag) *model.Tag {
	if tag == nil {
		return nil
	}
	copied := *tag
	return &copied
}
