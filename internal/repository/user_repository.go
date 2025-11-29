package repository

import (
	"bejalar-dasar/internal/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	FindByUsername(username string) (model.User, error)
	// Tambahkan fungsi lain (Create, Update) di sini jika diperlukan
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

// Implementasi FindByUsername
func (r *userRepository) FindByUsername(username string) (model.User, error) {
	var user model.User
	// Cari user berdasarkan username
	err := r.db.Where("username = ?", username).First(&user).Error
	return user, err
}