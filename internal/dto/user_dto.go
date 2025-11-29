package dto

// LoginRequest adalah DTO untuk menerima kredensial login
type LoginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}