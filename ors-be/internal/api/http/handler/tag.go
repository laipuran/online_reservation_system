package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"ors-be/internal/api/http/response"
	"ors-be/internal/service"
)

type TagHandler struct {
	tagSvc service.TagService
}

func NewTagHandler(tagSvc service.TagService) *TagHandler {
	return &TagHandler{tagSvc: tagSvc}
}

func (h *TagHandler) Create() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req service.TagInput
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail("无效的请求体"))
			return
		}

		tag, err := h.tagSvc.Create(r.Context(), req)
		if err != nil {
			writeTagError(w, err)
			return
		}

		response.JSON(w, http.StatusCreated, response.Created(tag))
	}
}

func (h *TagHandler) GetByID() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := parseIDParam(r, "id", "无效的标签ID")
		if err != nil {
			response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
			return
		}

		tag, err := h.tagSvc.GetByID(r.Context(), id)
		if err != nil {
			writeTagError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(tag))
	}
}

func (h *TagHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tags, err := h.tagSvc.List(r.Context())
		if err != nil {
			writeTagError(w, err)
			return
		}

		response.JSON(w, http.StatusOK, response.OK(tags))
	}
}

func writeTagError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, service.ErrTagNameRequired),
		errors.Is(err, service.ErrTagNameTooLong):
		response.JSON(w, http.StatusBadRequest, response.Fail(err.Error()))
	case errors.Is(err, service.ErrTagAlreadyExists):
		response.JSON(w, http.StatusConflict, response.Conflict(err.Error()))
	case errors.Is(err, service.ErrTagNotFound):
		response.JSON(w, http.StatusNotFound, response.NotFound(err.Error()))
	default:
		response.JSON(w, http.StatusInternalServerError, response.ServerError("标签操作失败"))
	}
}
