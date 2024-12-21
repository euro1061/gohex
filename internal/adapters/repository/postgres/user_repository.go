package postgres

import (
	"fmt"

	"github.com/euro1061/gohex/internal/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	// Auto Migrate the schema
	if err := db.AutoMigrate(&domain.User{}); err != nil {
		panic(fmt.Sprintf("error migrating database: %v", err))
	}

	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	result := r.db.Create(user)
	if result.Error != nil {
		return fmt.Errorf("error creating user: %v", result.Error)
	}
	return nil
}

func (r *UserRepository) GetByID(id uint) (*domain.User, error) {
	var user domain.User
	result := r.db.First(&user, id)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting user: %v", result.Error)
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting user: %v", result.Error)
	}
	return &user, nil
}

func (r *UserRepository) GetByEmail(email string) (*domain.User, error) {
	var user domain.User
	result := r.db.Where("email = ?", email).First(&user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("error getting user: %v", result.Error)
	}
	return &user, nil
}

func (r *UserRepository) Update(user *domain.User) error {
	result := r.db.Save(user)
	if result.Error != nil {
		return fmt.Errorf("error updating user: %v", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("user not found")
	}
	return nil
}
