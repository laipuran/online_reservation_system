package handler

import (
	"net/http"

	"ors-be/internal/api/http/response"
	"ors-be/internal/service"
)

type CategoryHandler struct {
	categorySvc service.CategoryService
}

func NewCategoryHandler(categorySvc service.CategoryService) *CategoryHandler {
	return &CategoryHandler{categorySvc: categorySvc}
}

func (h *CategoryHandler) List() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		categories, err := h.categorySvc.List(r.Context())
		if err != nil {
			response.JSON(w, http.StatusInternalServerError, response.ServerError("分类列表查询失败"))
			return
		}

		response.JSON(w, http.StatusOK, response.OK(categories))
	}
}
