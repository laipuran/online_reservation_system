package service

import (
	"context"
	"errors"
	"testing"

	"ors-be/internal/model"
)

type mockUserInterestRepo struct {
	tagsByUser map[int64][]*model.Tag
}

func newMockUserInterestRepo() *mockUserInterestRepo {
	return &mockUserInterestRepo{tagsByUser: make(map[int64][]*model.Tag)}
}

func (m *mockUserInterestRepo) ReplaceByUserID(ctx context.Context, userID int64, tagIDs []int64) error {
	tags := make([]*model.Tag, 0, len(tagIDs))
	for _, tagID := range tagIDs {
		tags = append(tags, &model.Tag{ID: tagID, Name: "标签"})
	}
	m.tagsByUser[userID] = tags
	return nil
}

func (m *mockUserInterestRepo) ListByUserID(ctx context.Context, userID int64) ([]*model.Tag, error) {
	tags := make([]*model.Tag, 0, len(m.tagsByUser[userID]))
	for _, tag := range m.tagsByUser[userID] {
		tags = append(tags, cloneTag(tag))
	}
	return tags, nil
}

func newTestUserInterestService() UserInterestService {
	tagRepo := newMockTagRepo()
	_ = tagRepo.Create(context.Background(), &model.Tag{Name: "放松"})
	_ = tagRepo.Create(context.Background(), &model.Tag{Name: "塑形"})
	return NewUserInterestService(tagRepo, newMockUserInterestRepo())
}

func TestUserInterestService_Replace_Success(t *testing.T) {
	svc := newTestUserInterestService()

	tags, err := svc.Replace(context.Background(), 1, UserInterestsInput{TagIDs: []int64{1, 2, 1}})
	if err != nil {
		t.Fatalf("Replace() error = %v", err)
	}
	if len(tags) != 2 {
		t.Fatalf("Replace() len = %d, want 2", len(tags))
	}

	listed, err := svc.List(context.Background(), 1)
	if err != nil {
		t.Fatalf("List() error = %v", err)
	}
	if len(listed) != 2 {
		t.Errorf("List() len = %d, want 2", len(listed))
	}
}

func TestUserInterestService_Replace_EmptyClearsTags(t *testing.T) {
	svc := newTestUserInterestService()

	if _, err := svc.Replace(context.Background(), 1, UserInterestsInput{TagIDs: []int64{1, 2}}); err != nil {
		t.Fatalf("Replace() error = %v", err)
	}
	tags, err := svc.Replace(context.Background(), 1, UserInterestsInput{TagIDs: []int64{}})
	if err != nil {
		t.Fatalf("Replace(empty) error = %v", err)
	}
	if len(tags) != 0 {
		t.Errorf("Replace(empty) len = %d, want 0", len(tags))
	}
}

func TestUserInterestService_Replace_InvalidTagID(t *testing.T) {
	svc := newTestUserInterestService()

	_, err := svc.Replace(context.Background(), 1, UserInterestsInput{TagIDs: []int64{0}})
	if !errors.Is(err, ErrUserInterestInvalidTag) {
		t.Errorf("Replace() error = %v, want %v", err, ErrUserInterestInvalidTag)
	}
}

func TestUserInterestService_Replace_TagNotFound(t *testing.T) {
	svc := newTestUserInterestService()

	_, err := svc.Replace(context.Background(), 1, UserInterestsInput{TagIDs: []int64{99}})
	if !errors.Is(err, ErrTagNotFound) {
		t.Errorf("Replace() error = %v, want %v", err, ErrTagNotFound)
	}
}
