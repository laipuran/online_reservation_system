package service

import (
	"context"
	"errors"

	"ors-be/internal/model"
	"ors-be/internal/repository"
)

var ErrUserInterestInvalidTag = errors.New("标签ID不正确")

type UserInterestsInput struct {
	TagIDs []int64 `json:"tag_ids"`
}

type UserInterestService interface {
	List(ctx context.Context, userID int64) ([]*model.Tag, error)
	Replace(ctx context.Context, userID int64, input UserInterestsInput) ([]*model.Tag, error)
}

type userInterestService struct {
	tagRepo      repository.TagRepository
	interestRepo repository.UserInterestRepository
}

func NewUserInterestService(
	tagRepo repository.TagRepository,
	interestRepo repository.UserInterestRepository,
) UserInterestService {
	return &userInterestService{
		tagRepo:      tagRepo,
		interestRepo: interestRepo,
	}
}

func (s *userInterestService) List(ctx context.Context, userID int64) ([]*model.Tag, error) {
	return s.interestRepo.ListByUserID(ctx, userID)
}

func (s *userInterestService) Replace(ctx context.Context, userID int64, input UserInterestsInput) ([]*model.Tag, error) {
	tagIDs, err := s.normalizeTagIDs(ctx, input.TagIDs)
	if err != nil {
		return nil, err
	}

	if err := s.interestRepo.ReplaceByUserID(ctx, userID, tagIDs); err != nil {
		return nil, err
	}
	return s.interestRepo.ListByUserID(ctx, userID)
}

func (s *userInterestService) normalizeTagIDs(ctx context.Context, tagIDs []int64) ([]int64, error) {
	seen := make(map[int64]struct{}, len(tagIDs))
	normalized := make([]int64, 0, len(tagIDs))

	for _, tagID := range tagIDs {
		if tagID <= 0 {
			return nil, ErrUserInterestInvalidTag
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
