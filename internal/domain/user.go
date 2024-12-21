package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Name        string     `json:"name" gorm:"not null"`
	Username    string     `json:"username" gorm:"uniqueIndex;not null"`
	Password    string     `json:"-" gorm:"not null"` // "-" means this field won't be included in JSON
	Gender      string     `json:"gender" gorm:"not null"`
	Email       string     `json:"email" gorm:"uniqueIndex;not null"`
	LastLoginAt *time.Time `json:"last_login_at"`
}
