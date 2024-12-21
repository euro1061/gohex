package application

import (
	"errors"
	"strings"

	"github.com/euro1061/gohex/internal/domain"
	"github.com/euro1061/gohex/internal/ports/repository"
)

var (
	ErrInvalidProductName        = errors.New("product name cannot be empty")
	ErrInvalidProductPrice       = errors.New("product price must be greater than 0")
	ErrInvalidProductDescription = errors.New("product description cannot be empty")
	ErrProductNotFound           = errors.New("product not found")
	ErrInvalidProductID          = errors.New("invalid product ID")
)

type ProductService struct {
	repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) *ProductService {
	return &ProductService{
		repo: repo,
	}
}

func (s *ProductService) validateProduct(name, description string, price float64) error {
	if strings.TrimSpace(name) == "" {
		return ErrInvalidProductName
	}
	if strings.TrimSpace(description) == "" {
		return ErrInvalidProductDescription
	}
	if price <= 0 {
		return ErrInvalidProductPrice
	}
	return nil
}

func (s *ProductService) CreateProduct(name, description string, price float64) (*domain.Product, error) {
	// Validate input
	if err := s.validateProduct(name, description, price); err != nil {
		return nil, err
	}

	// Clean input data
	name = strings.TrimSpace(name)
	description = strings.TrimSpace(description)

	product := &domain.Product{
		Name:        name,
		Description: description,
		Price:       price,
	}

	err := s.repo.Create(product)
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (s *ProductService) GetProduct(id uint) (*domain.Product, error) {
	if id == 0 {
		return nil, errors.New("invalid product ID")
	}

	product, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, ErrProductNotFound
	}
	return product, nil
}

func (s *ProductService) GetAllProducts() ([]domain.Product, error) {
	products, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	// If no products found, return empty slice instead of nil
	if products == nil {
		return []domain.Product{}, nil
	}
	return products, nil
}

func (s *ProductService) UpdateProduct(product *domain.Product) error {
	if product == nil {
		return errors.New("product cannot be nil")
	}

	if err := s.validateProduct(product.Name, product.Description, product.Price); err != nil {
		return err
	}

	// Check if product exists
	existing, err := s.repo.GetByID(product.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrProductNotFound
	}

	// Clean input data
	product.Name = strings.TrimSpace(product.Name)
	product.Description = strings.TrimSpace(product.Description)

	return s.repo.Update(product)
}

func (s *ProductService) DeleteProduct(id uint) error {
	if id == 0 {
		return ErrInvalidProductID
	}

	// Check if product exists before deleting
	existing, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrProductNotFound
	}

	return s.repo.Delete(id)
}
