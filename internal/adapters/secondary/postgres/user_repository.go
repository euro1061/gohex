package postgres

import (
	"github.com/euro1061/gohex/internal/domain/entity"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *userRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *entity.User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user.Password = string(hashedPassword)

	return r.db.Create(user).Error
}

func (r *userRepository) FindByID(id uint) (*entity.User, error) {
	var user entity.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *userRepository) FindAll() ([]entity.User, error) {
	var users []entity.User
	err := r.db.Find(&users).Error
	if err != nil {
		return nil, err
	}
	return users, nil
}

func (r *userRepository) Update(user *entity.User, id int) error {
	return r.db.Model(user).Where("id = ?", id).Updates(user).Error
}

func (r *userRepository) Delete(id int) error {
	return r.db.Unscoped().Delete(&entity.User{}, id).Error
}

func (r *userRepository) Login(email, password string) (*entity.User, error) {
	var user entity.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		return nil, err
	}

	return &user, nil
}
