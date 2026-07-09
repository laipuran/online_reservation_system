package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"ors-be/internal/api/http/middleware"
	"ors-be/internal/api/http/response"
	"ors-be/internal/service"
)

type ServiceProviderHandler struct {
	providerSvc service.ServiceProviderService
}

func NewServiceProviderHandler(providerSvc service.ServiceProviderService) *ServiceProviderHandler {
	return &ServiceProviderHandler{providerSvc: providerSvc}
}

func (h *ServiceProviderHandler) CreateMine() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r.Context())
		if claims == nil {
			response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
			return
		}

		var req service.ServiceProviderInput
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail("无效的请求体"))
			return
		}

		provider, err := h.providerSvc.Create(r.Context(), claims.UserID, req)
		if err != nil {
			writeProviderError(w, err)
			return
		}

		response.JSON(w, http.StatusCreated, response.Created(provider))
	}
}

func (h *ServiceProviderHandler) GetMine() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r.Context())
		if claims == nil {
			response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
			return
		}

		provider, err := h.providerSvc.GetMine(r.Context(), claims.UserID)
		if err != nil {
			writeProviderError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(provider))
	}
}

func (h *ServiceProviderHandler) UpdateMine() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r.Context())
		if claims == nil {
			response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
			return
		}

		var req service.ServiceProviderInput
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail("无效的请求体"))
			return
		}

		provider, err := h.providerSvc.UpdateMine(r.Context(), claims.UserID, req)
		if err != nil {
			writeProviderError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(provider))
	}
}

func (h *ServiceProviderHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
		if err != nil || id <= 0 {
			response.JSON(w, http.StatusBadRequest, response.Fail("无效的服务提供者ID"))
			return
		}

		provider, err := h.providerSvc.GetByID(r.Context(), id)
		if err != nil {
			writeProviderError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(provider))
	}
}

func writeProviderError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrBusinessNameRequired),
		errors.Is(err, service.ErrInvalidEmail),
		errors.Is(err, service.ErrInvalidPhone):
		response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
	case errors.Is(err, service.ErrProviderAlreadyExists):
		response.JSON(w, http.StatusConflict, response.Conflict(err.Error()))
	case errors.Is(err, service.ErrProviderNotFound):
		response.JSON(w, http.StatusNotFound, response.NotFound(err.Error()))
	case errors.Is(err, service.ErrProviderForbidden):
		response.JSON(w, http.StatusForbidden, response.Forbidden(err.Error()))
	default:
		response.JSON(w, http.StatusInternalServerError, response.ServerError("服务提供者操作失败"))
	}
}
