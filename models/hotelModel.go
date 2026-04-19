package models

import "gorm.io/gorm"

type Hotel struct {
	gorm.Model
	Name        string `json:"name" gorm:"size:110"`
	Address     string `json:"address" gorm:"size:110"`
	Description string `json:"description" gorm:"type:text"`
	Phone       string `json:"phone" gorm:"size:150"`
	Email       string `json:"email" gorm:"size:255;uniqueIndex"`

	HotelImages []HotelImage `json:"hotel_images" gorm:"foreignKey:HotelID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
