package service

import (
	"context"
	"errors"
	"strings"
	"unicode/utf8"

	"ors-be/internal/model"
	"ors-be/internal/repository"
)

var (
	ErrTagNotFound      = errors.New("标签不存在")
	ErrTagNameRequired  = errors.New("标签名称不能为空")
	ErrTagNameTooLong   = errors.New("标签名称不能超过50个字符")
	ErrTagAlreadyExists = errors.New("标签已存在")
)

type TagInput struct {
	Name string `json:"name"`
}

type TagService interface {
	Create(ctx context.Context, input TagInput) (*model.Tag, error)
	GetByID(ctx context.Context, id int64) (*model.Tag, error)
	List(ctx context.Context) ([]*model.Tag, error)
}

type tagService struct {
	tagRepo repository.TagRepository
}

func NewTagService(tagRepo repository.TagRepository) TagService {
	return &tagService{tagRepo: tagRepo}
}

func (s *tagService) Create(ctx context.Context, input TagInput) (*model.Tag, error) {
	name := strings.TrimSpace(input.Name)
	if name == "" {
		return nil, ErrTagNameRequired
	}
	if utf8.RuneCountInString(name) > 50 {
		return nil, ErrTagNameTooLong
	}

	existing, err := s.tagRepo.GetByName(ctx, name)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, ErrTagAlreadyExists
	}

	tag := &model.Tag{Name: name}
	if err := s.tagRepo.Create(ctx, tag); err != nil {
		return nil, err
	}
	return tag, nil
}

func (s *tagService) GetByID(ctx context.Context, id int64) (*model.Tag, error) {
	tag, err := s.tagRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if tag == nil {
		return nil, ErrTagNotFound
	}
	return tag, nil
}

func (s *tagService) List(ctx context.Context) ([]*model.Tag, error) {
	return s.tagRepo.List(ctx)
}
