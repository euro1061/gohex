package repository

import "github.com/euro1061/gohex/internal/domain"

type ProductRepository interface {
    Create(product *domain.Product) error
    GetByID(id uint) (*domain.Product, error)
    GetAll() ([]domain.Product, error)
    Update(product *domain.Product) error
    Delete(id uint) error
}
