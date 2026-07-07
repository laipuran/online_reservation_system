package service

import (
	"context"

	"ors-be/internal/model"
	"ors-be/internal/repository"
)

type CategoryService interface {
	List(ctx context.Context) ([]*model.Category, error)
}

type categoryService struct {
	categoryRepo repository.CategoryRepository
}

func NewCategoryService(categoryRepo repository.CategoryRepository) CategoryService {
	return &categoryService{categoryRepo: categoryRepo}
}

func (s *categoryService) List(ctx context.Context) ([]*model.Category, error) {
	return s.categoryRepo.List(ctx)
}
