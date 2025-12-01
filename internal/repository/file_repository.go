package repository

import (
	"bejalar-dasar/internal/model"

	"gorm.io/gorm"
)

// Definisikan antarmuka (interface) untuk Repository
type FileRepository interface {
	Create(file *model.File) error
	GetByID(id string) (*model.File, error)
}

type fileRepository struct {
	db *gorm.DB
}

// Konstruktor
func NewFileRepository(db *gorm.DB) FileRepository {
	return &fileRepository{db}
}

// Implementasi metode Create: Menyimpan metadata file ke DB
func (r *fileRepository) Create(file *model.File) error {
	return r.db.Create(file).Error
}

func (r *fileRepository) GetByID(id string) (*model.File, error) {
	file := &model.File{}
	// Mencari berdasarkan Primary Key (ID)
	err := r.db.Where("id = ?", id).First(file).Error
	return file, err
}