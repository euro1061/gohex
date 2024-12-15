package input

import (
	"github.com/euro1061/gohex/internal/domain/entity"
	"github.com/euro1061/gohex/internal/domain/repository"
)

type UserService interface {
	CreateUser(user *entity.User) error
	GetUser(id uint) (*entity.User, error)
	GetAllUsers() ([]entity.User, error)
	UpdateUser(user *entity.User, id int) error
	DeleteUser(id int) error
	Login(email, password string) (*entity.User, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		userRepo: repo,
	}
}

// Implement methods
func (s *userService) CreateUser(user *entity.User) error {
	return s.userRepo.Create(user)
}

func (s *userService) GetUser(id uint) (*entity.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *userService) GetAllUsers() ([]entity.User, error) {
	return s.userRepo.FindAll()
}

func (s *userService) UpdateUser(user *entity.User, id int) error {
	return s.userRepo.Update(user, id)
}

func (s *userService) DeleteUser(id int) error {
	return s.userRepo.Delete(id)
}

func (s *userService) Login(email, password string) (*entity.User, error) {
	return s.userRepo.Login(email, password)
}
