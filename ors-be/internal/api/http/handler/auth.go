package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"ors-be/internal/api/http/response"
	"ors-be/internal/service"
)

type AuthHandler struct {
	authSvc service.AuthService
}

func NewAuthHandler(authSvc service.AuthService) *AuthHandler {
	return &AuthHandler{authSvc: authSvc}
}

type registerRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

func (h *AuthHandler) Register() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req registerRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail("无效的请求体"))
			return
		}

		result, err := h.authSvc.Register(r.Context(), req.Email, req.Password, req.Name, req.Role)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrEmailAlreadyRegistered):
				response.JSON(w, http.StatusConflict, response.Conflict(err.Error()))
			case errors.Is(err, service.ErrInvalidEmail),
				errors.Is(err, service.ErrWeakPassword),
				errors.Is(err, service.ErrNameRequired),
				errors.Is(err, service.ErrInvalidRole):
				response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
			default:
				response.JSON(w, http.StatusInternalServerError, response.ServerError("注册失败"))
			}
			return
		}

		response.JSON(w, http.StatusCreated, response.Created(result))
	}
}

type loginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (h *AuthHandler) Login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req loginRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail("无效的请求体"))
			return
		}

		result, err := h.authSvc.Login(r.Context(), req.Email, req.Password)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrInvalidCredentials):
				response.JSON(w, http.StatusUnauthorized, response.Unauthorized(err.Error()))
			default:
				response.JSON(w, http.StatusInternalServerError, response.ServerError("登录失败"))
			}
			return
		}

		response.JSON(w, http.StatusOK, response.OK(result))
	}
}
