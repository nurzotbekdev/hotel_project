package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name     string `json:"name" gorm:"size:90"`
	Surname  string `json:"surname" gorm:"size:90"`
	Phone    string `json:"phone" gorm:"size:30"`
	Email    string `json:"email" gorm:"size:255;uniqueIndex"`
	Password string `json:"password" gorm:"size:255"`
	Role     string `json:"role" gorm:"size:20;default:'user'"`
}
