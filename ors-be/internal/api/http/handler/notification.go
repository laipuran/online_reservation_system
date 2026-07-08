package handler

import (
	"errors"
	"net/http"

	"ors-be/internal/api/http/middleware"
	"ors-be/internal/api/http/response"
	"ors-be/internal/service"
)

type NotificationHandler struct {
	notificationSvc service.NotificationService
}

func NewNotificationHandler(notificationSvc service.NotificationService) *NotificationHandler {
	return &NotificationHandler{notificationSvc: notificationSvc}
}

func (h *NotificationHandler) ListMine() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r.Context())
		if claims == nil {
			response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
			return
		}

		isRead, err := parseOptionalBool(r.URL.Query().Get("is_read"))
		if err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
			return
		}

		page, pageSize, limit, offset, err := parsePagination(r)
		if err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
			return
		}

		notifications, err := h.notificationSvc.ListMine(r.Context(), claims.UserID, isRead, limit, offset)
		if err != nil {
			writeNotificationError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(newListResponse(notifications, page, pageSize)))
	}
}

func (h *NotificationHandler) CountUnread() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r.Context())
		if claims == nil {
			response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
			return
		}

		count, err := h.notificationSvc.CountUnread(r.Context(), claims.UserID)
		if err != nil {
			writeNotificationError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(map[string]int64{"count": count}))
	}
}

func (h *NotificationHandler) MarkRead() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r.Context())
		if claims == nil {
			response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
			return
		}

		id, err := parseIDParam(r, "id", "无效的通知ID")
		if err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
			return
		}

		notification, err := h.notificationSvc.MarkRead(r.Context(), claims.UserID, id)
		if err != nil {
			writeNotificationError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(notification))
	}
}

func (h *NotificationHandler) MarkAllRead() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		claims := middleware.GetClaims(r.Context())
		if claims == nil {
			response.JSON(w, http.StatusUnauthorized, response.Unauthorized("缺少认证信息"))
			return
		}

		count, err := h.notificationSvc.MarkAllRead(r.Context(), claims.UserID)
		if err != nil {
			writeNotificationError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(map[string]int64{"updated_count": count}))
	}
}

func parseOptionalBool(value string) (*bool, error) {
	switch value {
	case "":
		return nil, nil
	case "true":
		parsed := true
		return &parsed, nil
	case "false":
		parsed := false
		return &parsed, nil
	default:
		return nil, errors.New("无效的已读状态")
	}
}

func writeNotificationError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrNotificationInvalidInput),
		errors.Is(err, service.ErrNotificationInvalidType):
		response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
	case errors.Is(err, service.ErrNotificationNotFound):
		response.JSON(w, http.StatusNotFound, response.NotFound(err.Error()))
	default:
		response.JSON(w, http.StatusInternalServerError, response.ServerError("通知操作失败"))
	}
}
