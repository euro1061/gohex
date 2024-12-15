package repository

import "github.com/euro1061/gohex/internal/domain/entity"

type UserRepository interface {
	Create(user *entity.User) error
	FindByID(id uint) (*entity.User, error)
	FindAll() ([]entity.User, error)
	Update(user *entity.User, id int) error
	Delete(id int) error
	Login(email, password string) (*entity.User, error)
}
