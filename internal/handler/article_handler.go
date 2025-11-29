package handler

import (
	"bejalar-dasar/internal/dto"
	"bejalar-dasar/internal/service"
	"bejalar-dasar/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ArticleHandler struct {
	service service.ArticleService
}

func NewArticleHandler(service service.ArticleService) *ArticleHandler {
	return &ArticleHandler{service}
}

func (h *ArticleHandler) GetAll(c *gin.Context) {
	articles, err := h.service.GetAll()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch articles", err.Error())
		return
	}
	// cek data artikel
	if len(articles) == 0{
		// jika masih kosong
		response.Success(c, http.StatusOK, "No article yet", articles)
		return
	}
	response.Success(c, http.StatusOK, "List of articles", articles)
}

func (h *ArticleHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	article, err := h.service.GetByID(id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "Article not found", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Article detail", article)
}

func (h *ArticleHandler) Create(c *gin.Context) {
	var req dto.CreateArticleRequest
	// Validasi payload sesuai DTO binding
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validation error", err.Error())
		return
	}

	// Ambil UserID dari Context yang diset oleh Middleware
	userID, exists := c.Get("userID")
	if !exists {
		response.Error(c, http.StatusUnauthorized, "Authentication required", "UserID not found in context")
		return
	}

	// Lakukan type assertion ke string
    userIDStr, ok := userID.(string)
    if !ok {
        response.Error(c, http.StatusInternalServerError, "Internal error", "Invalid UserID type")
		return
    }

	_, err := h.service.Create(req, userIDStr)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create article", err.Error())
		return
	}
	response.Success(c, http.StatusCreated, "Article created", nil)
}

func (h *ArticleHandler) Update(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateArticleRequest
	
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validation error", err.Error())
		return
	}

	_, err := h.service.Update(id, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to update article", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Article updated", nil)
}

func (h *ArticleHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(id); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete article", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Article deleted", nil)
}