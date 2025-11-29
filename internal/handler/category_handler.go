package handler

import (
	"bejalar-dasar/internal/dto"
	"bejalar-dasar/internal/service"
	"bejalar-dasar/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	service service.CategoryService
}

func NewCategoryHandler(service service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service}
}

func (h *CategoryHandler) GetAll(c *gin.Context) {
	categories, err := h.service.GetAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch categories", err.Error())
		return
	}
	// cek apakah ada datanya
	if len(categories) == 0 {
		response.Success(c, http.StatusOK, "No category yet", categories)
		return
	}
	response.Success(c, http.StatusOK, "List of categories", categories)
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req dto.CreateCategoryRequest
	// Validasi form
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validation error", err.Error())
		return
	}

	_, err := h.service.Create(req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create category", err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "Category created", nil)
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validation error", err.Error())
		return
	}

	_, err := h.service.Update(id, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update category", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Category updated", nil)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete category", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Category deleted", nil)
}