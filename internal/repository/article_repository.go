package repository

import (
	"bejalar-dasar/internal/model"

	"gorm.io/gorm"
)

type ArticleRepository interface {
	FindAll() ([]model.Article, error)
	FindByID(id string) (model.Article, error)
	Create(article model.Article) (model.Article, error)
	Update(article model.Article) (model.Article, error)
	Delete(article model.Article) error
}

type articleRepository struct {
	db *gorm.DB
}

func NewArticleRepository(db *gorm.DB) ArticleRepository {
	return &articleRepository{db}
}

func (r *articleRepository) queryPreload() *gorm.DB {
    // Tambahkan Preload("User")
	return r.db.Preload("Category").Preload("User")
}

func (r *articleRepository) FindAll() ([]model.Article, error) {
	var articles []model.Article
	// Preload "Category" untuk mengambil detail kategori terkait
	err := r.queryPreload().Find(&articles).Error
	return articles, err
}

func (r *articleRepository) FindByID(id string) (model.Article, error) {
	var article model.Article
	err := r.queryPreload().Where("id = ?", id).First(&article).Error
	return article, err
}

func (r *articleRepository) Create(article model.Article) (model.Article, error) {
	err := r.db.Create(&article).Error
	// Load relasi Category setelah create agar response lengkap (opsional)
	r.db.Preload("Category").First(&article, article.ID)
	return article, err
}

func (r *articleRepository) Update(article model.Article) (model.Article, error) {
	err := r.db.Save(&article).Error
	// Load ulang agar data relasi terupdate di response
	r.db.Preload("Category").First(&article, article.ID)
	return article, err
}

func (r *articleRepository) Delete(article model.Article) error {
	return r.db.Delete(&article).Error
}