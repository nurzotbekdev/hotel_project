package helper

import (
	"restaurant_manager/models"

	"gorm.io/gorm"
)

func IsEmailUnique(db *gorm.DB, email string) bool {
	var count int64
	db.Model(&models.User{}).Where("email = ?", email).Count(&count)
	return count == 0
}
