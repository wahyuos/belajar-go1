package service

import (
	"bejalar-dasar/internal/dto"
	"bejalar-dasar/internal/model"
	"bejalar-dasar/internal/repository"
	"errors"

	"github.com/google/uuid"
)

type ArticleService interface {
	GetAll() ([]dto.ArticleListResponse, error)
	GetByID(id string) (dto.ArticleResponse, error)
	Create(req dto.CreateArticleRequest, userID string) (model.Article, error)
	Update(id string, req dto.UpdateArticleRequest) (model.Article, error)
	Delete(id string) error
}

type articleService struct {
	repo         repository.ArticleRepository
	categoryRepo repository.CategoryRepository
}

func NewArticleService(repo repository.ArticleRepository, catRepo repository.CategoryRepository) ArticleService {
	return &articleService{repo, catRepo}
}

// Fungsi pembantu untuk mapping tunggal
func mapToArticleResponse(article model.Article) dto.ArticleResponse {
    // Pastikan relasi Category dan User tidak nil
    categoryName := ""
    if article.Category != nil {
        categoryName = article.Category.Name
    }

	authorName := ""
    if article.User != nil { authorName = article.User.Username }

	return dto.ArticleResponse{
		ID:           article.ID,
		Title:        article.Title,
		Content:      article.Content,
		CreatedAt:    article.CreatedAt.Format("2025-11-29 15:15:15"),
		CategoryName: categoryName,
		AuthorName:   authorName,
	}
}

func (s *articleService) GetAll() ([]dto.ArticleListResponse, error) {
	// return s.repo.FindAll()

	articles, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	// Lakukan Mapping ke List Response DTO
	var responses []dto.ArticleListResponse
	for _, article := range articles {
        categoryName := ""
        if article.Category != nil {
            categoryName = article.Category.Name
        }
		authorName := ""
		if article.User != nil { authorName = article.User.Username }
		responses = append(responses, dto.ArticleListResponse{
			ID:           article.ID,
			Title:        article.Title,
			CategoryName: categoryName,
			AuthorName:   authorName,
		})
	}

	return responses, nil
}

func (s *articleService) GetByID(id string) (dto.ArticleResponse, error) {
	// return s.repo.FindByID(id)

	article, err := s.repo.FindByID(id)
	if err != nil {
		return dto.ArticleResponse{}, err
	}
	// Panggil fungsi mapping
	return mapToArticleResponse(article), nil
}

func (s *articleService) Create(req dto.CreateArticleRequest, userID string) (model.Article, error) {
	// Konversi String ke UUID
	catID, err := uuid.Parse(req.CategoryID)
	if err != nil {
		return model.Article{}, errors.New("invalid category ID format")
	}

	// Validasi apakah kategori ada (Opsional, tapi Good Practice)
	_, err = s.categoryRepo.FindByID(req.CategoryID)
	if err != nil {
		return model.Article{}, errors.New("category not found")
	}
	
	// Konversi UserID (string) ke UUID
    uID, err := uuid.Parse(userID)
	if err != nil {
		return model.Article{}, errors.New("invalid user ID format")
	}

	article := model.Article{
		Title:      req.Title,
		Content:    req.Content,
		CategoryID: catID,
		UserID: uID,
	}

	return s.repo.Create(article)
}

func (s *articleService) Update(id string, req dto.UpdateArticleRequest) (model.Article, error) {
	// Cari artikel lama
	article, err := s.repo.FindByID(id)
	if err != nil {
		return article, err
	}

	// Update field jika ada di request
	if req.Title != "" {
		article.Title = req.Title
	}
	if req.Content != "" {
		article.Content = req.Content
	}
	if req.CategoryID != "" {
		catID, err := uuid.Parse(req.CategoryID)
		if err == nil {
			// Validasi kategori jika diubah
			_, errCheck := s.categoryRepo.FindByID(req.CategoryID)
			if errCheck == nil {
				article.CategoryID = catID
			}
		}
	}

	return s.repo.Update(article)
}

func (s *articleService) Delete(id string) error {
	article, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(article)
}