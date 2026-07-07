package repository

import (
	"context"

	"ors-be/internal/model"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByEmail(ctx context.Context, email string) (*model.User, error)
	GetByID(ctx context.Context, id int64) (*model.User, error)
}

type ServiceProviderRepository interface {
	Create(ctx context.Context, provider *model.ServiceProvider) error
	GetByID(ctx context.Context, id int64) (*model.ServiceProvider, error)
	GetByUserID(ctx context.Context, userID int64) (*model.ServiceProvider, error)
	Update(ctx context.Context, provider *model.ServiceProvider) error
}
  
type CategoryRepository interface {
	Create(ctx context.Context, category *model.Category) error
	GetByID(ctx context.Context, id int64) (*model.Category, error)
	List(ctx context.Context) ([]*model.Category, error)
}
