package application

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/euro1061/gohex/internal/domain"
	"github.com/euro1061/gohex/internal/ports/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) Register(user *domain.User) error {
	// Check if username already exists
	existingUser, err := s.repo.GetByUsername(user.Username)
	if err != nil {
		return fmt.Errorf("error checking username: %v", err)
	}
	if existingUser != nil {
		return errors.New("username already exists")
	}

	// Check if email already exists
	existingUser, err = s.repo.GetByEmail(user.Email)
	if err != nil {
		return fmt.Errorf("error checking email: %v", err)
	}
	if existingUser != nil {
		return errors.New("email already exists")
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("error hashing password: %v", err)
	}
	user.Password = string(hashedPassword)

	// Create user
	if err := s.repo.Create(user); err != nil {
		return fmt.Errorf("error creating user: %v", err)
	}

	return nil
}

func (s *UserService) Login(username, password string) (string, error) {
	// Get user by username
	user, err := s.repo.GetByUsername(username)
	if err != nil {
		return "", fmt.Errorf("error getting user: %v", err)
	}
	if user == nil {
		return "", errors.New("invalid username or password")
	}

	// Check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", errors.New("invalid username or password")
	}

	// Update last login time
	now := time.Now()
	user.LastLoginAt = &now
	if err := s.repo.Update(user); err != nil {
		return "", fmt.Errorf("error updating last login time: %v", err)
	}

	claims := jwt.MapClaims{
		"id": user.ID,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", fmt.Errorf("error generating token: %v", err)
	}

	return tokenString, nil
}

func (s *UserService) Update(user *domain.User) error {
	// Check if username is taken by another user
	existingUser, err := s.repo.GetByUsername(user.Username)
	if err != nil {
		return fmt.Errorf("error checking username: %v", err)
	}
	if existingUser != nil && existingUser.ID != user.ID {
		return errors.New("username already taken")
	}

	// Check if email is taken by another user
	existingUser, err = s.repo.GetByEmail(user.Email)
	if err != nil {
		return fmt.Errorf("error checking email: %v", err)
	}
	if existingUser != nil && existingUser.ID != user.ID {
		return errors.New("email already taken")
	}

	// Clean input data
	user.Name = strings.TrimSpace(user.Name)
	user.Username = strings.TrimSpace(user.Username)
	user.Email = strings.TrimSpace(user.Email)
	user.Gender = strings.ToLower(strings.TrimSpace(user.Gender))

	// Update user
	return s.repo.Update(user)
}

func (s *UserService) GetUserFromToken(token string) (*domain.User, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})
	if err != nil {
		return nil, err
	}

	id, ok := claims["id"].(float64)
	if !ok {
		return nil, errors.New("invalid token")
	}

	user, err := s.repo.GetByID(uint(id))
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	return user, nil
}
