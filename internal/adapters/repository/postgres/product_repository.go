package postgres

import (
	"fmt"

	"github.com/euro1061/gohex/internal/domain"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	// Auto Migrate the schema
	if err := db.AutoMigrate(&domain.Product{}); err != nil {
		panic(fmt.Sprintf("error migrating database: %v", err))
	}

	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *domain.Product) error {
	result := r.db.Create(product)
	if result.Error != nil {
		return fmt.Errorf("error creating product: %v", result.Error)
	}
	return nil
}

func (r *ProductRepository) GetByID(id uint) (*domain.Product, error) {
	var product domain.Product
	result := r.db.First(&product, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting product: %v", result.Error)
	}
	return &product, nil
}

func (r *ProductRepository) GetAll() ([]domain.Product, error) {
	var products []domain.Product
	result := r.db.Find(&products)
	if result.Error != nil {
		return nil, fmt.Errorf("error getting products: %v", result.Error)
	}
	return products, nil
}

func (r *ProductRepository) Update(product *domain.Product) error {
	result := r.db.Save(product)
	if result.Error != nil {
		return fmt.Errorf("error updating product: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}

func (r *ProductRepository) Delete(id uint) error {
	result := r.db.Delete(&domain.Product{}, id)
	if result.Error != nil {
		return fmt.Errorf("error deleting product: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("product not found")
	}
	return nil
}
