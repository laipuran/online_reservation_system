package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"ors-be/internal/api/http/middleware"
	"ors-be/internal/api/http/response"
	"ors-be/internal/service"
)

type UserInterestHandler struct {
	interestSvc service.UserInterestService
}

func NewUserInterestHandler(interestSvc service.UserInterestService) *UserInterestHandler {
	return &UserInterestHandler{interestSvc: interestSvc}
}

func (h *UserInterestHandler) ListMine() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r.Context())
		if claims == nil {
			response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
			return
		}

		tags, err := h.interestSvc.List(r.Context(), claims.UserID)
		if err != nil {
			writeUserInterestError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(tags))
	}
}

func (h *UserInterestHandler) ReplaceMine() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r.Context())
		if claims == nil {
			response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
			return
		}

		var req service.UserInterestsInput
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail("无效的请求体"))
			return
		}

		tags, err := h.interestSvc.Replace(r.Context(), claims.UserID, req)
		if err != nil {
			writeUserInterestError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(tags))
	}
}

func writeUserInterestError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrUserInterestInvalidTag):
		response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
	case errors.Is(err, service.ErrTagNotFound):
		response.JSON(w, http.StatusNotFound, response.NotFound(err.Error()))
	default:
		response.JSON(w, http.StatusInternalServerError, response.ServerError("用户兴趣标签操作失败"))
	}
}
