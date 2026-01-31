package services

import (
	"kasir-api/domain"
	"kasir-api/repositories"
)

// CategoryService defines the interface for category business logic.
type CategoryService interface {
	GetAllCategories() ([]domain.Category, error)
	GetCategoryByID(id int) (domain.Category, error)
	CreateCategory(category domain.Category) (domain.Category, error)
	UpdateCategory(id int, category domain.Category) (domain.Category, error)
	DeleteCategory(id int) error
}

type categoryServiceImpl struct {
	categoryRepo repositories.CategoryRepository
	productRepo  repositories.ProductRepository
}

// NewCategoryService creates a new instance of CategoryService.
func NewCategoryService(categoryRepo repositories.CategoryRepository, productRepo repositories.ProductRepository) CategoryService {
	return &categoryServiceImpl{
		categoryRepo: categoryRepo,
		productRepo:  productRepo,
	}
}

func (s *categoryServiceImpl) GetAllCategories() ([]domain.Category, error) {
	categories, err := s.categoryRepo.GetAll()
	if err != nil {
		return nil, err
	}
	if len(categories) == 0 {
		return categories, nil
	}

	products, err := s.productRepo.GetAll()
	if err != nil {
		return nil, err
	}

	productMap := make(map[int][]domain.Produk)
	for _, p := range products {
		productMap[p.CategoryID] = append(productMap[p.CategoryID], p)
	}

	for i, c := range categories {
		if products, ok := productMap[c.ID]; ok {
			categories[i].Products = products
		}
	}

	return categories, nil
}

func (s *categoryServiceImpl) GetCategoryByID(id int) (domain.Category, error) {
	category, err := s.categoryRepo.GetByID(id)
	if err != nil {
		return category, err
	}

	products, err := s.productRepo.GetByCategoryID(id)
	if err != nil {
		return category, err
	}

	category.Products = products
	return category, nil
}

func (s *categoryServiceImpl) CreateCategory(category domain.Category) (domain.Category, error) {
	return s.categoryRepo.Create(category)
}

func (s *categoryServiceImpl) UpdateCategory(id int, category domain.Category) (domain.Category, error) {
	return s.categoryRepo.Update(id, category)
}

func (s *categoryServiceImpl) DeleteCategory(id int) error {
	return s.categoryRepo.Delete(id)
}
