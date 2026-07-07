package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"ors-be/internal/api/http/middleware"
	"ors-be/internal/api/http/response"
	"ors-be/internal/model"
	"ors-be/internal/service"
)

type ServiceHandler struct {
	serviceSvc service.ServiceService
}

func NewServiceHandler(serviceSvc service.ServiceService) *ServiceHandler {
	return &ServiceHandler{serviceSvc: serviceSvc}
}

func (h *ServiceHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		filter, err := parseServiceFilter(r)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
			return
		}
		filter.Status = "active"

		result, err := h.serviceSvc.List(r.Context(), filter)
		if err != nil {
			writeServiceError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(result))
	}
}

func (h *ServiceHandler) ListByProvider() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providerID, err := parseIDParam(r, "id", "无效的服务提供者ID")
		if err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
			return
		}

		filter, err := parseServiceFilter(r)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
			return
		}
		filter.ProviderID = &providerID
		filter.Status = "active"

		result, err := h.serviceSvc.List(r.Context(), filter)
		if err != nil {
			writeServiceError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(result))
	}
}

func (h *ServiceHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseIDParam(r, "id", "无效的服务ID")
		if err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
			return
		}

		serviceView, err := h.serviceSvc.GetByID(r.Context(), id)
		if err != nil {
			writeServiceError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(serviceView))
	}
}

func (h *ServiceHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r.Context())
		if claims == nil {
			response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
			return
		}

		var req service.ServiceInput
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail("无效的请求体"))
			return
		}

		serviceView, err := h.serviceSvc.Create(r.Context(), claims.UserID, req)
		if err != nil {
			writeServiceError(w, err)
			return
		}

		response.JSON(w, http.StatusCreated, response.Created(serviceView))
	}
}

func (h *ServiceHandler) Update() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r.Context())
		if claims == nil {
			response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
			return
		}

		id, err := parseIDParam(r, "id", "无效的服务ID")
		if err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
			return
		}

		var req service.ServiceInput
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail("无效的请求体"))
			return
		}

		serviceView, err := h.serviceSvc.Update(r.Context(), claims.UserID, id, req)
		if err != nil {
			writeServiceError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(serviceView))
	}
}

func (h *ServiceHandler) UpdateStatus() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r.Context())
		if claims == nil {
			response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
			return
		}

		id, err := parseIDParam(r, "id", "无效的服务ID")
		if err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
			return
		}

		var req service.ServiceStatusInput
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail("无效的请求体"))
			return
		}

		serviceView, err := h.serviceSvc.UpdateStatus(r.Context(), claims.UserID, id, req.Status)
		if err != nil {
			writeServiceError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(serviceView))
	}
}

func parseServiceFilter(r *http.Request) (model.ServiceFilter, error) {
	query := r.URL.Query()
	filter := model.ServiceFilter{
		Keyword:   query.Get("keyword"),
		SortBy:    query.Get("sort_by"),
		SortOrder: query.Get("sort_order"),
	}

	var err error
	if filter.CategoryID, err = parseOptionalInt64(query.Get("category_id")); err != nil {
		return filter, errors.New("无效的分类ID")
	}
	if filter.ProviderID, err = parseOptionalInt64(query.Get("provider_id")); err != nil {
		return filter, errors.New("无效的服务提供者ID")
	}
	if filter.MinPrice, err = parseOptionalFloat64(query.Get("min_price")); err != nil {
		return filter, errors.New("无效的最低价格")
	}
	if filter.MaxPrice, err = parseOptionalFloat64(query.Get("max_price")); err != nil {
		return filter, errors.New("无效的最高价格")
	}
	if filter.Page, err = parseOptionalInt(query.Get("page")); err != nil {
		return filter, errors.New("无效的页码")
	}
	if filter.PageSize, err = parseOptionalInt(query.Get("page_size")); err != nil {
		return filter, errors.New("无效的分页大小")
	}
	return filter, nil
}

func parseIDParam(r *http.Request, name, message string) (int64, error) {
	id, err := strconv.ParseInt(chi.URLParam(r, name), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New(message)
	}
	return id, nil
}

func parseOptionalInt64(value string) (*int64, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := strconv.ParseInt(value, 10, 64)
	if err != nil || parsed <= 0 {
		return nil, errors.New("invalid int64")
	}
	return &parsed, nil
}

func parseOptionalFloat64(value string) (*float64, error) {
	if value == "" {
		return nil, nil
	}
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil || parsed < 0 {
		return nil, errors.New("invalid float64")
	}
	return &parsed, nil
}

func parseOptionalInt(value string) (int, error) {
	if value == "" {
		return 0, nil
	}
	parsed, err := strconv.Atoi(value)
	if err != nil || parsed <= 0 {
		return 0, errors.New("invalid int")
	}
	return parsed, nil
}

func writeServiceError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrServiceTitleRequired),
		errors.Is(err, service.ErrServiceInvalidCategory),
		errors.Is(err, service.ErrServiceInvalidPrice),
		errors.Is(err, service.ErrServiceInvalidDuration),
		errors.Is(err, service.ErrServiceInvalidStatus),
		errors.Is(err, service.ErrBusinessNameRequired):
		response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
	case errors.Is(err, service.ErrProviderNotFound):
		response.JSON(w, http.StatusNotFound, response.NotFound(err.Error()))
	case errors.Is(err, service.ErrServiceNotFound):
		response.JSON(w, http.StatusNotFound, response.NotFound(err.Error()))
	case errors.Is(err, service.ErrServiceForbidden):
		response.JSON(w, http.StatusForbidden, response.Forbidden(err.Error()))
	default:
		response.JSON(w, http.StatusInternalServerError, response.ServerError("服务操作失败"))
	}
}
