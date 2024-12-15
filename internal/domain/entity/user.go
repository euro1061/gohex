package entity

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`
	Email     string `json:"email" gorm:"unique"`
	Password  string `json:"password"`
}

type Credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
