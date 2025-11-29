package repository

import (
	"bejalar-dasar/internal/model"

	"gorm.io/gorm"
)

type CategoryRepository interface {
	FindAll() ([]model.Category, error)
	FindByID(id string) (model.Category, error)
	Create(category model.Category) (model.Category, error)
	Update(category model.Category) (model.Category, error)
	Delete(category model.Category) error
}

type categoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) CategoryRepository {
	return &categoryRepository{db}
}

func (r *categoryRepository) FindAll() ([]model.Category, error) {
	var categories []model.Category
	err := r.db.Find(&categories).Error
	return categories, err
}

func (r *categoryRepository) FindByID(id string) (model.Category, error) {
	var category model.Category
	err := r.db.Where("id = ?", id).First(&category).Error
	return category, err
}

func (r *categoryRepository) Create(category model.Category) (model.Category, error) {
	err := r.db.Create(&category).Error
	return category, err
}

func (r *categoryRepository) Update(category model.Category) (model.Category, error) {
	err := r.db.Save(&category).Error
	return category, err
}

func (r *categoryRepository) Delete(category model.Category) error {
	return r.db.Delete(&category).Error
}