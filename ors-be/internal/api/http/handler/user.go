package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"ors-be/internal/api/http/middleware"
	"ors-be/internal/api/http/response"
	"ors-be/internal/service"
)

type UserHandler struct {
	userSvc service.UserService
}

func NewUserHandler(userSvc service.UserService) *UserHandler {
	return &UserHandler{userSvc: userSvc}
}

func (h *UserHandler) GetMine() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r.Context())
		if claims == nil {
			response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
			return
		}

		user, err := h.userSvc.GetMine(r.Context(), claims.UserID)
		if err != nil {
			writeUserError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(user))
	}
}

func (h *UserHandler) UpdateMine() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r.Context())
		if claims == nil {
			response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
			return
		}

		var req service.UserInput
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail("无效的请求体"))
			return
		}

		user, err := h.userSvc.UpdateMine(r.Context(), claims.UserID, req)
		if err != nil {
			writeUserError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(user))
	}
}

func (h *UserHandler) UpdatePassword() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r.Context())
		if claims == nil {
			response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
			return
		}

		var req service.UserPasswordInput
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail("无效的请求体"))
			return
		}

		if err := h.userSvc.UpdatePassword(r.Context(), claims.UserID, req); err != nil {
			writeUserError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(map[string]string{"message": "密码修改成功"}))
	}
}

func writeUserError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrNameRequired),
		errors.Is(err, service.ErrWeakPassword),
		errors.Is(err, service.ErrInvalidPhone):
		response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
	case errors.Is(err, service.ErrCurrentPasswordWrong):
		response.JSON(w, http.StatusUnauthorized, response.Unauthorized(err.Error()))
	case errors.Is(err, service.ErrUserNotFound):
		response.JSON(w, http.StatusNotFound, response.NotFound(err.Error()))
	default:
		response.JSON(w, http.StatusInternalServerError, response.ServerError("用户操作失败"))
	}
}
