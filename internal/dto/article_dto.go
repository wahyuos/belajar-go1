package dto

import "github.com/google/uuid"

// ArticleResponse adalah struktur yang HANYA berisi field yang ingin ditampilkan.
type ArticleResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	CreatedAt string    `json:"published_at"`
	CategoryName string `json:"category_name"`
	AuthorName string `json:"author_name"`
}

// Untuk respon list
type ArticleListResponse struct {
	ID        uuid.UUID `json:"id"`
	Title     string    `json:"title"`
	CategoryName string `json:"category_name"`
	AuthorName string `json:"author_name"`
}

type CreateArticleRequest struct {
	Title      string `json:"title" binding:"required,min=5"`
	Content    string `json:"content" binding:"required"`
	CategoryID string `json:"category_id" binding:"required,uuid"`
}

type UpdateArticleRequest struct {
	Title      string `json:"title" binding:"omitempty,min=5"`
	Content    string `json:"content" binding:"omitempty"`
	CategoryID string `json:"category_id" binding:"omitempty,uuid"`
}