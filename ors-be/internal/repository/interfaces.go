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

type ServiceTagRepository interface {
	ReplaceByServiceID(ctx context.Context, serviceID int64, tagIDs []int64) error
	ListByServiceID(ctx context.Context, serviceID int64) ([]*model.Tag, error)
}

type UserInterestRepository interface {
	ReplaceByUserID(ctx context.Context, userID int64, tagIDs []int64) error
	ListByUserID(ctx context.Context, userID int64) ([]*model.Tag, error)
}

type ServiceRepository interface {
	Create(ctx context.Context, service *model.Service) error
	GetByID(ctx context.Context, id int64) (*model.Service, error)
	GetViewByID(ctx context.Context, id int64) (*model.ServiceView, error)
	List(ctx context.Context, filter model.ServiceFilter) ([]*model.ServiceView, int, error)
	Update(ctx context.Context, service *model.Service) error
	UpdateStatus(ctx context.Context, id int64, status string) error
}

type ReservationRepository interface {
	Create(ctx context.Context, reservation *model.Reservation) error
	GetByID(ctx context.Context, id int64) (*model.Reservation, error)
	GetByIDForUser(ctx context.Context, id, userID int64) (*model.Reservation, error)
	GetByIDForProvider(ctx context.Context, id, providerID int64) (*model.Reservation, error)
	ListByUserID(ctx context.Context, userID int64, status string, limit, offset int) ([]*model.Reservation, error)
	ListByProviderID(ctx context.Context, providerID int64, status string, limit, offset int) ([]*model.Reservation, error)
	UpdateStatus(ctx context.Context, id int64, status string) (*model.Reservation, error)
}

type ReviewRepository interface {
	Create(ctx context.Context, review *model.Review) error
	GetByID(ctx context.Context, id int64) (*model.Review, error)
	GetByReservationID(ctx context.Context, reservationID int64) (*model.Review, error)
	ListByServiceID(ctx context.Context, serviceID int64, limit, offset int) ([]*model.Review, error)
	ListByUserID(ctx context.Context, userID int64, limit, offset int) ([]*model.Review, error)
	ListByProviderID(ctx context.Context, providerID int64, limit, offset int) ([]*model.Review, error)
}

type NotificationRepository interface {
	Create(ctx context.Context, notification *model.Notification) error
	GetByID(ctx context.Context, id int64) (*model.Notification, error)
	ListByUserID(ctx context.Context, userID int64, isRead *bool, limit, offset int) ([]*model.Notification, error)
	CountUnread(ctx context.Context, userID int64) (int64, error)
	MarkRead(ctx context.Context, id, userID int64) (*model.Notification, error)
	MarkAllRead(ctx context.Context, userID int64) (int64, error)
}
