package service

import (
	"bejalar-dasar/internal/repository"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	VerifyCredentials(username string, password string) (string, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

// VerifyCredentials mensimulasikan pencarian user dan validasi password.
// Dalam aplikasi nyata, Anda akan mencari user di DB, lalu membandingkan hash password (misalnya bcrypt).
func (s *userService) VerifyCredentials(username string, password string) (string, error) {
	// 1. Cari user berdasarkan username di database
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		// Jika user tidak ditemukan (GORM mengembalikan ErrRecordNotFound), kembalikan error "tidak cocok"
		return "", errors.New("username and password do not match") 
	}
	
	// 2. Bandingkan password yang dimasukkan dengan hash password dari database
	// bcrypt.CompareHashAndPassword menerima hash dari DB (user.Password) dan password plaintext (password)
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	
	if err != nil {
		// Jika perbandingan gagal (termasuk jika password tidak cocok), kembalikan error "tidak cocok"
		return "", errors.New("username and password do not match") 
	}
	
	// 3. Verifikasi berhasil, kembalikan User ID
	return user.ID.String(), nil
}