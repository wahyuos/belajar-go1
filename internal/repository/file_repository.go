package repository

import (
	"bejalar-dasar/internal/model"

	"gorm.io/gorm"
)

// Definisikan antarmuka (interface) untuk Repository
type FileRepository interface {
	Create(file *model.File) error
}

type fileRepository struct {
	DB *gorm.DB
}

// Konstruktor
func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{DB: db}
}

// Implementasi metode Create: Menyimpan metadata file ke DB
func (r *fileRepository) Create(file *model.File) error {
	return r.DB.Create(file).Error
}