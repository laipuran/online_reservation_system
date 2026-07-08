package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"ors-be/internal/api/http/middleware"
	"ors-be/internal/api/http/response"
	"ors-be/internal/repository"
	"ors-be/internal/service"
)

type ReservationHandler struct {
	reservationSvc service.ReservationService
}

func NewReservationHandler(reservationSvc service.ReservationService) *ReservationHandler {
	return &ReservationHandler{reservationSvc: reservationSvc}
}

func (h *ReservationHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r.Context())
		if claims == nil {
			response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
			return
		}

		var req service.ReservationInput
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail("无效的请求体"))
			return
		}

		reservation, err := h.reservationSvc.Create(r.Context(), claims.UserID, req)
		if err != nil {
			writeReservationError(w, err)
			return
		}

		response.JSON(w, http.StatusCreated, response.Created(reservation))
	}
}

func (h *ReservationHandler) ListMine() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r.Context())
		if claims == nil {
			response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
			return
		}

		page, pageSize, err := parseReservationPagination(r)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
			return
		}

		result, err := h.reservationSvc.ListMine(r.Context(), claims.UserID, r.URL.Query().Get("status"), page, pageSize)
		if err != nil {
			writeReservationError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(result))
	}
}

func (h *ReservationHandler) GetMine() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r.Context())
		if claims == nil {
			response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
			return
		}

		id, err := parseIDParam(r, "id", "无效的预约 ID")
		if err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
			return
		}

		reservation, err := h.reservationSvc.GetMine(r.Context(), claims.UserID, id)
		if err != nil {
			writeReservationError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(reservation))
	}
}

func (h *ReservationHandler) CancelMine() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r.Context())
		if claims == nil {
			response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
			return
		}

		id, err := parseIDParam(r, "id", "无效的预约 ID")
		if err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
			return
		}

		reservation, err := h.reservationSvc.CancelMine(r.Context(), claims.UserID, id)
		if err != nil {
			writeReservationError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(reservation))
	}
}

func (h *ReservationHandler) ListForProvider() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r.Context())
		if claims == nil {
			response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
			return
		}

		page, pageSize, err := parseReservationPagination(r)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
			return
		}

		result, err := h.reservationSvc.ListForProvider(r.Context(), claims.UserID, r.URL.Query().Get("status"), page, pageSize)
		if err != nil {
			writeReservationError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(result))
	}
}

func (h *ReservationHandler) ConfirmForProvider() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r.Context())
		if claims == nil {
			response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
			return
		}

		id, err := parseIDParam(r, "id", "无效的预约 ID")
		if err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
			return
		}

		reservation, err := h.reservationSvc.ConfirmForProvider(r.Context(), claims.UserID, id)
		if err != nil {
			writeReservationError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(reservation))
	}
}

func (h *ReservationHandler) RejectForProvider() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r.Context())
		if claims == nil {
			response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
			return
		}

		id, err := parseIDParam(r, "id", "无效的预约 ID")
		if err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
			return
		}

		reservation, err := h.reservationSvc.RejectForProvider(r.Context(), claims.UserID, id)
		if err != nil {
			writeReservationError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(reservation))
	}
}

func parseReservationPagination(r *http.Request) (int, int, error) {
	query := r.URL.Query()

	page, err := parseOptionalInt(query.Get("page"))
	if err != nil {
		return 0, 0, errors.New("无效的页码")
	}

	pageSize, err := parseOptionalInt(query.Get("page_size"))
	if err != nil {
		return 0, 0, errors.New("无效的分页大小")
	}

	return page, pageSize, nil
}

func writeReservationError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrReservationInvalidInput),
		errors.Is(err, service.ErrReservationInvalidStatus),
		errors.Is(err, service.ErrReservationCannotCancel),
		errors.Is(err, service.ErrReservationCannotConfirm),
		errors.Is(err, service.ErrReservationCannotReject):
		response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
	case errors.Is(err, repository.ErrReservationTimeConflict):
		response.JSON(w, http.StatusConflict, response.Conflict(err.Error()))
	case errors.Is(err, service.ErrServiceNotFound),
		errors.Is(err, service.ErrProviderNotFound),
		errors.Is(err, service.ErrReservationNotFound):
		response.JSON(w, http.StatusNotFound, response.NotFound(err.Error()))
	default:
		response.JSON(w, http.StatusInternalServerError, response.ServerError("预约操作失败"))
	}
}
