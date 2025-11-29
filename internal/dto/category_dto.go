package dto

import "github.com/google/uuid"

// CategoryResponse adalah struktur yang HANYA berisi field yang ingin ditampilkan.
type CategoryResponse struct {
	ID		uuid.UUID `json:"id"`
	Name	string `json:"name"`
}

// Untuk respon list
type CategoryListResponse struct {
	ID		uuid.UUID `json:"id"`
	Name	string `json:"name"`
}

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=3"`
}

type UpdateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=3"`
}