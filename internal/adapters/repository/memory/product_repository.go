package memory

import (
	"errors"
	"sync"

	"github.com/euro1061/gohex/internal/domain"
)

type ProductRepository struct {
	sync.RWMutex
	products map[uint]*domain.Product
	nextID   uint
}

func NewProductRepository() *ProductRepository {
	return &ProductRepository{
		products: make(map[uint]*domain.Product),
		nextID:   1,
	}
}

func (r *ProductRepository) Create(product *domain.Product) error {
	r.Lock()
	defer r.Unlock()

	product.ID = r.nextID
	r.products[product.ID] = product
	r.nextID++
	return nil
}

func (r *ProductRepository) GetByID(id uint) (*domain.Product, error) {
	r.RLock()
	defer r.RUnlock()

	if product, exists := r.products[id]; exists {
		return product, nil
	}
	return nil, nil
}

func (r *ProductRepository) GetAll() ([]domain.Product, error) {
	r.RLock()
	defer r.RUnlock()

	products := make([]domain.Product, 0, len(r.products))
	for _, product := range r.products {
		products = append(products, *product)
	}
	return products, nil
}

func (r *ProductRepository) Update(product *domain.Product) error {
	r.Lock()
	defer r.Unlock()

	if _, exists := r.products[product.ID]; !exists {
		return errors.New("product not found")
	}

	r.products[product.ID] = product
	return nil
}

func (r *ProductRepository) Delete(id uint) error {
	r.Lock()
	defer r.Unlock()

	if _, exists := r.products[id]; !exists {
		return errors.New("product not found")
	}

	delete(r.products, id)
	return nil
}
