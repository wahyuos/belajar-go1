package service

import (
	"bejalar-dasar/internal/dto"
	"bejalar-dasar/internal/model"
	"bejalar-dasar/internal/repository"
	"strings"
)

type CategoryService interface {
	GetAll() ([]dto.CategoryListResponse, error)
	GetByID(id string) (dto.CategoryResponse, error)
	Create(req dto.CreateCategoryRequest) (model.Category, error)
	Update(id string, req dto.UpdateCategoryRequest) (model.Category, error)
	Delete(id string) error
}

type categoryService struct {
	repo repository.CategoryRepository
}

func NewCategoryService(repo repository.CategoryRepository) CategoryService {
	return &categoryService{repo}
}

func (s *categoryService) GetAll() ([]dto.CategoryListResponse, error) {
	// return s.repo.FindAll()
	categories, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	// Lakukan Mapping ke List Response DTO
	var responses []dto.CategoryListResponse
	for _, category := range categories {
		responses = append(responses, dto.CategoryListResponse{
			ID:          category.ID,
			Name:        category.Name,
		})
	}

	return responses, nil
}

func (s *categoryService) GetByID(id string) (dto.CategoryResponse, error) {
	// return s.repo.FindByID(id)

	category, err := s.repo.FindByID(id)
	if err != nil {
		return dto.CategoryResponse{}, err
	}

	response := dto.CategoryResponse{
		ID:          category.ID,
		Name:        category.Name,
	}

	return response, nil
}

func (s *categoryService) Create(req dto.CreateCategoryRequest) (model.Category, error) {
	// Simple logic: Create slug from name
	slug := strings.ReplaceAll(strings.ToLower(req.Name), " ", "-")
	
	cat := model.Category{
		Name: req.Name,
		Slug: slug,
	}
	return s.repo.Create(cat)
}

func (s *categoryService) Update(id string, req dto.UpdateCategoryRequest) (model.Category, error) {
	cat, err := s.repo.FindByID(id)
	if err != nil {
		return cat, err
	}

	cat.Name = req.Name
	cat.Slug = strings.ReplaceAll(strings.ToLower(req.Name), " ", "-")
	
	return s.repo.Update(cat)
}

func (s *categoryService) Delete(id string) error {
	cat, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	return s.repo.Delete(cat)
}