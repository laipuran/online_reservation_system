package service

import (
	"context"
	"errors"
	"strings"
	"time"

	"ors-be/internal/model"
	"ors-be/internal/repository"
)

const (
	ReservationStatusPending   = "pending"
	ReservationStatusConfirmed = "confirmed"
	ReservationStatusCompleted = "completed"
	ReservationStatusCancelled = "cancelled"
	ReservationStatusRejected  = "rejected"
)

var (
	ErrReservationNotFound      = errors.New("预约不存在")
	ErrReservationInvalidStatus = errors.New("预约状态无效")
	ErrReservationInvalidInput  = errors.New("预约参数无效")
	ErrReservationCannotCancel  = errors.New("当前预约状态不可取消")
	ErrReservationCannotConfirm = errors.New("当前预约状态不可确认")
	ErrReservationCannotReject  = errors.New("当前预约状态不可拒绝")
)

type ReservationInput struct {
	ServiceID int64     `json:"service_id"`
	StartTime time.Time `json:"start_time"`
	Note      string    `json:"note"`
}

type ReservationListResult struct {
	Items    []*model.Reservation `json:"items"`
	Page     int                  `json:"page"`
	PageSize int                  `json:"page_size"`
}

type ReservationService interface {
	Create(ctx context.Context, userID int64, input ReservationInput) (*model.ReservationView, error)
	GetMine(ctx context.Context, userID, id int64) (*model.Reservation, error)
	ListMine(ctx context.Context, userID int64, status string, page, pageSize int) (*ReservationListResult, error)
	CancelMine(ctx context.Context, userID, id int64) (*model.Reservation, error)
	ListForProvider(ctx context.Context, userID int64, status string, page, pageSize int) (*ReservationListResult, error)
	ConfirmForProvider(ctx context.Context, userID, id int64) (*model.Reservation, error)
	RejectForProvider(ctx context.Context, userID, id int64) (*model.Reservation, error)
}

type reservationService struct {
	reservationRepo repository.ReservationRepository
	serviceRepo     repository.ServiceRepository
	providerRepo    repository.ServiceProviderRepository
}

func NewReservationService(
	reservationRepo repository.ReservationRepository,
	serviceRepo repository.ServiceRepository,
	providerRepo repository.ServiceProviderRepository,
) ReservationService {
	return &reservationService{
		reservationRepo: reservationRepo,
		serviceRepo:     serviceRepo,
		providerRepo:    providerRepo,
	}
}

func (s *reservationService) Create(ctx context.Context, userID int64, input ReservationInput) (*model.ReservationView, error) {
	if input.ServiceID <= 0 || input.StartTime.IsZero() {
		return nil, ErrReservationInvalidInput
	}

	serviceItem, err := s.serviceRepo.GetByID(ctx, input.ServiceID)
	if err != nil {
		return nil, err
	}
	if serviceItem == nil || serviceItem.Status != "active" {
		return nil, ErrServiceNotFound
	}

	reservation := &model.Reservation{
		UserID:    userID,
		ServiceID: input.ServiceID,
		StartTime: input.StartTime,
		EndTime:   input.StartTime.Add(time.Duration(serviceItem.DurationMinutes) * time.Minute),
		Status:    ReservationStatusPending,
		Note:      strings.TrimSpace(input.Note),
	}

	if err := s.reservationRepo.Create(ctx, reservation); err != nil {
		return nil, err
	}

	serviceView, err := s.serviceRepo.GetViewByID(ctx, input.ServiceID)
	if err != nil {
		return nil, err
	}
	if serviceView == nil {
		return nil, ErrServiceNotFound
	}

	return reservationToView(reservation, serviceView), nil
}

func (s *reservationService) GetMine(ctx context.Context, userID, id int64) (*model.Reservation, error) {
	reservation, err := s.reservationRepo.GetByIDForUser(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if reservation == nil {
		return nil, ErrReservationNotFound
	}
	return reservation, nil
}

func (s *reservationService) ListMine(ctx context.Context, userID int64, status string, page, pageSize int) (*ReservationListResult, error) {
	status = normalizeReservationStatusFilter(status)
	if !isReservationStatusFilter(status) {
		return nil, ErrReservationInvalidStatus
	}
	page, pageSize = normalizeReservationPagination(page, pageSize)
	items, err := s.reservationRepo.ListByUserID(ctx, userID, status, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	return &ReservationListResult{Items: items, Page: page, PageSize: pageSize}, nil
}

func (s *reservationService) ListForProvider(ctx context.Context, userID int64, status string, page, pageSize int) (*ReservationListResult, error) {
	provider, err := s.providerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if provider == nil {
		return nil, ErrProviderNotFound
	}

	status = normalizeReservationStatusFilter(status)
	if !isReservationStatusFilter(status) {
		return nil, ErrReservationInvalidStatus
	}
	page, pageSize = normalizeReservationPagination(page, pageSize)
	items, err := s.reservationRepo.ListByProviderID(ctx, provider.ID, status, pageSize, (page-1)*pageSize)
	if err != nil {
		return nil, err
	}
	return &ReservationListResult{Items: items, Page: page, PageSize: pageSize}, nil
}

func (s *reservationService) CancelMine(ctx context.Context, userID, id int64) (*model.Reservation, error) {
	reservation, err := s.GetMine(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	if reservation.Status != ReservationStatusPending && reservation.Status != ReservationStatusConfirmed {
		return nil, ErrReservationCannotCancel
	}
	return s.updateStatus(ctx, id, ReservationStatusCancelled)
}

func (s *reservationService) ConfirmForProvider(ctx context.Context, userID, id int64) (*model.Reservation, error) {
	reservation, err := s.getForProviderUser(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	if reservation.Status != ReservationStatusPending {
		return nil, ErrReservationCannotConfirm
	}
	return s.updateStatus(ctx, id, ReservationStatusConfirmed)
}

func (s *reservationService) RejectForProvider(ctx context.Context, userID, id int64) (*model.Reservation, error) {
	reservation, err := s.getForProviderUser(ctx, userID, id)
	if err != nil {
		return nil, err
	}
	if reservation.Status != ReservationStatusPending {
		return nil, ErrReservationCannotReject
	}
	return s.updateStatus(ctx, id, ReservationStatusRejected)
}

func (s *reservationService) getForProviderUser(ctx context.Context, userID, id int64) (*model.Reservation, error) {
	provider, err := s.providerRepo.GetByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if provider == nil {
		return nil, ErrProviderNotFound
	}

	reservation, err := s.reservationRepo.GetByIDForProvider(ctx, id, provider.ID)
	if err != nil {
		return nil, err
	}
	if reservation == nil {
		return nil, ErrReservationNotFound
	}
	return reservation, nil
}

func (s *reservationService) updateStatus(ctx context.Context, id int64, status string) (*model.Reservation, error) {
	reservation, err := s.reservationRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		return nil, err
	}
	if reservation == nil {
		return nil, ErrReservationNotFound
	}
	return reservation, nil
}

func isReservationStatusFilter(status string) bool {
	switch status {
	case "", ReservationStatusPending, ReservationStatusConfirmed, ReservationStatusCompleted, ReservationStatusCancelled, ReservationStatusRejected:
		return true
	default:
		return false
	}
}

func normalizeReservationStatusFilter(status string) string {
	status = strings.TrimSpace(strings.ToLower(status))
	if status == "all" {
		return ""
	}
	return status
}

func normalizeReservationPagination(page, pageSize int) (int, int) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}
	if pageSize > 100 {
		pageSize = 100
	}
	return page, pageSize
}

func reservationToView(reservation *model.Reservation, serviceView *model.ServiceView) *model.ReservationView {
	return &model.ReservationView{
		ID: reservation.ID,
		Service: model.ReservationServiceSummary{
			ID:    serviceView.ID,
			Title: serviceView.Title,
			Provider: model.ReservationServiceProviderSummary{
				ID:           serviceView.Provider.ID,
				BusinessName: serviceView.Provider.BusinessName,
			},
		},
		StartTime: reservation.StartTime,
		EndTime:   reservation.EndTime,
		Status:    reservation.Status,
		Note:      reservation.Note,
		CreatedAt: reservation.CreatedAt,
	}
}
