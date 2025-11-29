package middleware

import (
	"bejalar-dasar/pkg/response"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go/v4"
	"github.com/gin-gonic/gin"
)

// Claims JWT kustom Anda
type AuthClaims struct {
	UserID string `json:"user_id"`
	jwt.StandardClaims
}

// Secret key yang harus ada di .env Anda
var jwtSecret = []byte(os.Getenv("JWT_SECRET_KEY"))

// GenerateToken menciptakan token JWT baru
func GenerateToken(userID string) (string, error) {
	claims := AuthClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(time.Hour * 24)), // Token berlaku 24 jam
			IssuedAt:  jwt.At(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtSecret)
}

// AuthMiddleware adalah middleware yang memvalidasi token JWT
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Ambil header Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			response.Error(c, http.StatusUnauthorized, "Request rejected", "Authorization not found")
			c.Abort()
			return
		}

		// 2. Cek format Bearer Token
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			response.Error(c, http.StatusUnauthorized, "Request rejected", "Invalid token format")
			c.Abort()
			return
		}

		tokenString := parts[1]

		// 3. Verifikasi Token
		claims := &AuthClaims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("invalid signing method: %v", token.Header["alg"])
			}
			return jwtSecret, nil
		})

		if err != nil || !token.Valid {
			response.Error(c, http.StatusUnauthorized, "Request rejected", "Token is invalid or expired")
			c.Abort()
			return
		}

		// 4. Lolos: Simpan UserID di Context untuk digunakan di handler
		c.Set("userID", claims.UserID)
		c.Next()
	}
}