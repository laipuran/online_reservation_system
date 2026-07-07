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

type TagRepository interface {
	Create(ctx context.Context, tag *model.Tag) error
	GetByID(ctx context.Context, id int64) (*model.Tag, error)
	GetByName(ctx context.Context, name string) (*model.Tag, error)
	List(ctx context.Context) ([]*model.Tag, error)
}

type ServiceRepository interface {
	Create(ctx context.Context, service *model.Service) error
	GetByID(ctx context.Context, id int64) (*model.Service, error)
	GetViewByID(ctx context.Context, id int64) (*model.ServiceView, error)
	List(ctx context.Context, filter model.ServiceFilter) ([]*model.ServiceView, int, error)
	Update(ctx context.Context, service *model.Service) error
	UpdateStatus(ctx context.Context, id int64, status string) error
}
