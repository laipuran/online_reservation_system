package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"ors-be/internal/api/http/middleware"
	"ors-be/internal/api/http/response"
	"ors-be/internal/service"
)

type ReviewHandler struct {
	reviewSvc service.ReviewService
}

func NewReviewHandler(reviewSvc service.ReviewService) *ReviewHandler {
	return &ReviewHandler{reviewSvc: reviewSvc}
}

func (h *ReviewHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r.Context())
		if claims == nil {
			response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
			return
		}

		var req service.ReviewInput
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail("无效的请求体"))
			return
		}

		review, err := h.reviewSvc.Create(r.Context(), claims.UserID, req)
		if err != nil {
			writeReviewError(w, err)
			return
		}

		response.JSON(w, http.StatusCreated, response.Created(review))
	}
}

func (h *ReviewHandler) ListByService() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		serviceID, err := parseIDParam(r, "id", "无效的服务ID")
		if err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
			return
		}

		page, pageSize, limit, offset, err := parsePagination(r)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
			return
		}

		reviews, err := h.reviewSvc.ListByService(r.Context(), serviceID, limit, offset)
		if err != nil {
			writeReviewError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(newListResponse(reviews, page, pageSize)))
	}
}

func (h *ReviewHandler) ListMine() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r.Context())
		if claims == nil {
			response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
			return
		}

		page, pageSize, limit, offset, err := parsePagination(r)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
			return
		}

		reviews, err := h.reviewSvc.ListMine(r.Context(), claims.UserID, limit, offset)
		if err != nil {
			writeReviewError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(newListResponse(reviews, page, pageSize)))
	}
}

func (h *ReviewHandler) ListByProvider() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		providerID, err := parseIDParam(r, "id", "无效的服务提供者ID")
		if err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
			return
		}

		page, pageSize, limit, offset, err := parsePagination(r)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
			return
		}

		reviews, err := h.reviewSvc.ListByProvider(r.Context(), providerID, limit, offset)
		if err != nil {
			writeReviewError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(newListResponse(reviews, page, pageSize)))
	}
}

func writeReviewError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrReviewInvalidScope),
		errors.Is(err, service.ErrReviewInvalidInput),
		errors.Is(err, service.ErrReviewInvalidRating),
		errors.Is(err, service.ErrReviewReservationNotDone):
		response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
	case errors.Is(err, service.ErrReviewReservationNotFound):
		response.JSON(w, http.StatusNotFound, response.NotFound(err.Error()))
	case errors.Is(err, service.ErrReviewForbidden):
		response.JSON(w, http.StatusForbidden, response.Forbidden(err.Error()))
	case errors.Is(err, service.ErrReviewAlreadyExists):
		response.JSON(w, http.StatusConflict, response.Conflict(err.Error()))
	default:
		response.JSON(w, http.StatusInternalServerError, response.ServerError("评价操作失败"))
	}
}

func parsePagination(r *http.Request) (page int, pageSize int, limit int, offset int, err error) {
	query := r.URL.Query()

	page, err = parseOptionalInt(query.Get("page"))
	if err != nil {
		return 0, 0, 0, 0, errors.New("无效的页码")
	}
	if page == 0 {
		page = 1
	}

	pageSize, err = parseOptionalInt(query.Get("page_size"))
	if err != nil {
		return 0, 0, 0, 0, errors.New("无效的分页大小")
	}
	if pageSize == 0 {
		pageSize = 20
	}
	if pageSize > 50 {
		pageSize = 50
	}

	limit = pageSize
	offset = (page - 1) * pageSize
	return page, pageSize, limit, offset, nil
}

func newListResponse(items interface{}, page, pageSize int) map[string]interface{} {
	return map[string]interface{}{
		"items":     items,
		"page":      page,
		"page_size": pageSize,
	}
}
