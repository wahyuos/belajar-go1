package handler

import (
	"bejalar-dasar/internal/dto"
	"bejalar-dasar/internal/service"
	"bejalar-dasar/pkg/middleware"
	"bejalar-dasar/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	userService service.UserService
}

func NewAuthHandler(userService service.UserService) *AuthHandler {
	return &AuthHandler{userService}
}

// Login menangani permintaan otentikasi user
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	
	// 1. Validasi Request Body
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "Validation failed", err.Error())
		return
	}

	// 2. Verifikasi Kredensial
	userID, err := h.userService.VerifyCredentials(req.Username, req.Password)
	if err != nil {
		response.Error(c, http.StatusUnauthorized, "Authentication failed", err.Error())
		return
	}

	// 3. GENERATE TOKEN JWT
	token, err := middleware.GenerateToken(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to create token", err.Error())
		return
	}

	// 4. Kirim Respon Sukses
	response.Success(c, http.StatusOK, "Login successful", gin.H{
		"token": token,
		"type":  "Bearer",
		"expires_in": 3600, // Misal 1 jam atau sesuai konfigurasi token
	})
}