package services

import (
	"fmt"
	"kasir-api/domain"
	"kasir-api/repositories"
)

// ProductService defines the interface for product business logic.
type ProductService interface {
	GetAllProducts() ([]domain.Produk, error)
	GetProductByID(id int) (domain.Produk, error)
	CreateProduct(product domain.Produk) (domain.Produk, error)
	UpdateProduct(id int, product domain.Produk) (domain.Produk, error)
	DeleteProduct(id int) error
}

type productServiceImpl struct {
	productRepo  repositories.ProductRepository
	categoryRepo repositories.CategoryRepository
}

// NewProductService creates a new instance of ProductService.
func NewProductService(productRepo repositories.ProductRepository, categoryRepo repositories.CategoryRepository) ProductService {
	return &productServiceImpl{
		productRepo:  productRepo,
		categoryRepo: categoryRepo,
	}
}

func (s *productServiceImpl) GetAllProducts() ([]domain.Produk, error) {
	return s.productRepo.GetAll()
}

func (s *productServiceImpl) GetProductByID(id int) (domain.Produk, error) {
	return s.productRepo.GetByID(id)
}

func (s *productServiceImpl) CreateProduct(product domain.Produk) (domain.Produk, error) {
	if product.CategoryID == 0 {
		return domain.Produk{}, fmt.Errorf("category_id harus diisi")
	}
	if _, err := s.categoryRepo.GetByID(product.CategoryID); err != nil {
		return domain.Produk{}, fmt.Errorf("category dengan id %d tidak ditemukan", product.CategoryID)
	}
	return s.productRepo.Create(product)
}

func (s *productServiceImpl) UpdateProduct(id int, product domain.Produk) (domain.Produk, error) {
	if product.CategoryID != 0 {
		if _, err := s.categoryRepo.GetByID(product.CategoryID); err != nil {
			return domain.Produk{}, fmt.Errorf("category dengan id %d tidak ditemukan", product.CategoryID)
		}
	}
	return s.productRepo.Update(id, product)
}

func (s *productServiceImpl) DeleteProduct(id int) error {
	return s.productRepo.Delete(id)
}
